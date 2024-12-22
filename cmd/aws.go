/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/k0kubun/pp"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

// awsCmd represents the aws command
var awsCmd = &cobra.Command{
	Use:   "aws",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		pp.Println("aws called")
	},
}

func init() {
	rootCmd.AddCommand(awsCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// awsCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// awsCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

var s3WaitObjectExists = &cobra.Command{
	Use: "s3-wait-object-exists",
	RunE: func(cmd *cobra.Command, _ []string) error {
		ctx := cmd.Context()
		cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion("ap-northeast-1"))
		if err != nil {
			return errors.WithStack(err)
		}
		client := s3.NewFromConfig(cfg)

		bucket := "goro9-sandbox"
		key := "test2.txt"

		// b := bytes.NewBufferString("test dayo")
		// if _, err := client.PutObject(ctx, &s3.PutObjectInput{
		// 	Bucket: aws.String(bucket),
		// 	Key:    aws.String(key),
		// 	Body:   b,
		// }); err != nil {
		// 	return errors.WithStack(err)
		// }

		if err := s3.NewObjectExistsWaiter(client, func(options *s3.ObjectExistsWaiterOptions) {
			options.Retryable = func(_ context.Context, _ *s3.HeadObjectInput, _ *s3.HeadObjectOutput, err error) (bool, error) {
				if err == nil {
					return false, nil
				}

				var errorType *types.NotFound
				if errors.As(err, &errorType) {
					return true, nil
				}

				return false, err
			}
		}).Wait(ctx, &s3.HeadObjectInput{
			Bucket: aws.String(bucket),
			Key:    aws.String(key),
		}, 30*time.Second); err != nil {
			return errors.WithStack(err)
		}

		return nil
	},
}

func init() {
	awsCmd.AddCommand(s3WaitObjectExists)
}
