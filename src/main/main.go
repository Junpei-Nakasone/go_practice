package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

type JwtClaims struct {
	Name string `json:"name"`
	jwt.StandardClaims
}

func hello(c echo.Context) error {
	return c.String(http.StatusOK, "Hello world!")
}

func mainJwt(c echo.Context) error {
	user := c.Get("user")
	token := user.(*jwt.Token)

	claims := token.Claims.(jwt.MapClaims)

	log.Println("User Name: ", claims["name"], "User ID: ", claims["jti"])

	return c.String(http.StatusOK, "you are on the top secret jwt page!")
}

func login(c echo.Context) error {
	username := c.QueryParam("username")
	password := c.QueryParam("password")

	// check username and password against DB after hashing the password
	if username == "nakasone" && password == "password" {

		// create jwt token
		token, err := createJwtToken()
		if err != nil {
			log.Println("Error Creating JWT token", err)
			return c.String(http.StatusInternalServerError, "something went wrong")
		}

		JwtCookie := &http.Cookie{}

		JwtCookie.Name = "JWTCookie"
		JwtCookie.Value = token
		JwtCookie.Expires = time.Now().Add(48 * time.Hour)

		c.SetCookie(JwtCookie)

		return c.JSON(http.StatusOK, map[string]string{
			"message": "You were logged in!",
			"token":   token,
		})
	}

	return c.String(http.StatusUnauthorized, "Your username or password were wrong")
}

func createJwtToken() (string, error) {
	claims := JwtClaims{
		"nakasone",
		jwt.StandardClaims{
			Id:        "main_user_id",
			ExpiresAt: time.Now().Add(24 * time.Hour).Unix(),
		},
	}

	rawToken := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)

	token, err := rawToken.SignedString([]byte("mySecret"))
	if err != nil {
		return "", err
	}

	return token, nil
}

func main() {
	fmt.Println("Welcome to the server")

	e := echo.New()

	jwtGroup := e.Group("/jwt")

	jwtGroup.Use(middleware.JWTWithConfig(middleware.JWTConfig{
		SigningMethod: "HS512",
		SigningKey:    []byte("mySecret"),
		TokenLookup:   "cookie:JWTCookie",
	}))

	jwtGroup.GET("/main", mainJwt)

	e.GET("/login", login)
	e.GET("/", hello)

	e.Start(":8000")

}
