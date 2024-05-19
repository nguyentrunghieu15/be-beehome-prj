package convert

import (
	"fmt"
	"testing"

	authapi "github.com/nguyentrunghieu15/be-beehome-prj/api/auth-api"
)

func TestStructProtoToMap(t *testing.T) {
	request := authapi.SignUpRequest{
		Email:     "user@example.com",
		Password:  "secret",
		Phone:     "+1234567890",
		FirstName: "John",
		LastName:  "Doe",
	}

	dataMap, _ := StructProtoToMap(request)

	fmt.Println(dataMap)
}
