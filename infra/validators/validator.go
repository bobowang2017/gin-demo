package validators

import (
	"github.com/go-playground/validator/v10"
	"reflect"
	"regexp"
)

// 验证字符串是否包含空格
func NotAllowBlank(fl validator.FieldLevel) bool {
	verificationRole := `\s`
	field := fl.Field()
	switch field.Kind() {
	case reflect.String:
		re, err := regexp.Compile(verificationRole)
		if err != nil {
			return false
		}
		return !re.MatchString(field.String())
	default:
		return true
	}
}
