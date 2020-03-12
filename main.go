package main

import (
	"Dev/database"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/mongo"
)

type logn struct {
	Email string
	Pass  string
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
	r.HandleFunc("/Dev/profile/AddEducation", education).Methods("POST")
	r.HandleFunc("/Dev/profile/AddExperience", experience).Methods("POST")
	r.HandleFunc("/Dev/profile", profile).Methods("GET", "POST")
	r.HandleFunc("/Dev/MyProfile/{name}", Myprofile).Methods("GET")
	r.HandleFunc("/Dev/dashboard/{name}", dashboard).Methods("GET")
	r.HandleFunc("/Dev/Post",WritePost).Methods("GET","POST")
	http.Handle("/", r)
	http.ListenAndServe(":3000", nil)
}

var cl, cl1,cl2 *mongo.Collection
var c *mongo.Client

func init() {
	cl, cl1, cl2, c = database.Createdb()
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
	var result database.User
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
		u := database.Finddb(cl, user.Email)
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
		result.Token = tokenString
		result.PasswordHash = ""
		tkn := database.UpdateToken(cl, u.Email, tokenString)
		if tkn {
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
	test := mockSignup{}
	var user database.User
	body, _ := ioutil.ReadAll(r.Body)
	err := json.Unmarshal(body, &test)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"error": "body not parsed"}`))
		return
	}
	fmt.Println("test:", test)
	if test.Password == test.Cpassword {
		user = database.Newuser(test.Name, test.Email, test.Password, "")
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
	
	tokenString = strings.TrimPrefix(tokenString, "Bearer ")
	fmt.Println("token", tokenString)

	token, _ := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method")
		}
		return []byte("secret"), nil
	})
	// var result database.User
	var name, email string
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		name = claims["name"].(string)
		email = claims["email"].(string)
	}
	fmt.Println(name, email)
	p := database.Finddb(cl, email)
	fmt.Println(p)
	if p.Name != "" {
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

//education updates the education
func education(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	tokenString := r.Header.Get("Authorization")

	tokenString = strings.TrimPrefix(tokenString, "Bearer ")
	fmt.Println("token", tokenString)

	token, _ := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method")
		}
		return []byte("secret"), nil
	})
	// var result database.User
	var name, email string
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		name = claims["name"].(string)
		email = claims["email"].(string)
	}
	fmt.Println(name, email)
	p := database.Finddb(cl, email)
	fmt.Println(p)
	if p.Name != "" {
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
}

//experience updates the experience
func experience(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	tokenString := r.Header.Get("Authorization")

	tokenString = strings.TrimPrefix(tokenString, "Bearer ")
	fmt.Println("token", tokenString)

	token, _ := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method")
		}
		return []byte("secret"), nil
	})
	// var result database.User
	var name, email string
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		name = claims["name"].(string)
		email = claims["email"].(string)
	}
	fmt.Println(name, email)
	p := database.Finddb(cl, email)
	fmt.Println(p)
	if p.Name != "" {
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
}

//Myprofile displays the profile of the user
func Myprofile(w http.ResponseWriter,r *http.Request){
	params:=mux.Vars(r)
	user:=params["name"]
	w.Header().Set("Content-Type", "application/json")
	tokenString := r.Header.Get("Authorization")

	tokenString = strings.TrimPrefix(tokenString, "Bearer ")
	fmt.Println("token", tokenString)

	token, _ := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method")
		}
		return []byte("secret"), nil
	})
	// var result database.User
	var name, email string
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		name = claims["name"].(string)
		email = claims["email"].(string)
	}
	fmt.Print(name)
	p := database.Finddb(cl, email)
	if p.Name==user{
		pro:=database.Findprofile(cl1,email)
		json.NewEncoder(w).Encode(pro)
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte(`{"success": "created"}`))
		return
	}
	w.WriteHeader(http.StatusNotFound)
	w.Write([]byte(`{"error": "not token not matched"}`))

}

func dashboard(w http.ResponseWriter,r *http.Request){
	params:=mux.Vars(r)
	user:=params["name"]
	w.Header().Set("Content-Type", "application/json")
	tokenString := r.Header.Get("Authorization")

	tokenString = strings.TrimPrefix(tokenString, "Bearer ")
	fmt.Println("token", tokenString)

	token, _ := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method")
		}
		return []byte("secret"), nil
	})
	// var result database.User
	var name, email string
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		name = claims["name"].(string)
		email = claims["email"].(string)
	}
	fmt.Print(name)
	p := database.Finddb(cl, email)
	if p.Name==user{
		pro:=database.Findprofile(cl1,email)
		json.NewEncoder(w).Encode(pro.Edu)
		json.NewEncoder(w).Encode(pro.Exp)
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte(`{"success": "Data extracted"}`))
		return
	}
	w.WriteHeader(http.StatusNotFound)
	w.Write([]byte(`{"error": "not token not matched"}`))

}

func WritePost(w http.ResponseWriter,r *http.Request){

	w.Header().Set("Content-Type", "application/json")
	tokenString := r.Header.Get("Authorization")

	tokenString = strings.TrimPrefix(tokenString, "Bearer ")
	fmt.Println("token", tokenString)

	token, _ := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method")
		}
		return []byte("secret"), nil
	})
	// var result database.User
	var name, email string
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		name = claims["name"].(string)
		email = claims["email"].(string)
	}
	var set []database.Post
	set=database.FindPost(cl2)
    json.NewEncoder(w).Encode(set)
	if r.Method=="POST"{
		type p struct{
			Post string
		}
		var Post p
		body, _ := ioutil.ReadAll(r.Body)
		err := json.Unmarshal(body, &Post)
		var data string
		data=Post.Post
		fmt.Println(err)
        if err==nil{
          ok:=database.InsertPost(cl2,name,email,data)
          if ok{
          	w.WriteHeader(http.StatusCreated)
				w.Write([]byte(`{"success": "Data extracted"}`))
			}
		}
		http.Redirect(w,r,"/Dev/Post",http.StatusOK)
	}
}