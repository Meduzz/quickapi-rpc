package errorz

type (
	// TODO make something prettier with validation errors
	ErrorDTO struct {
		Code    string `json:"code"`
		Message string `json:"message"`
	}
)

const (
	CodeBadJson  = "JSON"
	CodeBadInput = "VALIDATION"
	CodeGeneric  = "GENERIC" // 😇
)

func CreateError(code string, err error) *ErrorDTO {
	return &ErrorDTO{code, err.Error()}
}

func (e *ErrorDTO) Error() string {
	return e.Message
}
