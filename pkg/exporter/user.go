package exporter

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/mrmoneyc/slack-exporter/pkg/config"
	log "github.com/sirupsen/logrus"
	"github.com/slack-go/slack"
)

func GetUsers(client *slack.Client) ([]slack.User, error) {
	users, err := client.GetUsers()
	if err != nil {
		log.Errorf("%s", err)
		return nil, err
	}

	log.Infof("total users: %d", len(users))

	return users, nil
}

func SaveUsers(cfg *config.Config, users []slack.User, now string) error {
	currDir, err := os.Getwd()
	if err != nil {
		return err
	}

	exportPath := filepath.Join(currDir, cfg.ExportBasePath, now)
	if err := os.MkdirAll(exportPath, 0700); err != nil {
		return err
	}

	log.Infoln("save users")
	log.Debugf("users export path: %s", exportPath)

	userJsonPath := filepath.Join(exportPath, "users.json")
	log.Debugf("export file path: %s", userJsonPath)

	userJson, err := json.MarshalIndent(users, "", "  ")
	if err != nil {
		return err
	}

	if err := ioutil.WriteFile(userJsonPath, userJson, 0644); err != nil {
		return err
	}

	return nil
}
