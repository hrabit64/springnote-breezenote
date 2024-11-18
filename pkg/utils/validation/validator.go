package validationUtil

import (
	"github.com/go-playground/validator/v10"
	"github.com/hrabit64/springnote-breezenote/pkg/utils"
	"regexp"
)

func validateRegex(regex, value string) bool {
	reg := regexp.MustCompile(regex)
	return reg.MatchString(value)
}

func ValidationBase64() validator.Func {
	return func(fl validator.FieldLevel) bool {
		if value, ok := fl.Field().Interface().(string); ok {
			return utils.IsBase64(value)
		}
		return true
	}
}

func CheckUploadImageExt(fileName string) bool {

	return validateRegex(utils.UploadExt, fileName)
}

func CheckDownloadImageExt(fileName string) bool {

	return validateRegex(utils.DlExt, fileName)
}
