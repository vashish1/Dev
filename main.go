package main

import (
	"Dev/database"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/mongo"
)

type logn struct{
Email string
Pass string
}

type mockSignup struct {
	Name      string `json:"name"`
	Email     string `json:"email"`
	Password  string `json:"password"`
	Cpassword string `json:"cpassword"`
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/Dev", index)
	r.HandleFunc("/Dev/login", login).Methods("POST")
	r.HandleFunc("/Dev/signup", signup).Methods("POST")
	r.HandleFunc("/Dev/profile/AddEducation/{name}", education).Methods("POST")
	r.HandleFunc("/Dev/profile/AddExperience/{name}", experience).Methods("POST")
	r.HandleFunc("/Dev/profile/", profile).Methods("GET", "POST")
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
	var  result database.User
	var user logn
	body, _ := ioutil.ReadAll(r.Body)
	err := json.Unmarshal(body, &user)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"error": "body not parsed"}`))
		return
	}
	ok := database.Findfromuserdb(cl, user.Email, user.Pass)
	if ok {
		u:=database.Finddb(cl,user.Email)
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"name":  u.Name,
			"email": u.Email,
			
		})

		tokenString, err := token.SignedString([]byte("secret"))

		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
		    w.Write([]byte(`{"error": "error in token string"}`))
		    return
		}
		result.Token=tokenString
		result.PasswordHash=""
		tkn:=database.UpdateToken(cl,u.Email,tokenString)
		if tkn{
			json.NewEncoder(w).Encode(result)
			w.WriteHeader(http.StatusCreated)
			w.Write([]byte(`{"success": "created token successfully"}`))
			return
		}
	}
	w.WriteHeader(http.StatusNotFound)
	w.Write([]byte(`{"error": "not created"}`))
}

//signup handles the login credentials
func signup(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	test:=mockSignup{}
	var user database.User
	body, _ := ioutil.ReadAll(r.Body)
	err := json.Unmarshal(body, &test)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"error": "body not parsed"}`))
		return
	}
	fmt.Println("test:",test)
    if test.Password==test.Cpassword{
    	user=database.Newuser(test.Name,test.Email,test.Password,"")
		ok := database.Insertintouserdb(cl, user)
		if ok {
			w.WriteHeader(http.StatusCreated)
			w.Write([]byte(`{"success": "created"}`))
			return
		}
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte(`{"error": "not created"}`))
		}
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(`{"error": "Passwords do not match"}`))
	}

//profile updates the profile
func profile(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	tokenString := r.Header.Get("Authorization")
	token,_ := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method")
		} 
		return []byte("secret"), nil
	})
	var result database.User
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		result.Name = claims["name"].(string)
		result.Email = claims["email"].(string)
		json.NewEncoder(w).Encode(result)
	} 
	p := database.Finddb(cl, result.Email)
	fmt.Println(p)
	if p.Name!=""{
         var pro database.Profile
	body, _ := ioutil.ReadAll(r.Body)
	err := json.Unmarshal(body, &pro)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"error": "body not parsed"}`))
		return
	}
	pro.Email = p.Email
	ok := database.Insertprofile(cl1, pro)
	if ok {

		w.WriteHeader(http.StatusCreated)
		w.Write([]byte(`{"success": "created"}`))
		return
	}
	}
	w.WriteHeader(http.StatusBadRequest)
	w.Write([]byte(`{"error": "Token not verified"}`))	

	}
	// params := mux.Vars(r)
	// fmt.Println(params)

	// var pro database.Profile
	// body, _ := ioutil.ReadAll(r.Body)
	// err := json.Unmarshal(body, &pro)
	// if err != nil {
	// 	w.WriteHeader(http.StatusBadRequest)
	// 	w.Write([]byte(`{"error": "body not parsed"}`))
	// 	return
	// }
	// pro.Email = p.Email
	// ok := database.Insertprofile(cl1, pro)
	// if ok {

	// 	w.WriteHeader(http.StatusCreated)
	// 	w.Write([]byte(`{"success": "created"}`))
	// 	return
	// }


//education updates the education
func education(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	fmt.Println(params)
	p := database.Finddb(cl, params["name"])
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
	params := mux.Vars(r)
	fmt.Println(params)
	p := database.Finddb(cl, params["name"])
	var exp database.Experience
	body, _ := ioutil.ReadAll(r.Body)
	err := json.Unmarshal(body, &exp)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"error": "body not parsed"}`))
		return
	}

	ok := database.Updateexperience(cl1, p.Email, exp)
	if ok {
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte(`{"success": "created"}`))
		return
	}
	w.WriteHeader(http.StatusNotFound)
	w.Write([]byte(`{"error": "not created"}`))

}
