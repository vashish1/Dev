package database

import (
	"crypto/sha1"
	"encoding/hex"

	"github.com/google/uuid"
)

//User ......
type User struct {
	UUID         string
    Name    string
	Email        string
	PasswordHash string
	Image        string
}

//Newuser .....
func Newuser(name string, email string, password string) User {

	Password := SHA256ofstring(password)
	U := User{UUID: GenerateUUID(), Name: name, Email: email, PasswordHash: Password, Image: ""}
	return U
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
