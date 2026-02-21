package helper

type BaseResponse struct {
	Status  string      `json:"status"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}

func SuccessResponse(message string, data interface{}) BaseResponse {
	return BaseResponse{
		Status:  "success",
		Message: message,
		Data:    data,
	}
}

func ErrorResponse(message string) BaseResponse {
	return BaseResponse{
		Status:  "error",
		Message: message,
	}
}