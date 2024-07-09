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

var providerRepository = mongox.Repository[model.Provider]{
	Client:     mongox.DefaultClient,
	Collection: "provider",
}

var reviewRepository = mongox.Repository[model.Review]{
	Client:     mongox.DefaultClient,
	Collection: "review",
}

func FindPros(c echo.Context) error {
	//all user can find pros
	return c.JSON(http.StatusOK, nil)
}

func JoinAsProvider(c echo.Context) error {
	// all user can join as provider
	return c.JSON(http.StatusOK, nil)
}

func AddPaymentMethodPro(c echo.Context) error {
	// all pass by check jwt
	return c.JSON(http.StatusOK, nil)
}

func AddServicePro(c echo.Context) error {
	// all pass by check jwt
	return c.JSON(http.StatusOK, nil)
}

func AddSocialMediaPro(c echo.Context) error {
	// all pass by check jwt
	return c.JSON(http.StatusOK, nil)
}

func DeleteServicePro(c echo.Context) error {
	// all pass by check jwt
	return c.JSON(http.StatusOK, nil)
}

func DeleteSocialMediaPro(c echo.Context) error {
	// all pass by check jwt
	return c.JSON(http.StatusOK, nil)
}

func GetProviderProfile(c echo.Context) error {
	// all user can view
	return c.JSON(http.StatusOK, "GetProviderProfile nil")
}

func ReplyReviewPro(c echo.Context) error {
	req := new(ReplyReviewProRequest)
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

	review, err := reviewRepository.FindOneByAtribute("review_id", req.ReviewID)
	if err != nil {
		return err
	}

	p, err := coverter.ToPrincipal(coverter.MongoUserToPrincipalInfor(*userReq))
	if err != nil {
		return err
	}

	r, err := coverter.ToResource(review)
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

func ReviewPro(c echo.Context) error {
	req := new(ReviewProRequest)
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
		cerbos.NewResource("review", "1"),
		cerbosx.UPDATE,
	)
	if err != nil || !hasPermission {
		return c.JSON(http.StatusNonAuthoritativeInfo, nil)
	}

	return c.JSON(http.StatusOK, nil)
}

func SignUpPro(c echo.Context) error {
	// all user can sign up
	return c.JSON(http.StatusOK, nil)
}

func UpdateSocialMediaPro(c echo.Context) error {
	// bypass when check jwt
	return c.JSON(http.StatusOK, nil)
}

func FindProById(c echo.Context) error {
	// all user can find pro
	return c.JSON(http.StatusOK, nil)
}

func DeleteProById(c echo.Context) error {
	id := c.Param("id")
	userReqId := c.Get("user_id")
	if userReqId == nil {
		return c.JSON(http.StatusBadRequest, nil)
	}

	userReq, err := userRepository.FindOneByAtribute("user_id", userReqId)
	if err != nil {
		return err
	}

	provider, err := providerRepository.FindOneByAtribute("provider_id", id)
	if err != nil {
		return err
	}

	p, err := coverter.ToPrincipal(coverter.MongoUserToPrincipalInfor(*userReq))
	if err != nil {
		return err
	}

	r, err := coverter.ToResource(provider)
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

func UpdatePro(c echo.Context) error {
	req := new(UpdateProRequest)
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

	provider, err := providerRepository.FindOneByAtribute("provider_id", req.ID)
	if err != nil {
		return err
	}

	p, err := coverter.ToPrincipal(coverter.MongoUserToPrincipalInfor(*userReq))
	if err != nil {
		return err
	}

	r, err := coverter.ToResource(provider)
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

func GetAllReviewsOfProvider(c echo.Context) error {
	// all user can view
	return c.JSON(http.StatusOK, nil)
}

func GetAllServiceOfProvider(c echo.Context) error {
	// all user can view
	return c.JSON(http.StatusOK, nil)
}

// Define request structures
type JoinAsProviderRequest struct{}

type AddPaymentMethodProRequest struct {
	Name string `json:"name" form:"name" query:"name"`
}

type AddServiceProRequest struct {
	ServicesID []string `json:"services_id" form:"services_id" query:"services_id"`
}

type AddSocialMediaProRequest struct {
	Name string `json:"name" form:"name" query:"name"`
	Link string `json:"link" form:"link" query:"link"`
}

type DeleteServiceProRequest struct {
	ServicesID []string `json:"services_id" form:"services_id" query:"services_id"`
}

type DeleteSocialMediaProRequest struct {
	ID string `json:"id" form:"id" query:"id"`
}

type ReplyReviewProRequest struct {
	ReviewID string `json:"review_id" form:"review_id" query:"review_id"`
	Reply    string `json:"reply"     form:"reply"     query:"reply"`
}

type ReviewProRequest struct {
	ProviderID string `json:"provider_id" form:"provider_id" query:"provider_id"`
	Rating     int32  `json:"rating"      form:"rating"      query:"rating"`
	Comment    string `json:"comment"     form:"comment"     query:"comment"`
	Note       string `json:"note"        form:"note"        query:"note"`
	UserName   string `json:"user_name"   form:"user_name"   query:"user_name"`
	HireID     string `json:"hire_id"     form:"hire_id"     query:"hire_id"`
}

type SignUpProRequest struct {
	Name         string `json:"name"         form:"name"         query:"name"`
	Introduction string `json:"introduction" form:"introduction" query:"introduction"`
	Years        int32  `json:"years"        form:"years"        query:"years"`
	Address      string `json:"address"      form:"address"      query:"address"`
}

type UpdateSocialMediaProRequest struct {
	ID   string `json:"id"   form:"id"   query:"id"`
	Name string `json:"name" form:"name" query:"name"`
	Link string `json:"link" form:"link" query:"link"`
}

type UpdateProRequest struct {
	ID           string `json:"id"           form:"id"           query:"id"`
	Name         string `json:"name"         form:"name"         query:"name"`
	Introduction string `json:"introduction" form:"introduction" query:"introduction"`
	Years        int32  `json:"years"        form:"years"        query:"years"`
	Address      string `json:"address"      form:"address"      query:"address"`
}
