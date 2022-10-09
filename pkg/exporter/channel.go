package exporter

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"

	"github.com/mrmoneyc/slack-exporter/pkg/config"
	log "github.com/sirupsen/logrus"
	"github.com/slack-go/slack"
)

func GetChannels(cfg *config.Config, client *slack.Client, users []slack.User) ([]slack.Channel, error) {
	var channels []slack.Channel
	var channels_raw []slack.Channel
	var nextCursor string

	for {
		log.Debugln("call conversations_list (slack api)")
		prm := slack.GetConversationsParameters{
			Cursor:          nextCursor,
			ExcludeArchived: false,
			Limit:           200,
			Types:           cfg.ChannelTypes,
		}
		conversationList, cursor, err := client.GetConversations(&prm)
		if err != nil {
			log.Errorf("%s", err)
			return nil, err
		}

		channels_raw = append(channels_raw, conversationList...)

		time.Sleep(time.Duration(cfg.RequestDelay) * time.Millisecond)

		nextCursor = cursor
		if nextCursor == "" {
			break
		}
		log.Debugf("\tnext cursor: %s", nextCursor)
	}

	log.Infof("total channels: %d", len(channels_raw))

	for _, x := range channels_raw {
		if x.IsIM {
			if x.IsExtShared {
				x.Name = x.User
			} else {
				for _, y := range users {
					if y.ID == x.User {
						x.Name = fmt.Sprintf("@%s", y.Name)
					}
				}
			}
		}

		channels = append(channels, x)
	}

	log.Infof("total channels: %d", len(channels))

	return channels, nil
}

func SaveChannels(cfg *config.Config, channels []slack.Channel, now string) error {
	currDir, err := os.Getwd()
	if err != nil {
		return err
	}

	exportPath := filepath.Join(currDir, cfg.ExportBasePath, now)
	if err := os.MkdirAll(exportPath, 0700); err != nil {
		return err
	}

	log.Infoln("save channels")
	log.Debugf("channels export path: %s", exportPath)

	channelJsonPath := filepath.Join(exportPath, "channels.json")
	log.Debugf("export file path: %s", channelJsonPath)

	channelJson, err := json.MarshalIndent(channels, "", "  ")
	if err != nil {
		return err
	}

	if err := ioutil.WriteFile(channelJsonPath, channelJson, 0644); err != nil {
		return err
	}

	return nil
}
