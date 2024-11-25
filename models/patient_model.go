package models

// PatientSwagger is used for Swagger documentation
type PatientSwagger struct {
	ID     int    `json:"id"`
	Name   string `json:"name"`
	Age    int    `json:"age"`
	Gender string `json:"gender"`
	Notes  string `json:"notes"`
}

// ErrorResponse represents an error response structure.
type ErrorResponse struct {
	Error string `json:"error"`
}
