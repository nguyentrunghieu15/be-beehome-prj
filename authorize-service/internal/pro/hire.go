package pro

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func FindAllHire(c echo.Context) error {
	return c.JSON(http.StatusOK, "FindAllHire response")
}

func CreateHire(c echo.Context) error {
	return c.JSON(http.StatusOK, "CreateHire response")
}

func DeleteHire(c echo.Context) error {
	return c.JSON(http.StatusOK, "DeleteHire response")
}

func UpdateStatusHire(c echo.Context) error {
	return c.JSON(http.StatusOK, "UpdateStatusHire response")
}
