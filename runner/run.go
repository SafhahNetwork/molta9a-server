package run

import (
	"archive/tar"
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"molta9a/server/utils"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	dimage "github.com/docker/docker/api/types/image"
	"github.com/docker/docker/api/types/registry"
	"github.com/docker/docker/api/types/volume"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/stdcopy"
)

var (
	ctx = context.Background()
	cli *client.Client

	labels = map[string]string{"com.molta9a.type": "test"}
)

const (
	REFRESH_DOCKER_IMAGE_INTERVAL = time.Second
)

func init() {
	var err error
	cli, err = client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil || !errors.Is(err, nil) {
		panic(err)
	}
}

var users_ids = struct {
	sync.Mutex
	m map[string]int
}{m: map[string]int{}}

var images_updates = struct {
	sync.Mutex
	m map[string]time.Time
}{m: map[string]time.Time{}}

func logDuration(tag string) func(text string) {
	t := time.Now()
	return func(text string) {
		now := time.Now()
		secs := float64(now.Sub(t)) / float64(time.Second)
		log.Println("", tag, "", fmt.Sprintf("%.3fs", secs), "", text)
	}
}

func Run(ip string, image string, env []string, args []string, files []utils.File) ([]byte, bool, error) {

	// keeps track of how many times user called func
	users_ids.Lock()
	users_ids.m[ip]++
	iip := "#" + strconv.Itoa(users_ids.m[ip])
	users_ids.Unlock()

	// print seconds from now until the func is done
	defer logDuration(iip)("total")

	// func to count time between events
	logDuration := logDuration(iip)

	// convert the files into a tar_archive
	files = append(files, utils.File{
		Name:  "/",
		Input: "",
	})
	tarBuffer, err := utils.CreateTar(files, &utils.CreateTarOptions{
		FileMode: 0777,
		FileUid:  1000,
		FileGid:  1001,
		DirMode:  0777,
		DirUid:   1000,
		DirGid:   1001,
	})
	if err != nil {
		log.Println(err)
		return nil, false, err
	}

	logDuration("tar archive creation")

	// Refresh docker image every 1min
	images_updates.Lock()
	if time.Since(images_updates.m[image]) > REFRESH_DOCKER_IMAGE_INTERVAL {
		// Try to pull the image
		err := func() error {
			var options dimage.PullOptions
			// Log in to pull the image from github
			if strings.HasPrefix(image, "ghcr.io/") {
				b, err := json.Marshal(registry.AuthConfig{
					Username: os.Getenv("GITHUB_USERNAME"),
					Password: os.Getenv("GITHUB_TOKEN"),
				})
				if err != nil {
					return err
				}
				options.RegistryAuth = base64.URLEncoding.EncodeToString(b)
			}
			resp, err := cli.ImagePull(ctx, image, options)
			if err != nil {
				return err
			}
			io.Copy(io.Discard, resp)
			resp.Close()
			return nil
		}()
		if err != nil {
			if _, _, err := cli.ImageInspectWithRaw(ctx, image); err != nil {
				// The image doesn't exist
				images_updates.Unlock()
				return nil, false, err
			}
			// Failed to pull the image but an old version exists
			log.Println("Failed to pull the image but an old version exists, err: ", err)
		}
		images_updates.m[image] = time.Now()
	}
	images_updates.Unlock()
	logDuration("Pulled docker image: " + image)

	// Create the volume that will contain the code to test
	volumeStudent, err := cli.VolumeCreate(ctx, volume.CreateOptions{Labels: labels})
	if err != nil {
		return nil, false, err
	}

	logError := func(event string, err error) {
		if err != nil {
			event += ": " + err.Error()
		}
		logDuration(event)
	}

	defer func() {
		logError("volume remove", cli.VolumeRemove(ctx, volumeStudent.Name, false))
	}()

	// helper to remove a container
	containerRemove := func(id string) {
		logError("container remove", cli.ContainerRemove(ctx, id, container.RemoveOptions{
			RemoveVolumes: true,
			Force:         true,
		}))
	}

	// Create a container  whose sole purpose is to copy the TAR archive to the volume
	resp, err := cli.ContainerCreate(ctx, &container.Config{
		Image:  image,
		Labels: labels,
	}, &container.HostConfig{
		Binds: []string{volumeStudent.Name + ":/data"},
	}, nil, nil, "")
	if err != nil {
		return nil, false, err
	}
	logDuration("container creation")

	err = cli.CopyToContainer(ctx, resp.ID, "/data", tarBuffer, types.CopyToContainerOptions{
		CopyUIDGID: true,
	})
	containerRemove(resp.ID)
	if err != nil {
		log.Println(err)
		return nil, false, err
	}
	logDuration("container copy")

	// Create the test container with security constraints
	maxPID := int64(256)
	hostConfig := container.HostConfig{
		LogConfig: container.LogConfig{
			Type: "json-file",
			Config: map[string]string{
				"max-size": "1m",
				"max-file": "2",
			},
		},
		ReadonlyRootfs: true,
		Resources: container.Resources{
			PidsLimit: &maxPID,
			Memory:    500e6,
			NanoCPUs:  2e9,
		},
		NetworkMode: "none",
		Tmpfs:       map[string]string{"/jail": "size=100M,noatime,exec,nodev,nosuid,uid=1000,gid=1001,nr_inodes=5k,mode=1770"},
		Binds:       []string{volumeStudent.Name + ":/jail/student:rw"},
	}
	// Redirect a domain name to localhost if needed
	for _, val := range env {
		if strings.HasPrefix(val, "DOMAIN=") {
			hostConfig.ExtraHosts = []string{
				strings.TrimPrefix(val, "DOMAIN=") + ":127.0.0.1"}
		}
	}
	args = []string{
		"accumulate",
		"/jail/student",
	}
	resp, err = cli.ContainerCreate(ctx, &container.Config{
		Image:      image,
		User:       "1000:1001",
		WorkingDir: "/jail",
		Env:        append(env, "HOME=/jail", "TMPDIR=/jail"),
		Labels:     labels,
		Cmd:        args,
	}, &hostConfig, nil, nil, "")
	if err != nil {
		return nil, false, err
	}
	logDuration("container creation")

	if err := cli.ContainerStart(ctx, resp.ID, container.StartOptions{}); err != nil {
		containerRemove(resp.ID)
		return nil, false, err
	}
	logDuration("container start")

	ctxTimeout, cancel := context.WithTimeout(ctx, 15*time.Second)
	defer cancel()

	statusCh, errCh := cli.ContainerWait(ctxTimeout, resp.ID, container.WaitConditionNotRunning)
	var ok bool
	select {
	case err := <-errCh:
		if errors.Is(err, context.DeadlineExceeded) {
			if err := cli.ContainerKill(ctx, resp.ID, "SIGKILL"); err != nil {
				containerRemove(resp.ID)
				return nil, false, err
			}
		}
		containerRemove(resp.ID)
		return nil, false, errors.New("timeout: Did you write an infinite loop? " + err.Error())
	case status := <-statusCh:
		ok = status.StatusCode == 0
		if !ok && status.Error != nil {
			panic(status.Error.Message)
		}
	}

	logDuration("container stop")

	var outResult []byte
	resRead, _, err := cli.CopyFromContainer(ctx, resp.ID, "/jail/student/results.json")
	if err != nil {
		log.Println(err)
	} else {
		tr := tar.NewReader(resRead)
		tr.Next()
		outResult, err = io.ReadAll(tr)
		if err != nil {
			log.Println(err)
			containerRemove(resp.ID)
		}
		log.Println(string(outResult))
	}

	// Read Logs

	logsStream, err := cli.ContainerLogs(ctx, resp.ID, container.LogsOptions{
		ShowStdout: true,
		ShowStderr: true,
	})
	if err != nil {
		log.Println(err)
		containerRemove(resp.ID)
		return nil, false, nil
	}

	logDuration("container logs")
	containerRemove(resp.ID)

	// demux stream, see: https://docs.docker.com/engine/api/v1.41/#operation/ContainerAttach
	var outLogs bytes.Buffer
	if _, err := stdcopy.StdCopy(&outLogs, &outLogs, logsStream); err != nil {
		containerRemove(resp.ID)
		return nil, false, err
	}

	b := outLogs.Bytes()
	if len(b) > http.DefaultMaxHeaderBytes {
		b = b[:http.DefaultMaxHeaderBytes]
		b = append(b, []byte(" ... TRUNCATED")...)
	}

	return outResult, ok, nil
}
