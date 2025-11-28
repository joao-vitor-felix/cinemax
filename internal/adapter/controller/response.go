package controller

type Response struct {
	Data   any
	Status int
}

type ErrorResponse struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}
