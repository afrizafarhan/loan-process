package responses

type Response struct {
	Status  string      `json:"status"`
	Code    int         `json:"code"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Error   interface{} `json:"error,omitempty"`
}

const (
	M_BAD_REQUEST           = "BAD_REQUEST"
	M_CREATED               = "CREATED"
	M_OK                    = "OK"
	M_INTERNAL_SERVER_ERROR = "INTERNAL_SERVER_ERROR"
	M_UNPROCESSABLE_ENTITY  = "UNPROCESSABLE_ENTITY"
	M_NOT_FOUND             = "NOT_FOUND"
)

func SuccessResponse(status string, code int, message string) *Response {
	return &Response{
		Status:  status,
		Code:    code,
		Message: message,
	}
}

func SuccessResponseWithData(status string, code int, payload interface{}) *Response {
	return &Response{
		Status: status,
		Code:   code,
		Data:   payload,
	}
}

func ErrorResponse(status string, code int, err error) *Response {
	return &Response{
		Status: status,
		Code:   code,
		Error:  err.Error(),
	}
}
