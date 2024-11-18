package validationUtil

import "github.com/go-playground/validator/v10"

// CreateValidationErrorDetail 유효성 검사 에러 메시지를 생성합니다.
// message: 요청 데이터가 올바르지 않습니다. {필드} 필드 (은)는 {태그}이어야 합니다.
func CreateValidationErrorDetail(err validator.ValidationErrors) string {
	message := "요청 데이터가 올바르지 않습니다. "
	for _, e := range err {
		message += e.Field() + " 필드 (은)는 " + e.Tag() + "이어야 합니다. "
	}
	return message
}
