package user

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func GetProfile(c echo.Context) error {
	// Your logic here
	profile := UserInfor{}
	return c.JSON(http.StatusOK, profile)
}

func DeactiveAccount(c echo.Context) error {
	// Your logic here
	return c.JSON(http.StatusOK, map[string]interface{}{})
}

func AddCard(c echo.Context) error {
	req := new(AddCardRequest)
	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	// Your logic here
	return c.JSON(http.StatusOK, map[string]interface{}{})
}

func ChangeEmail(c echo.Context) error {
	req := new(ChangeEmailRequest)
	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	// Your logic here
	return c.JSON(http.StatusOK, UserInfor{})
}

func ChangeName(c echo.Context) error {
	req := new(ChangeNameRequest)
	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	// Your logic here
	return c.JSON(http.StatusOK, UserInfor{})
}
