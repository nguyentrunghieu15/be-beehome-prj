package validator

import (
	"errors"
	"fmt"

	"github.com/go-playground/validator/v10"
)

type IValidator interface {
	Validate(interface{}) error
	RegisterRules(rules map[string]string, typeStruct interface{})
}

type ValidatorStuctMap struct {
	validate *validator.Validate
}

func (v *ValidatorStuctMap) RegisterRules(rules map[string]string, typeStruct interface{}) {
	v.validate.RegisterStructValidationMapRules(rules, typeStruct)
}

func (v *ValidatorStuctMap) Validate(e interface{}) error {
	return v.validate.Struct(e)
}

func (*ValidatorStuctMap) Init() interface{} {
	return &ValidatorStuctMap{}
}

func ValidateMap(rules map[string]interface{}, data map[string]interface{}) error {
	validate := validator.New()
	err := validate.ValidateMap(data, rules)
	if len(err) > 0 {
		return errors.New(fmt.Sprintln(err))
	}
	return nil
}
