package main

import (
	"Dev/database"
	"encoding/json"
	"io/ioutil"
	"net/http"
    "fmt"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/mongo"
)

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/Dev", index)
	r.HandleFunc("/Dev/login", login).Methods("GET")
	r.HandleFunc("/Dev/signup", signup).Methods("POST")
	r.HandleFunc("/Dev/profile/AddEducation/{name}", education).Methods("POST")
	r.HandleFunc("/Dev/profile/AddExperience/{name}", experience).Methods("POST")
	r.HandleFunc("/Dev/profile/{name}", profile).Methods("GET", "POST")
	http.Handle("/", r)
	http.ListenAndServe(":80", nil)
}

var cl, cl1 *mongo.Collection
var c *mongo.Client

func init() {
	cl, cl1, c = database.Createdb()
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
	var user database.User
	body, _ := ioutil.ReadAll(r.Body)
	err := json.Unmarshal(body, &user)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"error": "body not parsed"}`))
		return
	}

	ok := database.Insertintouserdb(cl, user)
	if ok {
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte(`{"success": "created"}`))
		return
	}
	w.WriteHeader(http.StatusNotFound)
	w.Write([]byte(`{"error": "not created"}`))
}

//profile updates the profile
func profile(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params:=mux.Vars(r)
	fmt.Println(params)
	p:=database.Finddb(cl,params["name"])
	fmt.Println(p)

	var pro database.Profile
	body, _ := ioutil.ReadAll(r.Body)
	err := json.Unmarshal(body, &pro)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"error": "body not parsed"}`))
		return
	}
	pro.Email=p.Email;
	ok:=database.Insertprofile(cl1,pro)
	if ok{
		
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte(`{"success": "created"}`))
		return
	}

}

//education updates the education
func education(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params:=mux.Vars(r)
	fmt.Println(params)
	p:=database.Finddb(cl,params["name"])
	var edu database.Education
	body, _ := ioutil.ReadAll(r.Body)
	err := json.Unmarshal(body, &edu)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"error": "body not parsed"}`))
		return
	}
	ok := database.Updateeducation(cl1, p.Email, edu)
	if ok {
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte(`{"success": "created"}`))
		return
	}
	w.WriteHeader(http.StatusNotFound)
	w.Write([]byte(`{"error": "not created"}`))
}

//experience updates the experience
func experience(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params:=mux.Vars(r)
	fmt.Println(params)
	p:=database.Finddb(cl,params["name"])
	var exp database.Experience
	body, _ := ioutil.ReadAll(r.Body)
	err := json.Unmarshal(body, &exp)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"error": "body not parsed"}`))
		return
	}

	ok := database.Updateexperience(cl1, p.Email,exp)
	if ok {
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte(`{"success": "created"}`))
		return
	}
	w.WriteHeader(http.StatusNotFound)
	w.Write([]byte(`{"error": "not created"}`))

}
