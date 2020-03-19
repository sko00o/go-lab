package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/graphql-go/graphql"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/gommon/log"
	"github.com/mitchellh/mapstructure"
)

// User ...
type User struct {
	ID       string `json:"id,omitempty"`
	Username string `json:"username,omitempty"`
	Password string `json:"password,omitempty"`
}

// Blog ...
type Blog struct {
	ID        string `json:"id,omitempty"`
	Title     string `json:"title,omitempty"`
	Content   string `json:"content,omitempty"`
	Author    string `json:"author,omitempty"`
	Pageviews int32  `json:"pageviews,omitempty"`
}

var jwtSecret = []byte("#asdVQSD")

var accountsMock = []User{
	User{
		ID:       "1",
		Username: "asd",
		Password: "asd",
	},
	User{
		ID:       "2",
		Username: "qwe",
		Password: "qwe",
	},
}

var blogsMock = []Blog{
	Blog{
		ID:        "1",
		Author:    "asd",
		Title:     "Sample Article",
		Content:   "This is a sample article weitten by asd",
		Pageviews: 1000,
	},
}

// ValidateJWT ...
func ValidateJWT(t string) (interface{}, error) {
	if t == "" {
		return nil, errors.New("auth error")
	}
	token, _ := jwt.Parse(t, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("There was an error")
		}
		return jwtSecret, nil
	})
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		var decodedToken interface{}
		mapstructure.Decode(claims, &decodedToken)
		return decodedToken, nil
	}
	return nil, errors.New("Invalid authorization token")
}

// CreateTokenEndpoint ...
func CreateTokenEndpoint(w http.ResponseWriter, r *http.Request) {
	var user User
	_ = json.NewDecoder(r.Body).Decode(&user)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": user.Username,
		"password": user.Password,
	})
	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		log.Error(err)
	}
	w.Header().Set("content-type", "application/json")
	w.Write([]byte(`
	{
		"toekn": "` + tokenString + `"
	}
	`))
}

type key string

const (
	keyPrincipalID key = "token"
)

func main() {
	rootQuery := graphql.NewObject(graphql.ObjectConfig{
		Name: "Query",
		Fields: graphql.Fields{
			"account": &graphql.Field{
				Type: accountType,
				Resolve: func(params graphql.ResolveParams) (interface{}, error) {
					account, err := ValidateJWT(params.Context.Value("token").(string))
					if err != nil {
						return nil, err
					}
					for _, act := range accountsMock {
						if act.Username == account.(User).Username {
							return act, nil
						}
					}
					return &User{}, nil
				},
			},
			"blogs": &graphql.Field{
				Type: graphql.NewList(blogType),
				Resolve: func(params graphql.ResolveParams) (interface{}, error) {
					return blogsMock, nil
				},
			},
		},
	})

	schema, _ := graphql.NewSchema(graphql.SchemaConfig{
		Query: rootQuery,
	})

	http.HandleFunc("/graphql", func(w http.ResponseWriter, r *http.Request) {
		res := graphql.Do(graphql.Params{
			Schema:        schema,
			RequestString: r.URL.Query().Get("query"),
			Context:       context.WithValue(context.Background(), keyPrincipalID, r.URL.Query().Get("token")),
		})
		json.NewEncoder(w).Encode(res)
	})
	fmt.Println("Starting the application at :12345")
	http.HandleFunc("/login", CreateTokenEndpoint)
	http.ListenAndServe(":12345", nil)
}

var accountType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "Account",
		Fields: graphql.Fields{
			"id": &graphql.Field{
				Type: graphql.String,
			},
			"username": &graphql.Field{
				Type: graphql.String,
			},
			"password": &graphql.Field{
				Type: graphql.String,
			},
		},
	})

var blogType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "Blog",
		Fields: graphql.Fields{
			"id": &graphql.Field{
				Type: graphql.String,
			},
			"title": &graphql.Field{
				Type: graphql.String,
			},
			"content": &graphql.Field{
				Type: graphql.String,
			},
			"author": &graphql.Field{
				Type: graphql.String,
			},
			"pageviews": &graphql.Field{
				Type: graphql.Int,
				Resolve: func(params graphql.ResolveParams) (interface{}, error) {
					_, err := ValidateJWT(params.Context.Value("token").(string))
					if err != nil {
						return nil, err
					}
					return params.Source.(Blog).Pageviews, nil
				},
			},
		},
	})
