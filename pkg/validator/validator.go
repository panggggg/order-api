package validator

import (
	"reflect"
	"strings"

	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	en_translations "github.com/go-playground/validator/v10/translations/en"
)

var (
	Validate *validator.Validate
	Trans    ut.Translator
)

func NewValidator() error {
	var (
		en  = en.New()
		uni = ut.New(en, en)
	)

	Validate = validator.New()
	Trans, _ = uni.GetTranslator("en")

	en_translations.RegisterDefaultTranslations(Validate, Trans)

	Validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})

	// Custom Translation (Example)
	errors := []error{
		Validate.RegisterTranslation(
			"required",
			Trans,
			func(ut ut.Translator) error {
				return ut.Add("required", "{0} is a required field", true)
			},
			func(ut ut.Translator, fe validator.FieldError) string {
				t, _ := ut.T("required", fe.Field())
				return t
			},
		),
	}

	for _, err := range errors {
		if err != nil {
			// Should log here
			return err
		}
	}

	return nil
}
