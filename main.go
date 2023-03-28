package main

import "github.com/labstack/echo/v4"

func main() {
	echo := echo.New()

	echo.Logger.Fatal(echo.Start(":8080"))
}
