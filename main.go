package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/vashish1/Dev/database"

	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/mongo"
)

var secret = os.Getenv("BlockKey")
var Port = os.Getenv("PORT")

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/Dev", index)
	r.HandleFunc("/Dev/login", login).Methods("POST")
	r.HandleFunc("/Dev/signup", signup).Methods("POST")
	r.HandleFunc("/Dev/profile/add/Education", addEducation).Methods("POST")
	r.HandleFunc("/Dev/profile/add/Experience", addExperience).Methods("POST")
	r.HandleFunc("/Dev/profile/update", updateProfile).Methods("POST")
	r.HandleFunc("/Dev/Profile/{id}", Myprofile).Methods("GET")
	r.HandleFunc("/Dev/Dashboard", dashboard).Methods("GET")
	r.HandleFunc("/Dev/Developers", developers).Methods("GET")
	r.HandleFunc("/Dev/Post", writePost).Methods("GET","POST")
	r.HandleFunc("/Dev/Post/comment/{id}", comment).Methods("GET")
	r.HandleFunc("/like/{id}", like).Methods("GET")
	r.HandleFunc("/dislike/{id}", dislike).Methods("GET")
	http.Handle("/", r)
	http.ListenAndServe(":"+Port, nil)
}

var cl, cl1, cl2 *mongo.Collection
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
	var result LoginResponse
	var user logn
	body, _ := ioutil.ReadAll(r.Body)
	err := json.Unmarshal(body, &user)

	//Handles the error whhile parsing the Request JSON
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"error": "Body not parsed"}`))
		return
	}

	//If user Exists in DB, Creating Token for response else returning error
	ok := database.Findfromuserdb(cl, user.Email, user.Password)
	if ok {
		u := database.Finddb(cl, user.Email)
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"name": u.Name,
			"uid":  u.UUID,
		})

		tokenString, err := token.SignedString([]byte(secret))

		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(`{"error": "Error While Creating Token"}`))
			return
		}
		result.Token = tokenString
		result.Success = true
		tkn := database.UpdateToken(cl, u.Email, tokenString)
		if tkn {
			json.NewEncoder(w).Encode(result)
			w.WriteHeader(200)
			return
		}
	}
	w.WriteHeader(400)
	w.Write([]byte(`{"error": "Login denied"}`))
}

//signup handles the signup of a new user in the DevConnect
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
	} else {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(`{"error": "Passwords do not match"}`))
	}
}

//updateProfile updates the profile of the user
func updateProfile(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	tokenString := r.Header.Get("Authorization")

	tokenString = strings.TrimPrefix(tokenString, "Bearer ")
	fmt.Println("token", tokenString)
	token, _ := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method")
		}
		return []byte(secret), nil
	})
	// var result database.User
	var _, uid string
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		_ = claims["name"].(string)
		uid = claims["uid"].(string)
	}

	p := database.Finddb(cl, uid)

	if p.UUID != "" {
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

			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"success": "updated"}`))
			return
		}
	}
	w.WriteHeader(401)
	w.Write([]byte(`{"error": "User Unauthorized"}`))
}

//addEducation updates the education
func addEducation(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	tokenString := r.Header.Get("Authorization")

	tokenString = strings.TrimPrefix(tokenString, "Bearer ")
	fmt.Println("token", tokenString)

	token, _ := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method")
		}
		return []byte(secret), nil
	})
	var uid string
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		_ = claims["name"].(string)
		uid = claims["uid"].(string)
	}
	p := database.Finddb(cl, uid)
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
			w.Write([]byte(`{"success": "updated"}`))
			return
		}
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(`{"error": "could not update"}`))
		return
	}
	w.WriteHeader(401)
	w.Write([]byte(`{"error": "User Unauthorized"}`))
}

//addExperience updates the experience
func addExperience(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	tokenString := r.Header.Get("Authorization")

	tokenString = strings.TrimPrefix(tokenString, "Bearer ")
	fmt.Println("token", tokenString)

	token, _ := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method")
		}
		return []byte(secret), nil
	})
	var uid string
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		_ = claims["name"].(string)
		uid = claims["uid"].(string)
	}
	p := database.Finddb(cl, uid)
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
			w.Write([]byte(`{"success": "updated"}`))
			return
		}
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(`{"error": "could not update"}`))
		return
	}
	w.WriteHeader(401)
	w.Write([]byte(`{"error": "User Unauthorized"}`))
}

//Myprofile displays the profile of the user
func Myprofile(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["id"]
	w.Header().Set("Content-Type", "application/json")
	tokenString := r.Header.Get("Authorization")

	tokenString = strings.TrimPrefix(tokenString, "Bearer ")
	fmt.Println("token", tokenString)

	token, _ := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method")
		}
		return []byte(secret), nil
	})
	// var result database.User
	var _, uid string
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		_ = claims["name"].(string)
		uid = claims["uid"].(string)
	}
	p := database.Finddb(cl, id)
	if p.UUID == id {
		pro := database.Findprofile(cl1, uid)
		json.NewEncoder(w).Encode(pro)
		w.WriteHeader(200)
		return
	}
	w.WriteHeader(401)
	w.Write([]byte(`{"error": "User Unauthorized"}`))

}

func dashboard(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	tokenString := r.Header.Get("Authorization")

	tokenString = strings.TrimPrefix(tokenString, "Bearer ")
	fmt.Println("token", tokenString)
	var result Dasboard
	token, _ := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method")
		}
		return []byte(secret), nil
	})
	// var result database.User
	var uid string
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		_ = claims["name"].(string)
		uid = claims["uid"].(string)
	}
	p := database.Finddb(cl, uid)
	if p.UUID != "" {
		pro := database.Findprofile(cl1, uid)
		result.Education = pro.Edu
		result.Experience = pro.Exp
		w.WriteHeader(http.StatusOK)
		return
	}
	w.WriteHeader(401)
	w.Write([]byte(`{"error": "User Unauthorised"}`))

}

func writePost(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")
	tokenString := r.Header.Get("Authorization")

	tokenString = strings.TrimPrefix(tokenString, "Bearer ")
	fmt.Println("token", tokenString)

	token, _ := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method")
		}
		return []byte(secret), nil
	})
	var uid string
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		_ = claims["name"].(string)
		uid = claims["uid"].(string)
	}
	user:=database.Finddb(cl,uid)
    if user.UUID!=""{
		if r.Method == "POST" {

			var postdata p
			body, _ := ioutil.ReadAll(r.Body)
			err := json.Unmarshal(body, &postdata)
			var Post database.Post
			Post.Id = database.PostId()
			Post.Email = user.Email
			Post.UserName = user.Name
			Post.Text = postdata.Data
			var ok, okk bool
			if err == nil {
				ok = database.InsertPost(cl2, Post)
				okk = database.UpdateUserPostId(cl, user.Email, Post.Id)
				if !ok || !okk {
					w.WriteHeader(http.StatusCreated)
					w.Write([]byte(`{"error": "post not created"}`))
					return
				}
					w.WriteHeader(http.StatusOK)
					w.Write([]byte(`{"success": "created"}`))
					return
				}
				http.Redirect(w, r, "/Dev/Post", http.StatusOK)
			}else{
			Total := database.FindPost(cl2)
			json.NewEncoder(w).Encode(Total)
			w.WriteHeader(200)
			return
		 }
	}
	
	w.WriteHeader(401)
	w.Write([]byte(`{"error": "User Unauthorised"}`))
}

func developers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	tokenString := r.Header.Get("Authorization")

	tokenString = strings.TrimPrefix(tokenString, "Bearer ")
	fmt.Println("token", tokenString)

	token, _ := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method")
		}
		return []byte(secret), nil
	})
	var uid string
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		_ = claims["name"].(string)
		uid = claims["uid"].(string)
	}
	p:=database.Finddb(cl,uid)
	if  p.UUID!=""{
		developer := database.FindDevelopers(cl1)
		json.NewEncoder(w).Encode(developer)
		w.WriteHeader(http.StatusOK)
		return
	}
    w.WriteHeader(401)
	w.Write([]byte(`{"error": "User Unauthorised"}`))
}

func comment(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	tokenString := r.Header.Get("Authorization")

	tokenString = strings.TrimPrefix(tokenString, "Bearer ")
	fmt.Println("token", tokenString)

	token, _ := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method")
		}
		return []byte(secret), nil
	})
	var uid string
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		_ = claims["name"].(string)
		uid = claims["uid"].(string)
	}
	user:=database.Finddb(cl,uid)
	if user.UUID!=""{
	params := mux.Vars(r)
	pro := params["id"]
	id, _ := strconv.Atoi(pro)
	var x p
	body, _ := ioutil.ReadAll(r.Body)
	err := json.Unmarshal(body, &x)
	if err == nil {
			ok := database.UpdateComments(cl2, id, x.Data)
			if ok {
				Total := database.FindComment(cl2, user.Email, id)
			    json.NewEncoder(w).Encode(Total)
			    w.WriteHeader(http.StatusOK)
			    return
			}
			w.WriteHeader(400)
			w.Write([]byte(`{"error": "Error while Commenting"}`))
			return	
		} 
			w.WriteHeader(400)
			w.Write([]byte(`{"error": "Body not parsed"}`))
			return
	}
    w.WriteHeader(401)
	w.Write([]byte(`{"error": "User Unauthorised"}`))
}

func like(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	tokenString := r.Header.Get("Authorization")

	tokenString = strings.TrimPrefix(tokenString, "Bearer ")
	fmt.Println("token", tokenString)

	token, _ := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method")
		}
		return []byte(secret), nil
	})
	var uid string
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		_ = claims["name"].(string)
		uid = claims["uid"].(string)
	}
	user:=database.Finddb(cl,uid)
	if user.UUID!=""{
	params := mux.Vars(r)
	pro := params["id"]
	id, _ := strconv.Atoi(pro)
	ok := database.UpdateLikes(cl2, user.Email, id)
	if ok {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"success": "updated"}`))
		return
	} 
		w.WriteHeader(400)
		w.Write([]byte(`{"error": "Error while Liking the post"}`))
		return
	}
	w.WriteHeader(401)
	w.Write([]byte(`{"error": "User Unauthorised"}`))
}

func dislike(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	tokenString := r.Header.Get("Authorization")

	tokenString = strings.TrimPrefix(tokenString, "Bearer ")
	token, _ := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method")
		}
		return []byte(secret), nil
	})
	var uid string
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		_ = claims["name"].(string)
		uid = claims["uid"].(string)
	}
	user:=database.Finddb(cl,uid)
	if user.UUID!=""{
	params := mux.Vars(r)
	pro := params["id"]
	id, _ := strconv.Atoi(pro)
	ok := database.UpdateDisLikes(cl2, user.Email, id)
	if ok {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"success": "updated"}`))
		return
	} 
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"error": "Error while disliking the post"}`))
	   return
}
	w.WriteHeader(401)
	w.Write([]byte(`{"error": "User Unauthorised"}`))
}