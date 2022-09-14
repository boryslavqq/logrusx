package logrusx

import (
	"errors"
	"github.com/sirupsen/logrus"
	"os"
)

type LogField struct {
	Key   string
	Value interface{}
}

type logger struct {
	logrusLogging *logrus.Logger
	fields        logrus.Fields
}

// Create a new logger with JSON configuration and custom service name,
// returns error if service name is invalid
func New(serviceName string) (Logging, error) {
	fieldMap := logrus.FieldMap{}
	fieldMap[logrus.FieldKeyMsg] = "message"

	Logger := logrus.New()
	Fields := logrus.Fields{}
	Logger.SetOutput(os.Stdout)
	Logger.SetFormatter(&logrus.JSONFormatter{
		PrettyPrint: false,
		FieldMap:    fieldMap,
	})
	if serviceName != "" {
		Fields["service"] = serviceName
	} else {
		return nil, errors.New("invalid service name")
	}

	return &logger{logrusLogging: Logger, fields: Fields}, nil
}

type Logging interface {
	Info(msg string)
	Error(msg string, fields ...LogField)
	Fatal(msg string, fields ...LogField)
	addValue(key string, value interface{})
	deleteValue(key string)
	fillFields(fields []LogField)
	deleteFields(fields []LogField)
}

func (l *logger) Info(msg string) {
	l.logrusLogging.WithFields(l.fields).Info(msg)
}

func (l *logger) Error(msg string, fields ...LogField) {
	defer l.deleteFields(fields)

	l.fillFields(fields)
	l.logrusLogging.WithFields(l.fields).Error(msg)
}

func (l *logger) Fatal(msg string, fields ...LogField) {
	defer l.deleteFields(fields)

	l.fillFields(fields)
	l.logrusLogging.WithFields(l.fields).Fatal(msg)
}

func (l *logger) addValue(key string, value interface{}) {
	l.fields[key] = value
}

func (l *logger) deleteValue(key string) {
	delete(l.fields, key)
}

func (l *logger) fillFields(fields []LogField) {
	for _, field := range fields {
		l.addValue(field.Key, field.Value)
	}
}

func (l *logger) deleteFields(fields []LogField) {
	for _, field := range fields {
		l.deleteValue(field.Key)
	}
}
