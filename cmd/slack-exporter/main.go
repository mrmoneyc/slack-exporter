package main

import (
	"io"
	"os"

	"github.com/mrmoneyc/slack-exporter/pkg/cmd"
	log "github.com/sirupsen/logrus"
)

func init() {
	log.SetFormatter(&log.TextFormatter{
		DisableColors: false,
		ForceColors:   true,
		FullTimestamp: true,
	})
}

func main() {
	logFile, err := os.OpenFile("slack-exporter.log", os.O_RDWR|os.O_APPEND|os.O_CREATE, 0660)
	if err != nil {
		log.WithError(err).Warnf("unable to create log file")
	}
	defer logFile.Close()

	mw := io.MultiWriter(os.Stdout, logFile)
	log.SetOutput(mw)

	cmd.Execute()
}
