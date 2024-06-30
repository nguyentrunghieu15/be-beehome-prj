package pro

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func ListServices(c echo.Context) error {
	return c.JSON(http.StatusOK, "ListServices response")
}

func CreateService(c echo.Context) error {
	return c.JSON(http.StatusOK, "CreateService response")
}

func FulltextSearchServices(c echo.Context) error {
	return c.JSON(http.StatusOK, "FulltextSearchServices response")
}

func GetService(c echo.Context) error {
	return c.JSON(http.StatusOK, "GetService response")
}

func DeleteService(c echo.Context) error {
	return c.JSON(http.StatusOK, "DeleteService response")
}

func UpdateService(c echo.Context) error {
	return c.JSON(http.StatusOK, "UpdateService response")
}
