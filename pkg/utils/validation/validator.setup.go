package validationUtil

import (
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

// SetupValidator gin validator 설정
func SetupValidator() {
	// Register custom validator
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("base64", ValidationBase64())
	}
}
