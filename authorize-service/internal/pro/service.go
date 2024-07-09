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

var serviceRepository = mongox.Repository[model.Service]{
	Client:     mongox.DefaultClient,
	Collection: "service",
}
var userRepository = mongox.Repository[model.User]{
	Client:     mongox.DefaultClient,
	Collection: "user",
}

func ListServices(c echo.Context) error {
	// all user can view service
	return c.JSON(http.StatusOK, nil)
}

func CreateService(c echo.Context) error {
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

func FulltextSearchServices(c echo.Context) error {
	// all user can find
	return c.JSON(http.StatusOK, nil)
}

func GetService(c echo.Context) error {
	//all user can find
	return c.JSON(http.StatusOK, nil)
}

func DeleteService(c echo.Context) error {
	id := c.Param("id")
	userReqId := c.Get("user_id")
	if userReqId == nil {
		return c.JSON(http.StatusBadRequest, nil)
	}

	userReq, err := userRepository.FindOneByAtribute("user_id", userReqId)
	if err != nil {
		return err
	}

	service, err := serviceRepository.FindOneByAtribute("service_id", id)
	if err != nil {
		return err
	}

	p, err := coverter.ToPrincipal(coverter.MongoUserToPrincipalInfor(*userReq))
	if err != nil {
		return err
	}

	r, err := coverter.ToResource(service)
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

func UpdateService(c echo.Context) error {
	req := new(UpdateServiceRequest)
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

	service, err := serviceRepository.FindOneByAtribute("service_id", req.Id)
	if err != nil {
		return err
	}

	p, err := coverter.ToPrincipal(coverter.MongoUserToPrincipalInfor(*userReq))
	if err != nil {
		return err
	}

	r, err := coverter.ToResource(service)
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

type CreateServiceRequest struct {
	Name           string `json:"name"             form:"name"             query:"name"`
	Detail         string `json:"detail"           form:"detail"           query:"detail"`
	GroupServiceID string `json:"group_service_id" form:"group_service_id" query:"group_service_id"`
}

type UpdateServiceRequest struct {
	Id     string `json:"id"     form:"id"     query:"id"`
	Name   string `json:"name"   form:"name"   query:"name"`
	Detail string `json:"detail" form:"detail" query:"detail"`
}
