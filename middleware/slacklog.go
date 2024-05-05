package middleware

import (
	"log/slog"
	"os"
	"strings"
	"time"

	slogslack "github.com/samber/slog-slack/v2"
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
