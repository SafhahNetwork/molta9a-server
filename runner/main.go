package main

import (
	"encoding/json"
	"flag"
	"io"
	"molta9a/server/runner/run"
	"net/http"
	"strings"
	"time"
)

func handle(rw http.ResponseWriter, r *http.Request) {
	// Parse URL path, query
	ip := strings.Split(r.Header.Get("X-Forwarded-For"), ", ")[0]
	if ip == "" {
		// No proxy detected, fallback to real IP (IPv4 or IPv6)
		parts := strings.Split(r.RemoteAddr, ":")
		ip = strings.Join(parts[:len(parts)-1], ":") // remove port
	}
	image := strings.Trim(r.URL.Path, "/")
	query := r.URL.Query()
	env := query["env"]
	args := query["args"]

	// Read body
	r.Body = http.MaxBytesReader(nil, r.Body, http.DefaultMaxHeaderBytes)
	body, err := io.ReadAll(r.Body)
	if err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(rw).Encode(struct {
			Output string
			Ok     bool
		}{err.Error(), false})
		return
	}

	b, ok, err := run.Run(ip, image, env, args, body)
	message := string(b)
	if err != nil {
		message = err.Error()
		rw.WriteHeader(http.StatusBadRequest)
	}
	json.NewEncoder(rw).Encode(struct {
		Output string
		Ok     bool
	}{message, ok})
}

func main() {
	http.HandleFunc("/", handle)
	port := flag.String("port", "6969", "listening port")
	flag.Parse()
	srv := &http.Server{
		Addr:         ":" + *port,
		ReadTimeout:  time.Minute,
		WriteTimeout: time.Minute,
	}
	srv.ListenAndServe()
}
