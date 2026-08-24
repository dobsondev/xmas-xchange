package cmd

import (
	"fmt"
	"os"

	"github.com/dobsondev/xmas-xchange/exchange"
	"github.com/dobsondev/xmas-xchange/output"
	"github.com/spf13/cobra"
)

var printAWSRegionFlag string

var printCmd = &cobra.Command{
	Use:   "print [file]",
	Short: "Print a previously generated exchange TOML file (local path or s3:// URI)",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		location := args[0]

		var data []byte
		var err error

		if bucket, key, ok := output.ParseS3URI(location); ok {
			region := printAWSRegionFlag
			if region == "" {
				region = os.Getenv("AWS_REGION")
			}
			if region == "" {
				return fmt.Errorf("aws region required: pass --aws-region or set AWS_REGION")
			}

			data, err = output.ReadS3(cmd.Context(), region, bucket, key)
		} else {
			data, err = os.ReadFile(location)
		}
		if err != nil {
			return err
		}

		exchanges, err := exchange.DecodeExchangeToml(data)
		if err != nil {
			return err
		}

		exchange.PrintExchange(exchanges)
		return nil
	},
}

func init() {
	printCmd.Flags().StringVar(&printAWSRegionFlag, "aws-region", "", "AWS region to use when reading an s3:// location (falls back to AWS_REGION env var)")

	rootCmd.AddCommand(printCmd)
}
