package user

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func GetProfile(c echo.Context) error {
	// Always by pass by jwt token
	return c.JSON(http.StatusOK, nil)
}

func DeactiveAccount(c echo.Context) error {
	// Always by pass by jwt token
	return c.JSON(http.StatusOK, map[string]interface{}{})
}

func AddCard(c echo.Context) error {
	// Always by pass by jwt token
	return c.JSON(http.StatusOK, map[string]interface{}{})
}

func ChangeEmail(c echo.Context) error {
	// Always by pass by jwt token
	return c.JSON(http.StatusOK, UserInfor{})
}

func ChangeName(c echo.Context) error {
	// Always by pass by jwt token
	return c.JSON(http.StatusOK, UserInfor{})
}
