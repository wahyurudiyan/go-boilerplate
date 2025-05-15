package common

import "github.com/gin-gonic/gin"

type RESTBody[T any] struct {
	Data    T              `json:"data,omitempty"`
	Message string         `json:"message,omitempty"`
	Error   *RESTBodyError `json:"error,omitempty"`
}

type RESTBodyError struct {
	Code   int
	Reason string
}

func BindJSON[T any](c *gin.Context) (RESTBody[T], error) {
	var body RESTBody[T]
	if err := c.BindJSON(&body); err != nil {
		return RESTBody[T]{}, err
	}

	return body, nil
}

func RESTSuccessResponse[T any](message string, data T) RESTBody[T] {
	return RESTBody[T]{
		Data:    data,
		Message: message,
	}
}

func RESTErrorResponse[T any](code int, reason string) RESTBody[T] {
	return RESTBody[T]{
		Error: &RESTBodyError{
			Code:   code,
			Reason: reason,
		},
	}
}
