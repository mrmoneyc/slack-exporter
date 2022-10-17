# slack-exporter

> The tool for exporting Slack conversation histories (with replies) and files written in Go.

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://github.com/mrmoneyc/slack-exporter/blob/master/LICENSE)

## Requirement

* Go 1.19+ (for development)
* [Slack App](https://api.slack.com/apps) User / Bot Token with the required scope listed below (refer to **Creating an app**, **Requesting scopes** and **Installing the app to a workspace** sections in this [tutorial](https://api.slack.com/authentication/basics))
  * `users:read`
  * `channels:history`, `channels:read`
  * `groups:history`, `groups:read`
  * `im:history`, `im:read`
  * `mpim:history`, `mpim:read`
  * `files:read`

## Configuration

Add your `config.yaml` file:

```yaml
---
SlackToken:
RequestDelay: 1200
ChannelTypes:
  # - "public_channel"
  # - "private_channel"
  # - "mpim"
  # - "im"
ExportBasePath: "./export"
SplitMessages: false
ArchiveData: true
DownloadFiles: true
IncludeChannel:
  # - "general"
  # - "random"
  # - "@someone"
  # - "mpdm-user1--user2--user3-1"
```

List of configuration values

| Key            | Type         | Description                                                  |
| -------------- | ------------ | ------------------------------------------------------------ |
| SlackToken     | string       | Your Slack app user / bot token.                             |
| RequestDelay   | integer      | Waiting time for the Slack  API call (millisecond).          |
| ChannelTypes   | string array | Specify channel types to export.<br />`public_channel`: Public Channel<br />`private_channel`: Private Channel<br />`mpim`: Group Message<br />`im`: Direct Message |
| ExportBasePath | string       | Data export target path.                                     |
| SplitMessages  | boolean      | Split message files by day if set to `true`.                 |
| ArchiveData    | boolean      | Make archive file for exported data if set to `true`.        |
| DownloadFiles  | boolean      | Download files in messages.                                  |
| IncludeChannel | string array | Specify channels to export.                                  |

## Usage

To run exporter:

```sh
slack-exporter
```

Windows user can use `run.bat` to run exporter.

If you want to switch to other configuration file, you can add `-c` or `--config` option:

```sh
slack-exporter -c ~/another_config.yaml

slack-exporter --config /home/someone/yet_another_config.yaml
```

To change log verbosity, you can add `-v` or `--verbosity` option:

```sh
slack-exporter -v debug

slack-exporter --verbosity warn
```

For more usage, use `help` or `--help`
