package logs_test

import (
	"context"
	"log/slog"
	"testing"

	"github.com/openinnovationai/k8s-amd-exporter/internal/application/logs"
	"github.com/openinnovationai/k8s-amd-exporter/internal/application/settings"
	"github.com/stretchr/testify/assert"
)

func TestMakeLogger(t *testing.T) {
	t.Parallel()
	// Given
	testCases := map[string]struct {
		conf *settings.Configuration
		want slog.Level
	}{
		"production": {
			conf: &settings.Configuration{
				LogLevel: "production",
			},
			want: slog.LevelInfo,
		},
		"development": {
			conf: &settings.Configuration{
				LogLevel: "development",
			},
			want: slog.LevelDebug,
		},
	}

	for testName, testData := range testCases {
		t.Run(testName, func(t *testing.T) {
			// When
			logger := logs.MakeLogger(testData.conf)
			// Then
			assert.True(t, logger.Handler().Enabled(context.TODO(), testData.want))
		})
	}
}
