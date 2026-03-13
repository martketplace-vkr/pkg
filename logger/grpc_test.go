package logger

import (
	"bytes"
	"context"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"google.golang.org/grpc"
	"strings"
	"testing"
)

func TestGRPCMiddleware(t *testing.T) {
	var logBuffer bytes.Buffer

	testLogger := zerolog.New(&logBuffer).With().Timestamp().Logger()

	logManager := &LogManager{}

	middleware := logManager.GRPCMiddleware(&testLogger)

	mockHandler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return "test_response", nil
	}

	mockErrorHandler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return nil, errors.New("test error")
	}

	mockInfo := &grpc.UnaryServerInfo{
		FullMethod: "/test.Service/TestMethod",
	}

	testReq := "test_request"

	tests := []struct {
		name              string
		handler           grpc.UnaryHandler
		expectError       bool
		expectLogContains []string
	}{
		{
			name:        "successful request",
			handler:     mockHandler,
			expectError: false,
			expectLogContains: []string{
				`"rpc.method":"TestMethod"`,
				`"level":"info"`,
				`"logging/structure":"access"`,
			},
		},
		{
			name:        "request with error",
			handler:     mockErrorHandler,
			expectError: true,
			expectLogContains: []string{
				`"rpc.method":"TestMethod"`,
				`"level":"info"`,
				`"logging/structure":"access"`,
				`"test error"`, // Проверяем наличие ошибки в логе
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logBuffer.Reset()

			ctx := context.Background()

			resp, err := middleware(ctx, testReq, mockInfo, tt.handler)

			if tt.expectError && err == nil {
				t.Error("Expected error but got nil")
			}
			if !tt.expectError && err != nil {
				t.Errorf("Expected no error but got: %v", err)
			}

			if !tt.expectError && resp != "test_response" {
				t.Errorf("Expected response 'test_response', got %v", resp)
			}

			logOutput := logBuffer.String()

			for _, expectedText := range tt.expectLogContains {
				if !strings.Contains(logOutput, expectedText) {
					t.Errorf("Log doesn't contain expected text '%s'. Log: %s",
						expectedText, logOutput)
				}
			}
		})
	}
}
