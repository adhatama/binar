package server

import "errors"

type SignupFormRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (r SignupFormRequest) Validate() map[string]interface{} {
	errs := make(map[string]interface{})

	nameErrs := []string{}
	if r.Name == "" {
		nameErrs = append(nameErrs, "Cannot be blank")
	}

	emailErrs := []string{}
	if r.Email == "" {
		emailErrs = append(emailErrs, "Cannot be blank")
	}

	passwordErrs := []string{}
	if r.Password == "" {
		passwordErrs = append(passwordErrs, "Cannot be blank")
	}

	if len(nameErrs) > 0 {
		errs["name"] = nameErrs
	}
	if len(emailErrs) > 0 {
		errs["email"] = emailErrs
	}
	if len(passwordErrs) > 0 {
		errs["password"] = passwordErrs
	}

	return errs
}

type LoginFormRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (r LoginFormRequest) Validate() map[string]interface{} {
	errs := make(map[string]interface{})

	if r.Email == "" {
		errs["email"] = []error{errors.New("Cannot be blank")}
	}
	if r.Password == "" {
		errs["password"] = []error{errors.New("Cannot be blank")}
	}

	return errs
}
