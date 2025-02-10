package conf

import (
	"errors"
	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	enTranslations "github.com/go-playground/validator/v10/translations/en"
	"path"
)

var (
	validate *validator.Validate
	trans    ut.Translator
)

func init() {
	fallback := en.New()
	validate = validator.New()
	uni := ut.New(fallback)
	trans, _ = uni.GetTranslator(fallback.Locale())
	_ = enTranslations.RegisterDefaultTranslations(validate, trans)
	_ = validate.RegisterValidation("dir", func(fl validator.FieldLevel) bool {
		return path.IsAbs(fl.Field().String())
	})
}
func ValidateStruct(s any) error {
	err := validate.Struct(s)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			return errors.New(err.Translate(trans))
		}
	}
	return err
}
