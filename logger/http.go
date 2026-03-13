package logger

import (
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog"
	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
)

func (l *LogManager) FiberMiddleware(
	logger zerolog.Logger,
	ipHeader string,
	traceIDHeader string,
	errHandler func(c *fiber.Ctx, err error, showUnknownErrorsInResponse bool),
	filter func(*fiber.Ctx) bool,
) fiber.Handler {
	return func(c *fiber.Ctx) error {
		if filter != nil && filter(c) {
			return c.Next()
		}

		start := time.Now()
		err := c.Next()
		errMessage := ""
		if err != nil {
			errMessage = err.Error()
		}
		if err != nil && errHandler != nil {
			errHandler(c, err, l.cfg.ShowUnknownErrorsInResponse)
		}

		ip := c.Get(ipHeader)
		if ipHeader == "" {
			ip = c.IP()
		}

		traceID, _ := c.Locals(traceIDHeader).(string)
		userID, _ := c.Locals("userID").(int)

		reqCopy := copyRequest(c, l.cfg.MaxHTTPBodySize, l.cfg.SecureReqJsonPaths)
		resCopy := copyResponse(c, l.cfg.MaxHTTPBodySize, l.cfg.SecureResJsonPaths)

		fields := &HttpFields{
			TraceID:      traceID,
			RemoteIP:     ip,
			Method:       c.Method(),
			Host:         c.Hostname(),
			Path:         c.OriginalURL(),
			Protocol:     c.Protocol(),
			UserAgent:    c.Get(fiber.HeaderUserAgent),
			RequestBody:  string(reqCopy),
			ResponseBody: string(resCopy),
			Error:        errMessage,
			StatusCode:   c.Response().StatusCode(),
			Latency:      time.Since(start).Milliseconds(),
			UserID:       userID,
		}

		var toLog zerolog.LogObjectMarshaler
		toLog = &LogStructure{
			HTTP: fields,
		}

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

func copyRequest(c *fiber.Ctx, bodySize int, securePath []string) []byte {
	reqCopy := make([]byte, 0)
	if len(c.Body()) >= bodySize {
		reqCopy, _ = sjson.SetBytes(
			[]byte{},
			"hidden",
			fmt.Sprintf(
				"because it's too big (body length: %d, max size: %d)",
				len(c.Body()),
				bodySize,
			),
		)
	} else {
		reqCopy = c.Request().Body()
		for _, r := range securePath {
			if gjson.GetBytes(reqCopy, r).Exists() {
				reqCopy, _ = sjson.SetBytes(reqCopy, r, "secured in middleware")
			}
		}
	}

	return reqCopy
}

func copyResponse(c *fiber.Ctx, bodySize int, securePath []string) []byte {
	resCopy := c.Response().Body()
	if len(resCopy) >= bodySize {
		resCopy, _ = sjson.SetBytes(
			[]byte{},
			"hidden",
			fmt.Sprintf(
				"because it's too big (body length: %d, max size: %d)",
				len(c.Body()),
				bodySize,
			),
		)
	} else {
		for _, r := range securePath {
			if gjson.GetBytes(resCopy, r).Exists() {
				resCopy, _ = sjson.SetBytes(resCopy, r, "secured in middleware")
			}
		}
	}

	return resCopy
}
