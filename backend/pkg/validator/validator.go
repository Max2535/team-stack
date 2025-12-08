package validator

import (
	"regexp"
	"strings"
)

var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)

type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

type ValidationErrors []ValidationError

func (v ValidationErrors) Error() string {
	if len(v) == 0 {
		return ""
	}
	var msgs []string
	for _, e := range v {
		msgs = append(msgs, e.Field+": "+e.Message)
	}
	return strings.Join(msgs, "; ")
}

func (v ValidationErrors) HasErrors() bool {
	return len(v) > 0
}

type LoginRequest struct {
	Email    string
	Password string
}

func ValidateLoginRequest(req LoginRequest) ValidationErrors {
	var errs ValidationErrors

	// Email validation
	email := strings.TrimSpace(req.Email)
	if email == "" {
		errs = append(errs, ValidationError{Field: "email", Message: "email is required"})
	} else if !emailRegex.MatchString(email) {
		errs = append(errs, ValidationError{Field: "email", Message: "invalid email format"})
	}

	// Password validation
	if req.Password == "" {
		errs = append(errs, ValidationError{Field: "password", Message: "password is required"})
	} else if len(req.Password) < 6 {
		errs = append(errs, ValidationError{Field: "password", Message: "password must be at least 6 characters"})
	}

	return errs
}
