package cmd

import (
	"fmt"
	"os"

	"github.com/dobsondev/xmas-xchange/exchange"
	"github.com/dobsondev/xmas-xchange/output"
	"github.com/dobsondev/xmas-xchange/participants"
	"github.com/spf13/cobra"
)

var (
	exchangeTomlFlag      string
	exchangeDryRunFlag    bool
	exchangeMaxAttempts   int
	exchangeFilenameFlag  string
	exchangeOutputFlag    string
	exchangeS3BucketFlag  string
	exchangeAWSRegionFlag string
)

var exchangeCmd = &cobra.Command{
	Use:   "exchange",
	Short: "Create a new gift exchange assignment and save it to a TOML file",
	RunE: func(cmd *cobra.Command, args []string) error {
		participantsList := participants.GetParticipantsFromToml(&exchangeTomlFlag)

		if exchangeDryRunFlag {
			fmt.Println("===\nDry run mode enabled. No sms messages will be sent.\n===")
		}

		exchanges := exchange.NewExchange(participantsList, exchangeMaxAttempts)

		data, err := exchange.EncodeExchangeToml(exchanges)
		if err != nil {
			return err
		}

		var location string
		switch exchangeOutputFlag {
		case "local":
			location, err = output.WriteLocal(exchangeFilenameFlag, data)
		case "s3":
			bucket := exchangeS3BucketFlag
			if bucket == "" {
				bucket = os.Getenv("S3_BUCKET")
			}
			if bucket == "" {
				return fmt.Errorf("s3 bucket required: pass --s3-bucket or set S3_BUCKET")
			}

			region := exchangeAWSRegionFlag
			if region == "" {
				region = os.Getenv("AWS_REGION")
			}
			if region == "" {
				return fmt.Errorf("aws region required: pass --aws-region or set AWS_REGION")
			}

			location, err = output.WriteS3(cmd.Context(), region, bucket, exchangeFilenameFlag, data)
		default:
			return fmt.Errorf("invalid --output %q: must be \"local\" or \"s3\"", exchangeOutputFlag)
		}
		if err != nil {
			return err
		}

		fmt.Printf("Exchange written to %s\n", location)
		return nil
	},
}

func init() {
	exchangeCmd.Flags().StringVar(&exchangeTomlFlag, "toml", "participants.toml", "Path to the TOML file containing participants' data")
	exchangeCmd.Flags().BoolVar(&exchangeDryRunFlag, "dry-run", true, "If set, the program will not perform any actions, just simulate")
	exchangeCmd.Flags().IntVar(&exchangeMaxAttempts, "max-attempts", 50, "Maximum number of attempts to create a valid exchange")
	exchangeCmd.Flags().StringVar(&exchangeFilenameFlag, "filename", "", "Filename (local path or S3 key) to write the exchange TOML output to")
	exchangeCmd.Flags().StringVar(&exchangeOutputFlag, "output", "local", `Where to write the exchange TOML output: "local" or "s3"`)
	exchangeCmd.Flags().StringVar(&exchangeS3BucketFlag, "s3-bucket", "", "S3 bucket to write to when --output=s3 (falls back to S3_BUCKET env var)")
	exchangeCmd.Flags().StringVar(&exchangeAWSRegionFlag, "aws-region", "", "AWS region to use when --output=s3 (falls back to AWS_REGION env var)")

	exchangeCmd.MarkFlagRequired("filename")

	rootCmd.AddCommand(exchangeCmd)
}
