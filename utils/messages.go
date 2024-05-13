package utils

import "errors"

// Errors
var (
	ErrConflictUsername                = errors.New("USERNAME_CONFLICT")
	ErrConflictEmail                   = errors.New("EMAIL_CONFLICT")
	ErrInternalServerError             = errors.New("INTERNAL_ERROR")
	ErrBadLogin                        = errors.New("BAD_LOGIN")
	ErrMissingAuthorizationToken       = errors.New("MISSING_AUTHORIZATION_TOKEN")
	ErrInvalidAuthorizationToken       = errors.New("INVALID_AUTHORIZATION_TOKEN")
	ErrInvalidAuthorizationTokenClaims = errors.New("INVALID_AUTHORIZATION_TOKEN_CLAIMS")
	ErrSomethingWentWrong              = errors.New("SOMETHING_WENT_WRONG")

	ErrExerciseNotFound   = errors.New("Exercise_NOT_FOUND")
	ErrDidntStartExercise = errors.New("EXERCISE_NOT_STARTED")

	ErrModuleFinished = errors.New("MODULE_FINISHED")

	ErrCourseNotFound       = errors.New("COURSE_NOT_FOUND")
	ErrAlreadyStartedCourse = errors.New("COURSE_ALREADY_STARTED")
	ErrDidntStartCourse     = errors.New("COURSE_NOT_STARTED")

	ErrUserNotFound = errors.New("USER_NOT_FOUND")
)

// Validation
var (
	InvalidField         = errors.New("BAD_FIELD")
	InvalidFieldRequired = errors.New("REQUIRED_FIELD")
	InvalidMinLength     = errors.New("UNMET_MINIMUM")
	InvalidMaxLength     = errors.New("UNMET_MAXIMUM")
	InvalidLength        = errors.New("UNMET_LENGTH")
	InvalidEmail         = errors.New("BAD_EMAIL")
)

// Info
const (
	InfoUserCreated = "USER_CREATED"
	InfoSuccess     = "SUCCESS"
)
