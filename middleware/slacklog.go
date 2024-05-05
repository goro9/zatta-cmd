package middleware

import (
	"log/slog"
	"os"
	"strings"
	"time"

	slogslack "github.com/samber/slog-slack/v2"
	"github.com/slack-go/slack"
	"github.com/spf13/cobra"
)

var slackChannelID string

func SlackLogFlags(cmd *cobra.Command) {
	cmd.Flags().StringVar(&slackChannelID, "slack", "sandbox", "slack channel id")
}

func SlackLog(next handler) handler {
	return func(cmd *cobra.Command, args []string) error {
		timeout := 3 * time.Second
		logger := slog.New(slogslack.Option{
			Level:    slog.LevelDebug,
			BotToken: os.Getenv("SLACK_BOT_TOKEN"),
			Channel:  slackChannelID,
			Timeout:  timeout,
		}.NewSlackHandler())

		slog.SetDefault(logger)

		argsString := strings.Join(os.Args, " ")
		defer func() {
			// panic handling
			r := recover()
			if r != nil {
				logger.With("panic", r).Error(":sos: " + argsString)
				panic(r)
			}
			time.Sleep(timeout)
		}()

		if err := next(cmd, args); err != nil {
			logger.With("error", err.Error()).Error(":x: " + argsString)
			return err
		}
		logger.Info(":white_check_mark: " + argsString)

		return nil
	}
}

func SlackThreadLog(next handler) handler {
	return func(cmd *cobra.Command, args []string) error {
		botToken := os.Getenv("SLACK_BOT_TOKEN")
		argsString := strings.Join(os.Args, " ")
		_, ts, err := slack.New(botToken).PostMessage(slackChannelID, slack.MsgOptionText("start "+argsString, true))
		if err != nil {
			return err
		}

		timeout := 3 * time.Second
		logger := slog.New(slogslack.Option{
			Level:           slog.LevelDebug,
			BotToken:        os.Getenv("SLACK_BOT_TOKEN"),
			Channel:         slackChannelID,
			Timeout:         timeout,
			ThreadTimestamp: ts,
			BroadcastLevel:  slog.LevelError,
		}.NewSlackHandler())

		slog.SetDefault(logger)

		defer func() {
			// panic handling
			r := recover()
			if r != nil {
				logger.With("panic", r).Error(":sos: exit with panic")
				panic(r)
			}
			time.Sleep(timeout)
		}()

		if err := next(cmd, args); err != nil {
			logger.With("error", err.Error()).Error(":x: exit with error")
			return err
		}
		logger.Info(":white_check_mark: exit")

		return nil
	}
}
