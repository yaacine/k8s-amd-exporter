package logs_test

import (
	"context"
	"log/slog"
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.com/openinnovationai/platform/infra/amd/amd_smi_exporter_v2/internal/application/logs"
	"gitlab.com/openinnovationai/platform/infra/amd/amd_smi_exporter_v2/internal/application/settings"
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
