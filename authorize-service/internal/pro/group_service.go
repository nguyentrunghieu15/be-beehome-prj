package pro

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func ListGroupServices(c echo.Context) error {
	return c.JSON(http.StatusOK, "ListGroupServices response")
}

func CreateGroupService(c echo.Context) error {
	return c.JSON(http.StatusOK, "CreateGroupService response")
}

func FulltextSearchGroupServices(c echo.Context) error {
	return c.JSON(http.StatusOK, "FulltextSearchGroupServices response")
}

func GetGroupService(c echo.Context) error {
	return c.JSON(http.StatusOK, "GetGroupService response")
}

func DeleteGroupService(c echo.Context) error {
	return c.JSON(http.StatusOK, "DeleteGroupService response")
}

func UpdateGroupService(c echo.Context) error {
	return c.JSON(http.StatusOK, "UpdateGroupService response")
}
