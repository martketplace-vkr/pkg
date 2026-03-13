package http

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/getsentry/sentry-go"
	"github.com/gofiber/fiber/v2"

	"go.opentelemetry.io/otel/attribute"

	"github.com/rs/zerolog"

	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"

	"github.com/martketplace-vkr/pkg/errs"
	"github.com/martketplace-vkr/pkg/logger/log"
	"github.com/martketplace-vkr/pkg/tracer"
)

const (
	UserLocalsKey    = "user"
	ErrNotFoundInCtx = "user not found in context"
)

type (
	response struct {
		TraceID string `json:"traceID"`
		Message string `json:"message"`
		Index   int    `json:"index"`
		Code    int    `json:"code"`
		Params  any    `json:"params"`
	}

	User struct {
		ID            int64  `json:"id"`
		Login         string `json:"login"`
		FirstName     string `json:"first_name"`
		SecondName    string `json:"last_name"`
		Role          string `json:"role"`
		PermissionKey int64
	}
)

func ErrorsMiddleware(
	logger zerolog.Logger,
	ipHeader string,
	getUserFromCtxFunc func(c *fiber.Ctx) any,
	secureReqJsonPaths []string,
	secureResJsonPaths []string,
	showInternalErrorsInResponse bool,
) fiber.Handler {
	return func(c *fiber.Ctx) error {
		_, span := tracer.StartFiberTrace(c)
		defer span.End()

		start := time.Now()
		err := c.Next()

		errMessage := ""
		statusCode := c.Response().StatusCode()
		traceID := tracer.GetTraceID(span)
		resp := response{TraceID: traceID, Message: ""}

		if err != nil {
			if fiberErr := (*fiber.Error)(nil); errors.As(err, &fiberErr) {
				errMessage = fiberErr.Message
				statusCode = fiberErr.Code
			} else if customErr, ok := errs.Parse(err); ok {
				if customErr.Code == errs.ErrCodeUnknown {
					customErr.Code = errs.ErrCodeInternal
				}

				errMessage = customErr.Message
				statusCode = int(customErr.Code)
				resp.Index = customErr.Index
				resp.Params = customErr.Params
			} else {
				errMessage = err.Error()
				statusCode = fiber.StatusInternalServerError
				sentry.CaptureException(
					errs.New(
						errs.ErrCodeUnknown,
						0,
						fmt.Sprintf("err: %s traceID: %s", errMessage, traceID),
					),
				)
			}

			if statusCode >= fiber.StatusInternalServerError {
				if showInternalErrorsInResponse {
					// Это для работы в деве
					resp.Message = fmt.Sprintf("на проде этого сообщения не будет: %s", err.Error())
					resp.Code = statusCode
				} else {
					// Скрываем инфу об ошибке на проде, оставляем только индекс
					resp.Message = fiber.ErrInternalServerError.Message
					resp.Code = fiber.StatusInternalServerError
					resp.Params = nil
					log.ErrorSentry(err, traceID)
				}
			} else {
				// Клиентские ошибки полностью выводятся клиенту
				resp.Message = errMessage
				resp.Code = statusCode
			}

			respJSON, errMarshalJSON := json.Marshal(resp)
			if errMarshalJSON != nil {
				log.Errorf("http_errors Init failed MarshalJson, error: %v", errMarshalJSON)
				return nil
			}

			c.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSON)
			c.Response().SetStatusCode(statusCode)
			c.Response().SetBody(respJSON)
		}

		if strings.Contains(c.OriginalURL(), "/health_check") {
			return nil
		}

		ip := c.Get(ipHeader)
		if ipHeader == "" {
			ip = c.IP()
		}

		reqCopy := c.Request().Body()
		for _, r := range secureReqJsonPaths {
			if gjson.GetBytes(reqCopy, r).Exists() {
				reqCopy, _ = sjson.SetBytes(reqCopy, r, "secured in middleware")
			}
		}

		resCopy := c.Response().Body()
		for _, r := range secureResJsonPaths {
			if gjson.GetBytes(resCopy, r).Exists() {
				resCopy, _ = sjson.SetBytes(resCopy, r, "secured in middleware")
			}
		}

		if len(reqCopy) >= 1000 {
			reqCopy = []byte("too long")
		}

		if len(resCopy) >= 1000 {
			resCopy = []byte("too long")
		}

		// if statusCode == fiber.StatusNotFound ||
		// 	statusCode == fiber.StatusMethodNotAllowed {
		// 	if alerter != nil {
		// 		alerter.Alert(c, Scrapping)
		// 	}
		// }

		fields := &LoggingData{
			TraceID:      traceID,
			RemoteIP:     ip,
			Method:       c.Method(),
			Host:         c.Hostname(),
			Path:         c.OriginalURL(),
			Protocol:     c.Protocol(),
			RequestBody:  string(reqCopy),
			ResponseBody: string(resCopy),
			Error:        errMessage,
			StatusCode:   c.Response().StatusCode(),
			Latency:      time.Since(start).Milliseconds(),
			User:         getUserFromCtxFunc(c),
		}
		toLog := zerolog.LogObjectMarshaler(fields)

		switch {
		case fields.StatusCode >= 500:
			logger.Error().EmbedObject(toLog).Msg("server error")
		case fields.StatusCode >= 400:
			logger.Error().EmbedObject(toLog).Msg("client error")
		case fields.StatusCode >= 300:
			logger.Warn().EmbedObject(toLog).Msg("redirect")
		case fields.StatusCode >= 200:
			logger.Info().EmbedObject(toLog).Msg("success")
		case fields.StatusCode >= 100:
			logger.Info().EmbedObject(toLog).Msg("informative")
		default:
			logger.Warn().EmbedObject(toLog).Msg("unknown status")
		}

		return nil
	}
}

func getUserFromCtxFunc(c *fiber.Ctx) any {
	user, ok := c.Locals(UserLocalsKey).(*User)
	if !ok {
		return UserNotFoundErr{
			Msg: ErrNotFoundInCtx,
		}
	}
	return user
}

func TraceMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		_, span := tracer.StartFiberTrace(c)
		defer span.End()

		span.SetAttributes(
			attribute.KeyValue{
				Key:   "http.method",
				Value: attribute.StringValue(c.Method()),
			},
			attribute.KeyValue{
				Key:   "http.path",
				Value: attribute.StringValue(c.Path()),
			},
		)

		return c.Next()
	}
}
