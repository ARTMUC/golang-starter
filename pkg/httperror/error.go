package httperror

import "fmt"

type ErrCustomError struct {
	Description string `json:"description,omitempty"`
	Metadata    string `json:"metadata,omitempty"`
	StatusCode  int    `json:"statusCode"`
}

func (e ErrCustomError) Error() string {
	return fmt.Sprintf("description: %s,  metadata: %s", e.Description, e.Metadata)
}

func NewHttpError(description, metadata string, statusCode int) ErrCustomError {
	return ErrCustomError{
		Description: description,
		Metadata:    metadata,
		StatusCode:  statusCode,
	}
}
