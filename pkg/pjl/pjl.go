package pjl

import (
	"bufio"
	"bytes"
	"encoding/json"
	"io"
	"time"

	"github.com/sirupsen/logrus"
)

var newLine = []byte("\n")

type PrettyJSONLogger struct {
	JSONLevelField    string
	JSONTimeField     string
	JSONTimeFormat    string
	JSONMsgField      string
	ExcludeJSONFields []string
	MaxScanLineSize   int
	Input             io.Reader
	Output            io.Writer
	Logger            *logrus.Logger
}

func (l *PrettyJSONLogger) Run() error {
	scanner := bufio.NewScanner(l.Input)
	scanner.Buffer(make([]byte, 4096), l.MaxScanLineSize)

	for scanner.Scan() {
		l.handleLine(scanner.Bytes())
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	return nil
}

func (l *PrettyJSONLogger) handleLine(b []byte) {
	fields, ok := parseJSONFields(b)
	if !ok {
		l.Output.Write(b)
		l.Output.Write(newLine)
		return
	}

	for _, field := range l.ExcludeJSONFields {
		delete(fields, field)
	}

	level := parseLevel(fields, l.JSONLevelField, logrus.InfoLevel)
	t := parseTime(fields, l.JSONTimeField, l.JSONTimeFormat)
	msg := parseMessage(fields, l.JSONMsgField)

	entry := l.Logger.WithFields(fields)
	entry.Time = t
	entry.Log(level, msg)
}

func parseJSONFields(b []byte) (map[string]interface{}, bool) {
	if !bytes.HasPrefix(b, []byte("{")) {
		return nil, false
	}

	fields := map[string]interface{}{}

	if err := json.Unmarshal(b, &fields); err != nil {
		return nil, false
	}

	return fields, true
}

func parseLevel(fields map[string]interface{}, levelField string, defaultLevel logrus.Level) logrus.Level {
	levelIf, ok := fields[levelField]
	if !ok {
		return defaultLevel
	}
	levelStr, ok := levelIf.(string)
	if !ok {
		return defaultLevel
	}
	level, err := logrus.ParseLevel(levelStr)
	if err != nil {
		return defaultLevel
	}
	delete(fields, levelField)
	return level
}

func parseTime(fields map[string]interface{}, timeField string, timeFormat string) time.Time {
	timeIf, ok := fields[timeField]
	if !ok {
		return time.Now()
	}
	delete(fields, timeField)
	timeStr, ok := timeIf.(string)
	if !ok {
		return time.Now()
	}
	t, err := time.Parse(timeFormat, timeStr)
	if err != nil {
		return time.Now()
	}
	return t
}

func parseMessage(fields map[string]interface{}, msgField string) string {
	msgIf, ok := fields[msgField]
	if !ok {
		return ""
	}
	delete(fields, msgField)
	msgStr, ok := msgIf.(string)
	if !ok {
		return ""
	}
	return msgStr
}
