package validator

import (
	"github.com/go-playground/validator/v10"
)

type IValidator interface {
	Validate(interface{}) error
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

func ContructValidateStructMap() *ValidatorStuctMap {
	return &ValidatorStuctMap{}
}
