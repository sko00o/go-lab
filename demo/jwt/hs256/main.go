package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	jwtmiddleware "github.com/auth0/go-jwt-middleware"
	jwt "github.com/dgrijalva/jwt-go"
	macaron "gopkg.in/macaron.v1"
)

const (
	SecretKey     = "helloworld"
	RightUsername = "asd"
	RightPassword = "123"
)

type LoginForm struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type Token struct {
	Token string `json:"token"`
}

func main() {
	// StartServer()
	str := "123123123"
	for _, v := range str {
		if v == '1' {
			fmt.Print("Y")
		}
		fmt.Printf("%v %T\n", v, v)
	}
}

func StartServer() {
	m := macaron.Classic()

	jwtMiddleware := jwtmiddleware.New(jwtmiddleware.Options{
		ValidationKeyGetter: func(token *jwt.Token) (interface{}, error) {
			return []byte(SecretKey), nil
		},
	})

	m.Post("/login", LoginHandler)
	m.Get("/resource", jwtMiddleware.CheckJWT, ProtectedHandler)
	m.Run()
}

func LoginHandler(w http.ResponseWriter, r *http.Request) (int, string) {

	var user LoginForm

	// Decode json to struct
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		return http.StatusForbidden, "Error in request"
	}

	// Check login
	if user.Username != RightUsername {
		return http.StatusForbidden, "Wrong username"
	} else if user.Password != RightPassword {
		return http.StatusForbidden, "Wrong password"
	}

	// Make token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": user.Username,
	})

	// Generate tokenString
	tokenString, err := token.SignedString([]byte(SecretKey))
	if err != nil {
		return http.StatusInternalServerError, "Error while signing the token"
	}

	return http.StatusOK, "Token: " + tokenString
}

func ProtectedHandler(w http.ResponseWriter, r *http.Request) (int, string) {
	return http.StatusOK, "Gained access to protected resource"
}

/*
reference:
https://github.com/auth0/go-jwt-middleware/blob/master/examples/martini-example/main.go?1529459038439
https://blog.csdn.net/wangshubo1989/article/details/74529333
*/
