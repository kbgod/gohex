package handler

import (
	"errors"
	"fmt"
	"net/http"

	domainErrors "app/internal/core/error"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/schema"
)

type ClientError struct {
	StatusCode int            `json:"-"`
	Code       int            `json:"code"`
	Message    string         `json:"message"`
	Data       map[string]any `json:"data,omitempty"`
}

func (c *ClientError) Error() string {
	return fmt.Sprintf("client error - code: %d; message: %s", c.Code, c.Message)
}

type ValidationError struct {
	ClientError

	Field string `json:"field"`
}

func (ve *ValidationError) Error() string {
	return fmt.Sprintf("validation error - field: %s; message: %s", ve.Field, ve.Message)
}

func newValidationError(field, message string, code ...int) *ValidationError {
	ve := &ValidationError{
		ClientError: ClientError{
			Message: message,
		},
		Field: field,
	}
	if len(code) > 0 {
		ve.Code = code[0]
	} else {
		ve.Code = http.StatusUnprocessableEntity
	}

	return ve
}

func newBindError(err error) error {
	unwrappedErr := errors.Unwrap(err)

	var multiErr schema.MultiError
	if errors.As(unwrappedErr, &multiErr) {
		var (
			conversionErr schema.ConversionError
			unknownKeyErr schema.UnknownKeyError
			emptyFieldErr schema.EmptyFieldError
		)

		for _, singleErr := range multiErr {
			switch {
			case errors.As(singleErr, &conversionErr):
				return newValidationError(conversionErr.Key, conversionErr.Err.Error())
			case errors.As(singleErr, &unknownKeyErr):
				return newBadRequest(unknownKeyErr.Error())
			case errors.As(singleErr, &emptyFieldErr):
				return newValidationError(emptyFieldErr.Key, emptyFieldErr.Error())
			}
		}
	}

	return newBadRequest(err.Error())
}

func newBadRequest(message string, code ...int) *ClientError {
	ce := &ClientError{
		Message: message,
	}
	if len(code) > 0 {
		ce.Code = code[0]
	} else {
		ce.Code = http.StatusBadRequest
	}

	if ce.Code >= 400 && ce.Code < 500 {
		ce.StatusCode = ce.Code
	} else {
		ce.StatusCode = http.StatusBadRequest
	}

	return ce
}

func ErrorHandler(ctx fiber.Ctx, err error) error {
	var (
		clientErr     *ClientError
		validationErr *ValidationError
		fiberErr      *fiber.Error
		domainErr     *domainErrors.DomainError
	)

	switch {
	case errors.As(err, &clientErr):
		return ctx.Status(clientErr.StatusCode).JSON(clientErr)
	case errors.As(err, &validationErr):
		return ctx.Status(http.StatusUnprocessableEntity).JSON(validationErr)
	case errors.As(err, &fiberErr):
		return ctx.Status(fiberErr.Code).JSON(&ClientError{
			Code:    fiberErr.Code,
			Message: fiberErr.Message,
		})
	case errors.As(err, &domainErr):
		statusCode := http.StatusBadRequest

		errResponse := &ClientError{
			Message: domainErr.Message(),
		}

		if domainErr.Code() == 0 {
			errResponse.Code = http.StatusBadRequest
		} else {
			errResponse.Code = domainErr.Code()
			if domainErr.Code() >= http.StatusBadRequest && domainErr.Code() < http.StatusInternalServerError {
				statusCode = domainErr.Code()
			}
		}

		if len(domainErr.Args()) > 0 {
			errResponse.Data = make(map[string]any, len(domainErr.Args()))
			for _, arg := range domainErr.Args() {
				errResponse.Data[arg.Key] = arg.Value
			}
		}

		return ctx.Status(statusCode).JSON(errResponse)
	default:
		// h.observer.Logger.
		//	Error().
		//	Ctx(ctx.Context()).
		//	Err(err).
		//	Msg("unhandled error")
		// LOGGED BY ROUTER LOGGER
		ctx.Res().Status(http.StatusInternalServerError)

		return nil
	}
}

type MessageResponse struct {
	Message string `json:"message" example:"ok"`
}
