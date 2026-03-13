package logger

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/rs/zerolog"

	"github.com/martketplace-vkr/pkg/logger/log"
)

type LogManager struct {
	cfg Config
}

func NewLogger(cfg Config) *LogManager {
	return &LogManager{cfg: cfg}
}

func (l *LogManager) InitLogger() {
	lvl, ok := loggerLevelMap[l.cfg.Level]
	if !ok {
		lvl = zerolog.InfoLevel
	}

	if l.cfg.PrettyLogging {
		log.Logger = zerolog.New(
			zerolog.MultiLevelWriter(zerolog.ConsoleWriter{
				Out:           os.Stdout,
				TimeFormat:    zerolog.TimeFieldFormat,
				FormatMessage: formatMessage,
				FormatCaller:  formatCaller,
			})).Level(lvl).With().
			CallerWithSkipFrameCount(l.cfg.SkipFrameCount).Timestamp().Logger()
	} else {
		log.Logger = zerolog.New(os.Stdout).Level(lvl).With().
			CallerWithSkipFrameCount(l.cfg.SkipFrameCount).Timestamp().Logger()
	}

	log.Infof("logger started with settings %+v", l.cfg)
}

func formatMessage(i interface{}) string {
	if i != nil {
		return fmt.Sprintf("| %s |", i)
	}
	return ""
}

func formatCaller(i interface{}) string {
	return filepath.Base(fmt.Sprintf("%s", i))
}
