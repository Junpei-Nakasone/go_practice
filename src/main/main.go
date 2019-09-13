package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

type JwtClaims struct {
	Name string `json:"name"`
	jwt.StandardClaims
}

func getMember(c echo.Context) error {
	memberName := c.QueryParam("name")
	memberAge := c.QueryParam("age")

	dataType := c.Param("data")

	if dataType == "string" {
		return c.String(http.StatusOK, fmt.Sprintf("The member name is %s\nand %s years old", memberName, memberAge))
	}

	if dataType == "json" {
		return c.JSON(http.StatusOK, map[string]string{
			"name": memberName,
			"type": memberAge,
		})
	}
	return c.JSON(http.StatusBadRequest, map[string]string{
		"error": "you need to specify the member's name/age.",
	})

}

func mainJwt(c echo.Context) error {
	user := c.Get("user")
	token := user.(*jwt.Token)

	claims := token.Claims.(jwt.MapClaims)

	log.Println("User Name: ", claims["name"], "User ID: ", claims["jti"])

	return c.String(http.StatusOK, "you are on the top secret jwt page!")
}

func login(c echo.Context) error {
	username := c.QueryParam("name")
	password := c.QueryParam("password")

	if username == "jack" && password == "1234" {
		cookie := &http.Cookie{}

		// this is the same
		//cookie := new(http.Cookie)

		cookie.Name = "sessionID"
		cookie.Value = "some_string"
		cookie.Expires = time.Now().Add(48 * time.Hour)

		c.SetCookie(cookie)

		// create JWT token
		token, err := createJwtToken()
		if err != nil {
			log.Println("Error Occured.", err)
			return c.String(http.StatusInternalServerError, "Something went wrong")
		}
		jwtCookie := &http.Cookie{}

		jwtCookie.Name = "JWTCookie"
		jwtCookie.Value = token
		jwtCookie.Expires = time.Now().Add(24 * time.Hour)

		c.SetCookie(jwtCookie)

		return c.JSON(http.StatusOK, map[string]string{
			"message": "You logged in.",
			"token":   token,
		})
	}
	return c.String(http.StatusUnauthorized, "Your name or password is invalid")
}

func createJwtToken() (string, error) {
	claims := JwtClaims{
		"jack",
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
	fmt.Println("WELCOME To SERVER")

	e := echo.New()

	jwtGroup := e.Group("/jwt")

	jwtGroup.Use(middleware.JWTWithConfig(middleware.JWTConfig{
		SigningMethod: "HS512",
		SigningKey:    []byte("mySecret"),
		TokenLookup:   "Cookie:JWTCookie",
	}))

	jwtGroup.GET("/main", mainJwt)

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "HELLO FROM WEBSite")
	})
	e.GET("/member/:data", getMember)
	e.GET("login", login)

	e.Start(":8000")
}
