package main

import (
	"flag"
	"os"
	"strings"
	"time"

	"github.com/bancek/pjl/pkg/pjl"
	"github.com/sirupsen/logrus"
)

func main() {
	var jsonLevelField string
	var jsonTimeField string
	var jsonTimeFormat string
	var jsonMsgField string
	var excludeJSONFieldsStr string
	var timeFormat string
	var levelStr string
	var maxScanLineSize int

	flag.StringVar(&jsonLevelField, "jsonLevelField", "level", "Log level JSON field")
	flag.StringVar(&jsonTimeField, "jsonTimeField", "time", "Log time JSON field")
	flag.StringVar(&jsonTimeFormat, "jsonTimeFormat", time.RFC3339Nano, "JSON time format")
	flag.StringVar(&jsonMsgField, "jsonMsgField", "msg", "Log msg JSON field")
	flag.StringVar(&excludeJSONFieldsStr, "excludeJsonFields", "", "Exclude fields (comma separated)")
	flag.StringVar(&timeFormat, "timeFormat", "2006-01-02T15:04:05.000Z07:00", "Output time format")
	flag.StringVar(&levelStr, "level", "debug", "Output level")
	flag.IntVar(&maxScanLineSize, "maxScanLineSize", 1*1024*1024, "Max scan line size in bytes")

	flag.Parse()

	var excludeJSONFields []string

	if excludeJSONFieldsStr != "" {
		excludeJSONFields = strings.Split(excludeJSONFieldsStr, ",")
	}

	level := logrus.DebugLevel

	if lvl, err := logrus.ParseLevel(levelStr); err == nil {
		level = lvl
	}

	formatter := &logrus.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: timeFormat,
	}

	logger := logrus.StandardLogger()
	logger.SetOutput(os.Stdout)
	logger.Formatter = formatter
	logger.Level = level

	pjl := &pjl.PrettyJSONLogger{
		JSONLevelField:    jsonLevelField,
		JSONTimeField:     jsonTimeField,
		JSONTimeFormat:    jsonTimeFormat,
		JSONMsgField:      jsonMsgField,
		ExcludeJSONFields: excludeJSONFields,
		MaxScanLineSize:   maxScanLineSize,
		Input:             os.Stdin,
		Output:            os.Stdout,
		Logger:            logger,
	}

	pjl.Run()
}
