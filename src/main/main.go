package main

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo"
)

func getCats(c echo.Context) error {
	catName := c.QueryParam("name")
	catType := c.QueryParam("type")

	dataType := c.Param("data")

	if dataType == "string" {
		return c.String(http.StatusOK, fmt.Sprintf("Your cat name is %s\nand  type is %s", catName, catType))
	}

	if dataType == "json" {
		return c.JSON(http.StatusOK, map[string]string{
			"name": catName,
			"type": catType,
		})
	}
	return c.JSON(http.StatusBadRequest, map[string]string{
		"error": "you need to specify your cat's name and type.",
	})

}

func main() {
	fmt.Println("WELCOME To SERVER")

	e := echo.New()

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "HELLO FROM WEBSite")
	})
	e.GET("/cats/:data", getCats)

	e.Start(":8000")
}
