package apiresponse

type ApiResponseMessage string

const (
	Ok             ApiResponseMessage = "ok"
	InvalidRequest ApiResponseMessage = "invalid request"
	UnknownEvent   ApiResponseMessage = "unknown event"
	NotFound       ApiResponseMessage = "not found"
)

type ApiResponse struct {
	Message ApiResponseMessage `json:"message"`
}
