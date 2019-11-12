package database

import (
	"crypto/sha1"
	"encoding/hex"

	"github.com/google/uuid"
)

//User ......
type User struct {
	UUID         string
	Name         string
	Email        string
	PasswordHash string
	Image        string
}
//profile ....
type Profile struct{
	Email string
	Status string
	Org string
	Website string
	Location string
	Skills []string
	Gitname string
	Bio string
    Social map[string]string
	Edu []Education
	Exp []Experience
}

//education .....
type Education struct{
	School string
	Degree string
	Field string
	From string
	To string
	Achievemets string
}

//experience ....
type Experience struct{
Title string
Org string
Location string
From string
To string
Description string
}

//Newuser .....
func Newuser(name string, email string, password string, img string) User {

	Password := SHA256ofstring(password)
	U := User{UUID: GenerateUUID(), Name: name, Email: email, PasswordHash: Password, Image: img}
	return U
}

func Newprofile(a string,b string,c string,d string,e string,f []string,g string ,h string,i map[string]string,j Education,k Experience)Profile{
	var pro Profile
	pro=Profile{
		Email:    a,
		Status:   b,
		Org:      c,
		Website:  d,
		Location: e,
		Skills:   f,
		Gitname:  g,
		Bio:      h,
		Social:   i,
		Edu:      nil,
		Exp:      nil,
	}
	return pro
}

func Neweducation(a,b,c,d,e,f string)Education{
	m:=Education{
		School:      a,
		Degree:      b,
		Field:       c,
		From:        d,
		To:          e,
		Achievemets: f,
	}
	return m
}


func Newexperience(a,b,c,d,e,f string)Experience{
	m:=Experience{
		Title:       a,
		Org:         b,
		Location:    c,
		From:        d,
		To:          e,
		Description: f,
	}
	return m
}

//SHA256ofstring is a function which takes a string a reurns its sha256 hashed form
func SHA256ofstring(p string) string {
	h := sha1.New()
	h.Write([]byte(p))
	hash := hex.EncodeToString(h.Sum(nil))
	return hash
}

//GenerateUUID generates a unique id for every user.
func GenerateUUID() string {

	sd := uuid.New()
	return (sd.String())

}


