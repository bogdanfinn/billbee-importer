package main

import (
	"github.com/applike/gosoline/pkg/application"
	"github.com/applike/gosoline/pkg/cfg"
	"github.com/applike/gosoline/pkg/kernel"
	"github.com/applike/gosoline/pkg/mon"
	"os"
	"path/filepath"
	stockx_billbee_bridge "stockx-billbee-bridge/internal/stockx-billbee-bridge"
	"strings"
)

func main() {
	ex, _ := os.Executable()
	stockxCsvFilePath := filepath.Join(filepath.Dir(ex), "stockx.csv")

	app := MinifiedApp()
	app.Add("stockx-billbee-bridge", stockx_billbee_bridge.New(stockxCsvFilePath))
	app.Run()
}

func MinifiedApp(options ...application.Option) kernel.Kernel {
	ex, _ := os.Executable()
	configFilePath := filepath.Join(filepath.Dir(ex), "config.dist.yml")

	defaults := []application.Option{
		application.WithUTCClock(true),
		application.WithConfigErrorHandlers(defaultErrorHandler),
		application.WithConfigFile(configFilePath, "yml"),
		application.WithConfigFileFlag,
		application.WithConfigEnvKeyReplacer(strings.NewReplacer(".", "_", "-", "_")),
		application.WithConfigSanitizers(cfg.TimeSanitizer),
		application.WithLoggerFormat(mon.FormatGelfFields),
		application.WithLoggerApplicationTag,
		application.WithLoggerTagsFromConfig,
		application.WithLoggerSettingsFromConfig,
		application.WithLoggerContextFieldsMessageEncoder(),
		application.WithLoggerContextFieldsResolver(mon.ContextLoggerFieldsResolver),
		application.WithKernelSettingsFromConfig,
	}

	options = append(defaults, options...)

	return application.New(options...)
}

var defaultErrorHandler = func(err error, msg string, args ...interface{}) {
	logger := mon.NewLogger()
	options := []mon.LoggerOption{
		mon.WithFormat(mon.FormatJson),
		mon.WithTimestampFormat("2006-01-02T15:04:05.999Z07:00"),
		mon.WithHook(mon.NewMetricHook()),
	}

	if err := logger.Option(options...); err != nil {
		panic(err)
	}

	logger.Fatalf(err, msg, args...)
}
