package main

import (
	"net/http"
	"reqtemplates/Dev-Connect/database"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/mongo"
)

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/Dev", index)
	http.Handle("/", r)
	http.ListenAndServe(":80", nil)
}

var cl *mongo.Collection
var c *mongo.Client

func init() {
	cl, c = database.Createdb()
}

//index handles the main page
func index(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte("hello its dev connect"))
}

//login handles the login credentials
func login(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	err := r.ParseForm()
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"error": "body not parsed"}`))
		return
	}
	email := r.FormValue("email")
	pass := r.FormValue("password")
	ok := database.Findfromuserdb(cl, email, pass)
	if ok {
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte(`{"success": "created"}`))
		return
	}
	w.WriteHeader(http.StatusNotFound)
	w.Write([]byte(`{"error": "not created"}`))
}

//signup handles the login credentials
func signup(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	err := r.ParseForm()
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"error": "body not parsed"}`))
		return
	}
	var u database.User
	name := r.FormValue("name")
	email := r.FormValue("email")
	pass := r.FormValue("password")
	cnf := r.FormValue("confirm-password")
	if pass == cnf {
		u = database.Newuser(name, email, pass)
	}

	ok := database.Insertintouserdb(cl, u)
	if ok {
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte(`{"success": "created"}`))
		return
	}
	w.WriteHeader(http.StatusNotFound)
	w.Write([]byte(`{"error": "not created"}`))
}