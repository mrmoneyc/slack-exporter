package cmd

import (
	"os"
	"time"

	"github.com/mrmoneyc/slack-exporter/pkg/archiver"
	"github.com/mrmoneyc/slack-exporter/pkg/config"
	"github.com/mrmoneyc/slack-exporter/pkg/exporter"
	log "github.com/sirupsen/logrus"
	"github.com/slack-go/slack"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var rootCmd = &cobra.Command{
	Use:               "slack-exporter",
	Short:             "slack exporter",
	PersistentPreRunE: PersistentPreRunE,
	Run:               Run,
}

var (
	cfgFile   string
	verbosity string
)

func init() {
	rootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "config.yaml", "config file path")
	rootCmd.PersistentFlags().StringVarP(&verbosity, "verbosity", "v", log.InfoLevel.String(), "log level (debug, info, warn, error, fatal, panic")
	cobra.OnInitialize(initConfig)
}

func initConfig() {
	if _, err := os.Stat(cfgFile); err != nil {
		if os.IsNotExist(err) {
			log.WithField("configfile", cfgFile).WithError(err).Fatalf("config file not found")
		}
	}

	viper.SetConfigType("yaml")
	viper.SetConfigFile(cfgFile)

	if err := viper.ReadInConfig(); err != nil {
		log.WithField("configfile", cfgFile).WithError(err).Fatalf("unable to read config")
	}
}

func PersistentPreRunE(cmd *cobra.Command, args []string) error {
	lv, err := log.ParseLevel(verbosity)
	if err != nil {
		log.WithField("verbosity", verbosity).WithError(err).Fatalf("invalid log level")
		return err
	}

	log.SetLevel(lv)
	log.Debugf("set log level to: %s", lv.String())

	return nil
}

func Run(cmd *cobra.Command, args []string) {
	now := time.Now().Format("20060102_150405")

	cfg := &config.Config{}
	if err := viper.Unmarshal(&cfg); err != nil {
		log.WithField("configfile", cfgFile).WithError(err).Warnf("cannot unmarshal config file")
	}

	var slackClient = slack.New(cfg.SlackToken)

	log.Infoln("=== start slack export ===")

	users, err := exporter.GetUsers(slackClient)
	if err != nil {
		log.WithError(err).Fatalf("unable to get user")
	}

	channels, err := exporter.GetChannels(cfg, slackClient, users)
	if err != nil {
		log.WithError(err).Fatalf("unable to get channel")
	}

	err = exporter.SaveUsers(cfg, users, now)
	if err != nil {
		log.WithError(err).Fatalf("unable to store user")
	}

	err = exporter.SaveChannels(cfg, channels, now)
	if err != nil {
		log.WithError(err).Fatalf("unable to store channel")
	}

	for _, c := range channels {
		log.Infof("ID: %s, Name: %s", c.ID, c.Name)

		messages, err := exporter.GetMessages(cfg, slackClient, c.ID)
		if err != nil {
			log.Errorf("%v", err)
			continue
		}

		if err := exporter.SaveMessages(cfg, messages, c.Name, now); err != nil {
			log.Errorf("%v", err)
		}

		if err := exporter.SaveFiles(cfg, slackClient, messages, c.Name, now); err != nil {
			log.Errorf("%v", err)
		}
	}

	if cfg.ArchiveData {
		archiver.MakeArchive(cfg, now)
	}

	log.Infoln("=== finish slack export ===")
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.WithError(err).Fatalf("cannot execute command")
		os.Exit(1)
	}
}
