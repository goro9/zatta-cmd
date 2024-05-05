package middleware

import "github.com/spf13/cobra"

type handler func(cmd *cobra.Command, args []string) error
