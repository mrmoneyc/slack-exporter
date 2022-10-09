package exporter

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"time"

	"github.com/mrmoneyc/slack-exporter/pkg/config"
	log "github.com/sirupsen/logrus"
	"github.com/slack-go/slack"
)

func GetMessages(cfg *config.Config, client *slack.Client, channelId string) ([]slack.Message, error) {
	log.Infof("get messages of %s", channelId)

	var messages []slack.Message
	var nextCursor string

	for {
		log.Debugf("call conversations_history (slack api): %s", channelId)
		prm := slack.GetConversationHistoryParameters{
			ChannelID: channelId,
			Cursor:    nextCursor,
			Limit:     200,
		}
		conversationHistory, err := client.GetConversationHistory(&prm)
		if err != nil {
			log.Errorf("%s", err)
			return nil, err
		}

		messages = append(messages, conversationHistory.Messages...)

		time.Sleep(time.Duration(cfg.RequestDelay) * time.Millisecond)

		nextCursor = conversationHistory.ResponseMetaData.NextCursor
		if nextCursor == "" {
			break
		}
		log.Debugf("\tnext cursor: %s", nextCursor)
	}

	log.Infof("total messages: %d", len(messages))

	for _, msg := range messages {
		if msg.ThreadTimestamp == msg.Timestamp {
			var nextThreadCursor string

			for {
				log.Debugf("call conversations_replies (slack api): %s", msg.Timestamp)
				prm := slack.GetConversationRepliesParameters{
					ChannelID: channelId,
					Timestamp: msg.ThreadTimestamp,
					Cursor:    nextThreadCursor,
					Limit:     200,
				}
				conversationReplies, _, cursor, err := client.GetConversationReplies(&prm)
				if err != nil {
					log.Errorf("%s", err)
					return nil, err
				}

				for _, reply := range conversationReplies {
					if reply.Timestamp != reply.ThreadTimestamp {
						messages = append(messages, reply)
					}
				}

				time.Sleep(time.Duration(cfg.RequestDelay) * time.Millisecond)

				nextThreadCursor = cursor
				if nextThreadCursor == "" {
					break
				}
				log.Debugf("\tnext cursor: %s", nextThreadCursor)
			}

		}
	}

	log.Infof("total messages (with replies): %d", len(messages))

	sort.SliceStable(messages, func(i, j int) bool {
		return messages[i].Timestamp < messages[j].Timestamp
	})

	return messages, nil
}

func SaveMessages(cfg *config.Config, messages []slack.Message, channelName string, now string) error {
	currDir, err := os.Getwd()
	if err != nil {
		return err
	}

	exportPath := filepath.Join(currDir, cfg.ExportBasePath, now, channelName)
	if err := os.MkdirAll(exportPath, 0700); err != nil {
		return err
	}

	log.Infof("save messages of %s", channelName)
	log.Debugf("messages export path: %s", exportPath)

	if cfg.SplitMessages {
		var isExists = make(map[string]bool)
		var days = []string{}
		for _, m := range messages {
			day, err := formatTimestamp(m.Timestamp)
			if err != nil {
				return err
			}

			if isExists[day] {
				continue
			}

			days = append(days, day)
			isExists[day] = true
		}

		for _, d := range days {
			var dayMessages []slack.Message

			for _, m := range messages {
				day, err := formatTimestamp(m.Timestamp)
				if err != nil {
					return err
				}

				if day == d {
					dayMessages = append(dayMessages, m)
				}
			}

			msgJsonPath := filepath.Join(exportPath, fmt.Sprintf("%s.json", d))
			log.Debugf("export file path: %s", msgJsonPath)

			msgJson, err := json.MarshalIndent(dayMessages, "", "  ")
			if err != nil {
				return err
			}

			if err := ioutil.WriteFile(msgJsonPath, msgJson, 0644); err != nil {
				return err
			}
		}
	} else {
		msgJsonPath := filepath.Join(exportPath, "messages.json")
		log.Debugf("export file path: %s", msgJsonPath)

		msgJson, err := json.MarshalIndent(messages, "", "  ")
		if err != nil {
			return err
		}

		if err := ioutil.WriteFile(msgJsonPath, msgJson, 0644); err != nil {
			return err
		}
	}

	return nil
}

func formatTimestamp(timeStamp string) (string, error) {
	ts, err := strconv.ParseFloat(timeStamp, 64)
	if err != nil {
		return "", err
	}

	sec := int64(math.Round(ts))
	epoch := time.Unix(sec, 0)

	return epoch.Format("2006-01-02"), nil
}
