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
	Name       string                `json:"name,omitempty"`
	Education  []database.Education  `json:"education,omitempty"`
	Experience []database.Experience `json:"experience,omitempty"`
}
type p struct {
	Text string `json:"data"`
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


