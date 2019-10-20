package main

import (
	"errors"
	"github.com/sirupsen/logrus"
)

type LoggerConfig struct {
	ReportCaller bool `yaml:"caller"`
	Level string `yaml:"level"`
	Formatter string `yaml:"formatter"`
}

func NewLogger(conf *LoggerConfig) (*logrus.Logger, error) {
	lg := logrus.New()
	lg.SetReportCaller(conf.ReportCaller)
	level, err := logrus.ParseLevel(conf.Level)
	if err != nil {
		return nil, err
	}
	lg.SetLevel(level)
	switch conf.Formatter {
	case "json":
		lg.SetFormatter(&logrus.JSONFormatter{})
	case "text":
		lg.SetFormatter(&logrus.TextFormatter{})
	default:
		return nil, errors.New("unknown logger formatter")
	}
	return lg, nil
}