package user

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
)

func ListUsers(c echo.Context) error {
	// Your logic here
	response := ListUsersResponse{}
	return c.JSON(http.StatusOK, response)
}

func CreateUser(c echo.Context) error {
	req := new(CreateUserRequest)
	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	// Your logic here
	return c.JSON(http.StatusOK, UserInfor{})
}

func BlockUser(c echo.Context) error {
	id := c.Param("id")
	fmt.Println(id)
	req := new(UserServiceBlockUserBody)
	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	// Your logic here
	return c.JSON(http.StatusOK, map[string]interface{}{})
}

func GetUser(c echo.Context) error {
	id := c.Param("id")
	fmt.Println(id)
	// Your logic here
	user := UserInfor{}
	return c.JSON(http.StatusOK, user)
}

func DeleteUser(c echo.Context) error {
	id := c.Param("id")
	fmt.Println(id)
	// Your logic here
	return c.JSON(http.StatusOK, map[string]interface{}{})
}

func UpdateUser(c echo.Context) error {
	id := c.Param("id")
	req := new(UserServiceUpdateUserBody)
	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	fmt.Println(id)
	// Your logic here
	return c.JSON(http.StatusOK, UserInfor{})
}
