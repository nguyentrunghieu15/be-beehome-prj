package pro

import (
	"context"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/nguyentrunghieu15/be-beehome-prj/authorize-service/internal/cerbosx"
	"github.com/nguyentrunghieu15/be-beehome-prj/authorize-service/internal/coverter"
	"github.com/nguyentrunghieu15/be-beehome-prj/authorize-service/internal/model"
	"github.com/nguyentrunghieu15/be-beehome-prj/internal/mongox"
)

var hireRepository mongox.Repository[model.Hire]

func FindAllHire(c echo.Context) error {
	// alway pass by check jwt
	return c.JSON(http.StatusOK, nil)
}

func CreateHire(c echo.Context) error {
	// Anyone can create hire request
	return c.JSON(http.StatusOK, nil)
}

func DeleteHire(c echo.Context) error {
	hireID := c.Param("hire_id")
	userReqId := c.Get("user_id")
	if userReqId == nil {
		return c.JSON(http.StatusBadRequest, nil)
	}

	userReq, err := userRepository.FindOneByAtribute("user_id", userReqId)
	if err != nil {
		return err
	}

	hire, err := hireRepository.FindOneByAtribute("hire_id", hireID)
	if err != nil {
		return err
	}

	p, err := coverter.ToPrincipal(*userReq)
	if err != nil {
		return err
	}

	r, err := coverter.ToResource(*hire)
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

func UpdateStatusHire(c echo.Context) error {
	reqId := c.Param("hireId")

	userReqId := c.Get("user_id")
	if userReqId == nil {
		return c.JSON(http.StatusBadRequest, nil)
	}

	userReq, err := userRepository.FindOneByAtribute("user_id", userReqId)
	if err != nil {
		return err
	}

	hire, err := hireRepository.FindOneByAtribute("hire_id", reqId)
	if err != nil {
		return err
	}

	p, err := coverter.ToPrincipal(*userReq)
	if err != nil {
		return err
	}

	r, err := coverter.ToResource(*hire)
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

type CreateHireRequest struct {
	ProviderID      string `json:"provider_id"       form:"provider_id"       query:"provider_id"`
	ServiceID       string `json:"service_id"        form:"service_id"        query:"service_id"`
	WorkTimeFrom    string `json:"work_time_from"    form:"work_time_from"    query:"work_time_from"`
	WorkTimeTo      string `json:"work_time_to"      form:"work_time_to"      query:"work_time_to"`
	Status          string `json:"status"            form:"status"            query:"status"`
	PaymentMethodID string `json:"payment_method_id" form:"payment_method_id" query:"payment_method_id"`
	Issue           string `json:"issue"             form:"issue"             query:"issue"`
}

type UpdateStatusHireRequest struct {
	HireID    string `json:"hire_id"    form:"hire_id"    query:"hire_id"`
	NewStatus string `json:"new_status" form:"new_status" query:"new_status"`
}
