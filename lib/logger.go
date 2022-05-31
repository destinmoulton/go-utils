package lib

import (
	"io"
	"os"
	"path/filepath"

	log "github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
)

type LoggerData struct{}

var Logger LoggerData

func (l *LoggerData) SetupLogger(logpath string, util string) {

	lumberjackLogger := &lumberjack.Logger{
		// Log file abbsolute path, os agnostic
		Filename:   filepath.ToSlash(filepath.Join(logpath, util+".log")),
		MaxSize:    5, // MB
		MaxBackups: 10,
		MaxAge:     30,   // days
		Compress:   true, // disabled by default
	}

	// Fork writing into two outputs
	multiWriter := io.MultiWriter(os.Stderr, lumberjackLogger)

	log.SetFormatter(&log.JSONFormatter{})
	log.SetLevel(log.InfoLevel)
	log.SetOutput(multiWriter)
}
