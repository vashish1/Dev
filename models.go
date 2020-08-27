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
	Name       string                `json:"name"`
	Education  []database.Education  `json:"education"`
	Experience []database.Experience `json:"experience"`
}
type p struct {
	Text string `json:"text"`
}

type mockSignup struct {
	Name      string `json:"name"`
	Email     string `json:"email"`
	Password  string `json:"password"`
	Cpassword string `json:"cpassword"`
}
