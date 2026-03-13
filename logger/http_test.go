package logger

import (
	"encoding/json"
	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestFiberMiddleware(t *testing.T) {
	var logged strings.Builder
	logger := zerolog.New(&logged).With().Timestamp().Logger()

	app := fiber.New()

	lm := NewLogger(Config{
		SkipFrameCount:     3,
		SecureReqJsonPaths: []string{"password", "token"},
		SecureResJsonPaths: []string{"password", "token"},
		MaxHTTPBodySize:    500,
	})

	app.Use(lm.FiberMiddleware(
		logger,
		"",
		"trace-id",
		func(c *fiber.Ctx, err error, show bool) {
		},
		nil,
	))

	app.Post("/test", func(c *fiber.Ctx) error {
		c.Locals("trace-id", "12345")
		c.Locals("userID", 42)

		return c.JSON(fiber.Map{
			"token": "super-secret-token",
		})
	})

	reqBody := `{"login":"admin","password":"qwerty"}`
	req := httptest.NewRequest(http.MethodPost, "/test", strings.NewReader(reqBody))

	req.Header.Set("User-Agent", "GoTest")
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req, -1)

	require.NoError(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)

	var parsedLog struct {
		HTTP struct {
			RequestBody  string `json:"requestBody"`
			ResponseBody string `json:"responseBody"`
		} `json:"HTTP"`
	}

	err = json.Unmarshal([]byte(logged.String()), &parsedLog)
	require.NoError(t, err)

	var requestData map[string]interface{}
	err = json.Unmarshal([]byte(parsedLog.HTTP.RequestBody), &requestData)
	require.NoError(t, err)

	var responseData map[string]interface{}
	err = json.Unmarshal([]byte(parsedLog.HTTP.ResponseBody), &responseData)
	require.NoError(t, err)

	assert.Equal(t, "secured in middleware", requestData["password"])
	assert.Equal(t, "admin", requestData["login"])

	assert.Equal(t, "secured in middleware", responseData["token"])
}
