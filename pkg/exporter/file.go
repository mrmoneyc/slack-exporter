package exporter

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/mrmoneyc/slack-exporter/pkg/config"
	log "github.com/sirupsen/logrus"
	"github.com/slack-go/slack"
)

func SaveFiles(cfg *config.Config, client *slack.Client, messages []slack.Message, channelName string, now string) error {
	currDir, err := os.Getwd()
	if err != nil {
		return err
	}

	exportPath := filepath.Join(currDir, cfg.ExportBasePath, now, channelName, "files")
	if err := os.MkdirAll(exportPath, 0700); err != nil {
		return err
	}

	log.Infof("save files of %s", channelName)
	log.Debugf("files export path: %s", exportPath)

	for _, msg := range messages {
		for _, f := range msg.Files {
			log.Debugf("\t* download %s (%s)", f.Name, f.Mode)

			fileName := fmt.Sprintf("%s_%s", f.ID, f.Name)
			fWriter, err := os.Create(filepath.Join(exportPath, fileName))
			if err != nil {
				log.Errorf("%v", err)
				continue
			}

			if err := client.GetFile(f.URLPrivate, fWriter); err != nil {
				log.Errorf("%v", err)
				continue
			}

			time.Sleep(time.Duration(cfg.RequestDelay) * time.Millisecond)
		}
	}

	return nil
}
