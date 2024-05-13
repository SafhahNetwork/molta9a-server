package db

import (
	"github.com/lib/pq"
	"gorm.io/gorm"
)

type Course struct {
	gorm.Model

	Name        string
	Description *string
	Image       *string

	Difficulty *uint8         `gorm:"default:1"`
	Languages  pq.StringArray `gorm:"type:text[]"`
	Tags       pq.StringArray `gorm:"type:text[]"`

	TotalXp        uint
	TotalExercises uint
	TotalUsers     uint
}

type Quest struct {
	gorm.Model
	CourseID uint

	Name string
}

type Exercise struct {
	gorm.Model

	Index uint // ID relative to the course

	CourseID uint
	QuestID  uint

	Name string
	Type uint8 // 0=lesson 1=exercise

	Markdown       *string
	ExerciseConfig *string

	TotalXp    uint
	TotalUsers uint
}

type UserCourse struct {
	gorm.Model
	UserID   uint
	CourseID uint
	Status   string // active, completed, frozen
	Progress uint8

	Course *Course `gorm:"foreignkey:course_id"`
}

type UserExercise struct {
	gorm.Model
	UserID     uint
	CourseID   uint
	QuestID    uint
	ExerciseID uint
	TotalXp    uint

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
	TotalCourses   uint
}

// type UserActivity struct {
// 	UserID    uint
// 	ExerciseID  uint
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

type ExerciseConfig struct {
	DockerImage string
	Env         []string
	Args        []string
	Files       map[string]*File
}
