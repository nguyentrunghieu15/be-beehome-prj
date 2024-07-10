package coverter

import (
	"errors"
	"reflect"
	"strings"

	"github.com/cerbos/cerbos-sdk-go/cerbos"
)

func ToPrincipal(obj interface{}) (*cerbos.Principal, error) {
	// Use reflect to access the object's fields
	val := reflect.ValueOf(obj)
	// Check if the provided obj is a struct
	if val.Kind() != reflect.Struct {
		return nil, errors.New("ToPrincipal: obj must be a struct")
	}
	// Get the 'id' and 'role' fields
	idField := val.FieldByName("UserId")
	roleField := val.FieldByName("Role")

	// Check if the fields were found and are of type string
	if !idField.IsValid() || idField.Kind() != reflect.String {
		return nil, errors.New("ToPrincipal: id field is missing or not a string")
	}
	if !roleField.IsValid() || roleField.Kind() != reflect.String {
		return nil, errors.New("ToPrincipal: role field is missing or not a string")
	}

	// Convert the fields to strings
	id := idField.String()
	role := roleField.String()
	result := cerbos.NewPrincipal(id, role)

	typ := reflect.TypeOf(obj)
	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		fieldType := typ.Field(i)
		// Get the bson tag
		bsonTag := fieldType.Tag.Get("bson")
		if fieldType.Name != "Role" {
			result = result.WithAttr(bsonTag, field.Interface())
		}
	}
	return result, nil
}

func ToResource(obj interface{}) (*cerbos.Resource, error) {
	// Use reflect to access the object's fields
	val := reflect.ValueOf(obj)
	// Check if the provided obj is a struct
	if val.Kind() != reflect.Struct {
		return nil, errors.New("ToPrincipal: obj must be a struct")
	}

	// Convert the fields to strings
	id := "abcxyz"
	typ := reflect.TypeOf(obj)
	structName := strings.ToLower(typ.Name())
	result := cerbos.NewResource(structName, id)

	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		fieldType := typ.Field(i)
		// Get the bson tag
		bsonTag := fieldType.Tag.Get("bson")
		result = result.WithAttr(bsonTag, field.Interface())
	}
	return result, nil
}
