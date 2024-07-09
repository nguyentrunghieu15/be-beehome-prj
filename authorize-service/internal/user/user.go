package user

import (
	"context"
	"net/http"

	"github.com/cerbos/cerbos-sdk-go/cerbos"
	"github.com/labstack/echo/v4"
	"github.com/nguyentrunghieu15/be-beehome-prj/authorize-service/internal/cerbosx"
	"github.com/nguyentrunghieu15/be-beehome-prj/authorize-service/internal/coverter"
	"github.com/nguyentrunghieu15/be-beehome-prj/authorize-service/internal/model"
	"github.com/nguyentrunghieu15/be-beehome-prj/internal/mongox"
)

var userRepository = mongox.Repository[model.User]{
	Client:     mongox.DefaultClient,
	Collection: "user",
}

func ListUsers(c echo.Context) error {
	userId := c.Get("user_id")
	if userId == nil {
		return c.JSON(http.StatusNonAuthoritativeInfo, nil)
	}

	user, err := userRepository.FindOneByAtribute("user_id", userId)
	if err != nil {
		return err
	}

	p, err := coverter.ToPrincipal(coverter.MongoUserToPrincipalInfor(*user))
	if err != nil {
		return err
	}
	hasPermission, err := cerbosx.DefaultClient.CanActive(
		context.Background(),
		p,
		cerbos.NewResource("user", "1"),
		cerbosx.READ,
	)

	if err != nil || !hasPermission {
		return c.JSON(http.StatusNonAuthoritativeInfo, nil)
	}

	return c.JSON(http.StatusOK, nil)
}

func CreateUser(c echo.Context) error {
	// any one can create account
	return c.JSON(http.StatusOK, nil)
}

func BlockUser(c echo.Context) error {
	id := c.Param("id")
	userId := c.Get("user_id")
	if userId == nil {
		return c.JSON(http.StatusNonAuthoritativeInfo, nil)
	}

	user, err := userRepository.FindOneByAtribute("user_id", userId)
	if err != nil {
		return err
	}

	blockedUser, err := userRepository.FindOneByAtribute("user_id", id)
	if err != nil {
		return err
	}

	p, err := coverter.ToPrincipal(coverter.MongoUserToPrincipalInfor(*user))
	if err != nil {
		return err
	}

	r, err := coverter.ToResource(blockedUser)
	if err != nil {
		return err
	}

	hasPermission, err := cerbosx.DefaultClient.CanActive(
		context.Background(),
		p,
		r,
		cerbosx.UPDATE,
	)

	if err != nil || !hasPermission {
		return c.JSON(http.StatusNonAuthoritativeInfo, nil)
	}
	return c.JSON(http.StatusOK, nil)
}

func GetUser(c echo.Context) error {

	return c.JSON(http.StatusOK, nil)
}

func DeleteUser(c echo.Context) error {
	id := c.Param("id")
	userId := c.Get("user_id")
	if userId == nil {
		return c.JSON(http.StatusNonAuthoritativeInfo, nil)
	}

	user, err := userRepository.FindOneByAtribute("user_id", userId)
	if err != nil {
		return err
	}

	userResource, err := userRepository.FindOneByAtribute("user_id", id)
	if err != nil {
		return err
	}

	p, err := coverter.ToPrincipal(coverter.MongoUserToPrincipalInfor(*user))
	if err != nil {
		return err
	}

	r, err := coverter.ToResource(userResource)
	if err != nil {
		return err
	}

	hasPermission, err := cerbosx.DefaultClient.CanActive(
		context.Background(),
		p,
		r,
		cerbosx.DELETE,
	)

	if err != nil || !hasPermission {
		return c.JSON(http.StatusNonAuthoritativeInfo, nil)
	}
	return c.JSON(http.StatusOK, nil)
}

func UpdateUser(c echo.Context) error {
	id := c.Param("id")
	userId := c.Get("user_id")
	if userId == nil {
		return c.JSON(http.StatusNonAuthoritativeInfo, nil)
	}

	user, err := userRepository.FindOneByAtribute("user_id", userId)
	if err != nil {
		return err
	}

	userResource, err := userRepository.FindOneByAtribute("user_id", id)
	if err != nil {
		return err
	}

	p, err := coverter.ToPrincipal(coverter.MongoUserToPrincipalInfor(*user))
	if err != nil {
		return err
	}

	r, err := coverter.ToResource(userResource)
	if err != nil {
		return err
	}

	hasPermission, err := cerbosx.DefaultClient.CanActive(
		context.Background(),
		p,
		r,
		cerbosx.UPDATE,
	)

	if err != nil || !hasPermission {
		return c.JSON(http.StatusNonAuthoritativeInfo, nil)
	}
	return c.JSON(http.StatusOK, nil)
}
