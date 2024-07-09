package pro

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

var groupServiceRepository = mongox.Repository[model.GroupService]{
	Client:     mongox.DefaultClient,
	Collection: "group_service",
}

func ListGroupServices(c echo.Context) error {
	// all user can view
	return c.JSON(http.StatusOK, nil)
}

func CreateGroupService(c echo.Context) error {
	// only admin
	req := new(CreateServiceRequest)
	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	userReqId := c.Get("user_id")
	if userReqId == nil {
		return c.JSON(http.StatusBadRequest, nil)
	}

	userReq, err := userRepository.FindOneByAtribute("user_id", userReqId)
	if err != nil {
		return err
	}

	p, err := coverter.ToPrincipal(coverter.MongoUserToPrincipalInfor(*userReq))
	if err != nil {
		return err
	}

	hasPermission, err := cerbosx.DefaultClient.CanActive(
		context.Background(),
		p,
		cerbos.NewResource("service", "1"),
		cerbosx.CREATE,
	)
	if err != nil || !hasPermission {
		return c.JSON(http.StatusNonAuthoritativeInfo, nil)
	}

	return c.JSON(http.StatusOK, nil)
}

func FulltextSearchGroupServices(c echo.Context) error {
	// all user can find
	return c.JSON(http.StatusOK, nil)
}

func GetGroupService(c echo.Context) error {
	// all user can get
	return c.JSON(http.StatusOK, nil)
}

func DeleteGroupService(c echo.Context) error {
	id := c.Param("id")
	userReqId := c.Get("user_id")
	if userReqId == nil {
		return c.JSON(http.StatusBadRequest, nil)
	}

	userReq, err := userRepository.FindOneByAtribute("user_id", userReqId)
	if err != nil {
		return err
	}

	grpservice, err := groupServiceRepository.FindOneByAtribute("group-service_id", id)
	if err != nil {
		return err
	}

	p, err := coverter.ToPrincipal(coverter.MongoUserToPrincipalInfor(*userReq))
	if err != nil {
		return err
	}

	r, err := coverter.ToResource(grpservice)
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

func UpdateGroupService(c echo.Context) error {
	req := new(UpdateGroupServiceRequest)
	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	userReqId := c.Get("user_id")
	if userReqId == nil {
		return c.JSON(http.StatusBadRequest, nil)
	}

	userReq, err := userRepository.FindOneByAtribute("user_id", userReqId)
	if err != nil {
		return err
	}

	grpservice, err := groupServiceRepository.FindOneByAtribute("group_service_id", req.Id)
	if err != nil {
		return err
	}

	p, err := coverter.ToPrincipal(coverter.MongoUserToPrincipalInfor(*userReq))
	if err != nil {
		return err
	}

	r, err := coverter.ToResource(grpservice)
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

type CreateGroupServiceRequest struct {
	Name   string `json:"name"   form:"name"   query:"name"`
	Detail string `json:"detail" form:"detail" query:"detail"`
}

type UpdateGroupServiceRequest struct {
	Id     string `json:"id"     form:"id"     query:"id"`
	Name   string `json:"name"   form:"name"   query:"name"`
	Detail string `json:"detail" form:"detail" query:"detail"`
}
