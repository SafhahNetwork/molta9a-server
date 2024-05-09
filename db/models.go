package db

import (
	"github.com/lib/pq"
	"gorm.io/gorm"
)

type Track struct {
	gorm.Model

	Name        string
	Description *string
	Image       *string

	Difficulty *uint8         `gorm:"default:1"`
	Languages  pq.StringArray `gorm:"type:text[]"`
	Tags       pq.StringArray `gorm:"type:text[]"`

	TotalXp      uint
	TotalModules uint
	TotalUsers   uint
}

type Quest struct {
	gorm.Model
	TrackID uint

	Name string
}

type Module struct {
	gorm.Model

	Index uint // ID relative to the track

	TrackID uint
	QuestID uint

	Name string
	Type uint8 // 0=lesson 1=exercise

	Markdown *string
	Exercise *string

	TotalXp    uint
	TotalUsers uint
}

type UserTrack struct {
	gorm.Model
	UserID   uint
	TrackID  uint
	Status   string // active, completed, frozen
	Progress uint8

	Track *Track `gorm:"foreignkey:track_id"`
}

type UserModule struct {
	gorm.Model
	UserID   uint
	TrackID  uint
	QuestID  uint
	ModuleID uint
	TotalXp  uint

	Status string
	Files  *string
	Output *string
}

type User struct {
	gorm.Model
	Username string
	Email    string
	Password string

	Avatar         *string
	TotalXp        uint
	TotalExercises uint
	TotalLessons   uint
	TotalTracks    uint
}

// type UserActivity struct {
// 	UserID    uint
// 	ModuleID  uint
// 	TotalXp   uint
// 	CreatedAt time.Time
// }

// Not directly in DB

type File struct {
	Name     string
	Value    string
	Language string
	ReadOnly bool
}

type Exercise struct {
	DockerImage string
	Env         []string
	Args        []string
	Files       map[string]*File
}
