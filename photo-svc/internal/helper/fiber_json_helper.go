package helper

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	errorcode "github.com/hervibest/be-yourmoments-backup/photo-svc/internal/enum/error"
	"github.com/hervibest/be-yourmoments-backup/photo-svc/internal/helper/logger"
	"github.com/hervibest/be-yourmoments-backup/photo-svc/internal/model"
	"github.com/oklog/ulid/v2"

	"github.com/gofiber/fiber/v2"
)

func StrictBodyParser(ctx *fiber.Ctx, request interface{}) error {
	decoder := json.NewDecoder(bytes.NewReader(ctx.Body()))
	decoder.DisallowUnknownFields()
	return decoder.Decode(request)
}

func ErrBodyParserResponseJSON(ctx *fiber.Ctx, err error) error {
	return ctx.Status(http.StatusBadRequest).JSON(model.BodyParseErrorResponse{
		Success: false,
		Message: "Invalid fields",
		Errors:  err.Error(),
	})
}

func ErrValidationResponseJSON(ctx *fiber.Ctx, validatonErrs *UseCaseValError) error {
	return ctx.Status(http.StatusUnprocessableEntity).JSON(model.ValidationErrorResponse{
		Success: false,
		Message: "Validation error",
		Errors:  validatonErrs.GetValidationErrors(),
	})
}

func ErrUseCaseResponseJSON(ctx *fiber.Ctx, msg string, err error, logs *logger.Log) error {
	if appErr, ok := err.(*AppError); ok {
		if appErr.Err != nil {
			logs.Error(fmt.Sprintf("Internal error in controller : %s [%s]: %v", msg, appErr.Code, appErr.Err.Error()))
		} else {
			logs.Log(fmt.Sprintf("Client error in controller : %s [%s]: %v", msg, appErr.Code, appErr.Message))
		}

		return ctx.Status(appErr.HTTPStatus()).JSON(model.ErrorResponse{
			Success: false,
			Message: appErr.Message,
		})
	}

	return fiber.NewError(fiber.StatusInternalServerError, "Something went wrong. Please try again later")
}

func MultipleULIDSliceParser(ulidSlice []string) error {
	invalidIds := make([]string, 0)
	for _, id := range ulidSlice {
		if _, err := ulid.Parse(id); err != nil {
			invalidIds = append(invalidIds, id)
		}
	}
	if len(invalidIds) != 0 {
		return NewUseCaseError(errorcode.ErrInvalidArgument, fmt.Sprintf("Invalid photo ids : %s", invalidIds))
	}
	return nil
}
