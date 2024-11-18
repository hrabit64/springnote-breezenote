package utils

import "github.com/go-playground/validator/v10"

// IsValidationError 유효성 검사 에러인지 확인합니다.
func IsValidationError(err error) bool {
	errs, ok := err.(validator.ValidationErrors)
	return ok && len(errs) > 0
}
