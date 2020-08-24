package main

import "github.com/vashish1/Dev/database"

type LoginResponse struct {
	Success bool   `json:"success"`
	Token   string `json:"token"`
}

type logn struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type Dasboard struct {
	Education []database.Education
	Experience []database.Experience
}
type p struct {
	Data string `json:"data"`
}

type mockSignup struct {
	Name      string `json:"name"`
	Email     string `json:"email"`
	Password  string `json:"password"`
	Cpassword string `json:"cpassword"`
}

type str struct {
	Str string
}