package config

type Config struct {
	SlackToken     string
	RequestDelay   int64
	ChannelTypes   []string
	ExportBasePath string
	SplitMessages  bool
	ArchiveData    bool
}
