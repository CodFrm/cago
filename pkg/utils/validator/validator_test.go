package validator

import (
	"testing"

	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"
)

func Test_defaultValidator_registerValidation(t *testing.T) {
	v1, err := NewValidator()
	v := v1.(*DefaultValidator)
	assert.Nil(t, err)
	err = v.validate.Var("11111111", "mobile")
	assert.Equal(t, "手机号码格式不正确", TransError(err.(validator.ValidationErrors)))
	err = v.validate.Var("13000000000", "mobile")
	assert.Nil(t, err)
	s := struct {
		Mobile string `binding:"required,mobile"`
	}{Mobile: "1234567"}
	err = v.ValidateStruct(&s)
	assert.Equal(t, "手机号码格式不正确", TransError(err.(validator.ValidationErrors)))
}
