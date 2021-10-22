package validator

import (
	"reflect"

	"github.com/go-playground/locales/ja"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	ja_translations "github.com/go-playground/validator/v10/translations/ja"
)

var (
	uni      *ut.UniversalTranslator
	validate *validator.Validate
	trans    ut.Translator
)

func init() {
	ja := ja.New()
	uni = ut.New(ja, ja)
	t, _ := uni.GetTranslator("ja")
	trans = t
	validate = validator.New()
	validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		fieldName := fld.Tag.Get("ja")
		if fieldName == "-" {
			return ""
		}
		return fieldName
	})
	ja_translations.RegisterDefaultTranslations(validate, trans)
}

func Get() *validator.Validate {
	return validate
}

func GetStructFieldByFieldName(s interface{}, field_name string) (reflect.StructField, bool) {
	sf, ok := reflect.TypeOf(s).Elem().FieldByName(field_name)
	if !ok {
		return sf, false
	}

	return sf, true
}

func GetStructTagByFieldName(s interface{}, field_name string) (reflect.StructTag, bool) {
	if sf, found := GetStructFieldByFieldName(s, field_name); found {
		return sf.Tag, true

	} else {
		return "", false
	}
}

func GetJaTagByFieldName(s interface{}, field_name string) (string, bool) {
	if sf, found := GetStructFieldByFieldName(s, field_name); found {
		ja := sf.Tag.Get("ja")
		if ja == "" {
			return "", false
		} else {
			return ja, true
		}
	} else {
		return "", false
	}
}

func GetErrorMessages(s interface{}, err error) []string {
	ret := []string{}
	if err == nil {
		return ret
	}

	for _, err_field := range err.(validator.ValidationErrors) {
		ret = append(ret, err_field.Translate(trans))
	}

	return ret
}
