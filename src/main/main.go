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

// Jwtトークンを発行する時に使うstruct
type JwtClaims struct {
	Name string `json:"name"`
	// jwt-goで定義されているstruct。７項目ある
	jwt.StandardClaims
}

func hello(c echo.Context) error {
	return c.String(http.StatusOK, "Hello world!")
}

// "jwt/main"にアクセスする時に呼び出される関数。先にログインしてJWTが発行されていないとエラーになる。
func mainJwt(c echo.Context) error {
	user := c.Get("user")
	token := user.(*jwt.Token)

	claims := token.Claims.(jwt.MapClaims)

	log.Println("User Name: ", claims["name"], "User ID: ", claims["jti"])

	return c.String(http.StatusOK, "you are on the top secret jwt page!")

	// todo: エラー処理追加する
}

// "localhost:8000/login"にアクセスすると呼び出される関数。
func login(c echo.Context) error {
	// URLの"username="の値をusernameに格納
	username := c.QueryParam("username")
	// URLの"password="の値をpasswordに格納
	password := c.QueryParam("password")

	// 有効なusernameとpasswordを指定するif文
	if username == "nakasone" && password == "password" {

		// jwtを発行する関数をtokenに格納
		token, err := createJwtToken()

		// エラー処理
		if err != nil {
			log.Println("Error Creating JWT token", err)
			return c.String(http.StatusInternalServerError, "something went wrong")
		}

		// http.CookieのstructをJwtCookieに代入
		JwtCookie := &http.Cookie{}

		JwtCookie.Name = "JWTCookie"

		// tokenは上記で記述されているようにcreateJatToken関数を格納している
		JwtCookie.Value = token

		// timeパッケージを使用してJwtCookieの期限を48時間に指定
		JwtCookie.Expires = time.Now().Add(48 * time.Hour)

		// CookieをHTTPレスポンスに加える
		c.SetCookie(JwtCookie)

		// messageとtokenをブラウザに表示
		return c.JSON(http.StatusOK, map[string]string{
			"message": "You were logged in!",
			"token":   token,
		})
	}

	// 上記の処理が行われなかった場合
	return c.String(http.StatusUnauthorized, "Your username or password were wrong")
}

// jwtを発行する関数
func createJwtToken() (string, error) {

	// JwtClaimsの内容を設定し、claimsに格納
	claims := JwtClaims{
		"nakasone",
		jwt.StandardClaims{
			Id:        "main_user_id",
			ExpiresAt: time.Now().Add(24 * time.Hour).Unix(),
		},
	}

	// rawTokenにNewWithClaims関数で暗号化メソッドと上記claimsを格納
	rawToken := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)

	// tokenに上記rowTokenにSingedStringを格納
	token, err := rawToken.SignedString([]byte("mySecret"))
	if err != nil {
		return "", err
	}

	return token, nil
}

func main() {
	fmt.Println("Welcome to the server")

	e := echo.New()

	// "localhost:8000/jwt"配下へのアクセスはjwtGroupが適用される
	jwtGroup := e.Group("/jwt")

	//jwtGroupはJWTWithConfigミドルウェアを使用する
	jwtGroup.Use(middleware.JWTWithConfig(middleware.JWTConfig{
		SigningMethod: "HS512",
		SigningKey:    []byte("mySecret"),
		// cookieを発行するよう指定
		TokenLookup: "cookie:JWTCookie",
	}))

	//"localhost:8000/jwt/main"にアクセスするとmainJwt関数が呼び出される
	jwtGroup.GET("/main", mainJwt)

	e.GET("/login", login)
	e.GET("/", hello)

	e.Start(":8000")

}
