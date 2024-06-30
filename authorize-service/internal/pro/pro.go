package pro

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func FindPros(c echo.Context) error {
	return c.JSON(http.StatusOK, "FindPros response")
}

func JoinAsProvider(c echo.Context) error {
	return c.JSON(http.StatusOK, "JoinAsProvider response")
}

func AddPaymentMethodPro(c echo.Context) error {
	return c.JSON(http.StatusOK, "AddPaymentMethodPro response")
}

func AddServicePro(c echo.Context) error {
	return c.JSON(http.StatusOK, "AddServicePro response")
}

func AddSocialMediaPro(c echo.Context) error {
	return c.JSON(http.StatusOK, "AddSocialMediaPro response")
}

func DeleteServicePro(c echo.Context) error {
	return c.JSON(http.StatusOK, "DeleteServicePro response")
}

func DeleteSocialMediaPro(c echo.Context) error {
	return c.JSON(http.StatusOK, "DeleteSocialMediaPro response")
}

func GetProviderProfile(c echo.Context) error {
	return c.JSON(http.StatusOK, "GetProviderProfile response")
}

func ReplyReviewPro(c echo.Context) error {
	return c.JSON(http.StatusOK, "ReplyReviewPro response")
}

func ReviewPro(c echo.Context) error {
	return c.JSON(http.StatusOK, "ReviewPro response")
}

func SignUpPro(c echo.Context) error {
	return c.JSON(http.StatusOK, "SignUpPro response")
}

func UpdateSocialMediaPro(c echo.Context) error {
	return c.JSON(http.StatusOK, "UpdateSocialMediaPro response")
}

func FindProById(c echo.Context) error {
	return c.JSON(http.StatusOK, "FindProById response")
}

func DeleteProById(c echo.Context) error {
	return c.JSON(http.StatusOK, "DeleteProById response")
}

func UpdatePro(c echo.Context) error {
	return c.JSON(http.StatusOK, "UpdatePro response")
}

func GetAllReviewsOfProvider(c echo.Context) error {
	return c.JSON(http.StatusOK, "GetAllReviewsOfProvider response")
}

func GetAllServiceOfProvider(c echo.Context) error {
	return c.JSON(http.StatusOK, "GetAllServiceOfProvider response")
}
