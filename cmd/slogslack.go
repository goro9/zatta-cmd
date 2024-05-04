/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"log/slog"
	"time"

	slogslack "github.com/samber/slog-slack/v2"
	"github.com/slack-go/slack"
	"github.com/spf13/cobra"
)

// slogslackCmd represents the slogslack command
var slogslackCmd = &cobra.Command{
	Use:   "slogslack",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("slogslack called")
	},
}

func init() {
	rootCmd.AddCommand(slogslackCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// slogslackCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// slogslackCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

var testWebhookCmd = &cobra.Command{
	Use:  "test_webhook",
	Args: cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		timeout := time.Duration(3000) * time.Millisecond
		logger := slog.New(slogslack.Option{
			Level:      slog.LevelDebug,
			WebhookURL: args[0],
			Channel:    args[1],
			Timeout:    timeout,
		}.NewSlackHandler())

		ctx := cmd.Context()
		logger.DebugContext(ctx, "test", slog.String("test", "test"), slog.String("test2", "test2"))

		time.Sleep(timeout)
		return nil
	},
}

var testBotCmd = &cobra.Command{
	Use:  "test_bot",
	Args: cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		botToken := args[0]
		channelID := args[1]

		ctx := cmd.Context()
		_, ts, err := slack.New(botToken).PostMessageContext(ctx, channelID,
			slack.MsgOptionText("start task", true),
		)
		if err != nil {
			return err
		}

		timeout := time.Duration(3000) * time.Millisecond
		logger := slog.New(slogslack.Option{
			Level:    slog.LevelDebug,
			BotToken: botToken,
			Channel:  channelID,
			Timeout:  timeout,
		}.NewSlackHandler())

		ctx = slogslack.WithThreadTimestamp(ctx, ts)

		logger.DebugContext(ctx, "test", slog.String("test", "test"), slog.String("test2", "test2"))
		logger.ErrorContext(slogslack.WithReplyBroadcast(ctx), "error", slog.String("error", "hogehoge"))
		logger.InfoContext(ctx, "end task")

		time.Sleep(timeout)
		return nil
	},
}

func init() {
	slogslackCmd.AddCommand(testWebhookCmd)
	slogslackCmd.AddCommand(testBotCmd)
}
