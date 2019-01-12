package server

import "errors"

type SignupFormRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (r SignupFormRequest) Validate() map[string][]error {
	errs := make(map[string][]error)

	if r.Name == "" {
		errs["name"] = append(errs["name"], errors.New("Name cannot be empty"))
	}
	if r.Email == "" {
		errs["email"] = append(errs["email"], errors.New("Email cannot be empty"))
	}
	if r.Password == "" {
		errs["password"] = append(errs["password"], errors.New("Password cannot be empty"))
	}

	return errs
}

type LoginFormRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (r LoginFormRequest) Validate() map[string][]error {
	errs := make(map[string][]error)

	if r.Email == "" {
		errs["email"] = append(errs["email"], errors.New("Email cannot be empty"))
	}
	if r.Password == "" {
		errs["password"] = append(errs["password"], errors.New("Password cannot be empty"))
	}

	return errs
}
