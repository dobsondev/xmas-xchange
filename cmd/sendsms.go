package cmd

import (
	"fmt"
	"os"

	"github.com/dobsondev/xmas-xchange/exchange"
	"github.com/dobsondev/xmas-xchange/output"
	"github.com/dobsondev/xmas-xchange/sms"
	"github.com/spf13/cobra"
	"github.com/twilio/twilio-go"
)

var (
	sendSMSAWSRegionFlag string
	sendSMSNameFlag      string
	sendSMSDryRunFlag    bool
	sendSMSQuietFlag     bool
	sendSMSTestFlag      bool
)

var sendSMSCmd = &cobra.Command{
	Use:   "sendsms [file]",
	Short: "Send SMS notifications for a saved exchange via Twilio (local path or s3:// URI)",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		location := args[0]

		var data []byte
		var err error

		if bucket, key, ok := output.ParseS3URI(location); ok {
			region := sendSMSAWSRegionFlag
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

		if sendSMSNameFlag != "" {
			exchanges, err = exchange.FilterByGiverName(exchanges, sendSMSNameFlag)
			if err != nil {
				return err
			}
		}

		from := os.Getenv("TWILIO_PHONE_NUMBER")

		var client *twilio.RestClient
		if sendSMSDryRunFlag {
			fmt.Println("===\nDry run mode enabled. No sms messages will be sent.\n===")
		} else {
			if os.Getenv("TWILIO_ACCOUNT_SID") == "" || os.Getenv("TWILIO_AUTH_TOKEN") == "" || from == "" {
				return fmt.Errorf("TWILIO_ACCOUNT_SID, TWILIO_AUTH_TOKEN, and TWILIO_PHONE_NUMBER must all be set to send real SMS messages")
			}
			client = sms.NewClient()
		}

		total := len(exchanges)
		sent := 0
		for i, ex := range exchanges {
			body := sms.BuildMessage(ex.Giver.Name, ex.Receiver.Name)
			if sendSMSTestFlag {
				body = "TEST: " + body
			}

			if sendSMSDryRunFlag {
				if sendSMSQuietFlag {
					fmt.Printf("[dry run] Would send message %d/%d\n", i+1, total)
				} else {
					fmt.Printf("[dry run] Would send to %s (%s): %s\n", ex.Giver.Name, sms.FormatPhoneNumber(ex.Giver.Number), body)
				}
				sent++
				continue
			}

			sid, status, err := sms.Send(client, from, ex.Giver.Number, body)
			if err != nil {
				if sendSMSQuietFlag {
					return fmt.Errorf("failed to send SMS %d/%d: %w", i+1, total, err)
				}
				return fmt.Errorf("failed to send SMS to %s: %w", ex.Giver.Name, err)
			}

			if sendSMSQuietFlag {
				fmt.Printf("Sent message %d/%d (sid=%s, status=%s)\n", i+1, total, sid, status)
			} else {
				fmt.Printf("Sent to %s (%s) [sid=%s, status=%s]\n", ex.Giver.Name, sms.FormatPhoneNumber(ex.Giver.Number), sid, status)
			}
			sent++
		}

		if sendSMSDryRunFlag {
			fmt.Printf("Simulated %d SMS message(s)\n", sent)
		} else {
			fmt.Printf("Sent %d SMS message(s)\n", sent)
		}

		return nil
	},
}

func init() {
	sendSMSCmd.Flags().StringVar(&sendSMSAWSRegionFlag, "aws-region", "", "AWS region to use when reading an s3:// location (falls back to AWS_REGION env var)")
	sendSMSCmd.Flags().StringVar(&sendSMSNameFlag, "name", "", "If set, only (re)send the SMS to this participant (case-insensitive match on giver name)")
	sendSMSCmd.Flags().BoolVar(&sendSMSDryRunFlag, "dry-run", true, "If set, simulate sending without calling Twilio")
	sendSMSCmd.Flags().BoolVar(&sendSMSQuietFlag, "quiet", false, "Suppress participant names/numbers in output, printing message indices instead (for shared/CI logs)")
	sendSMSCmd.Flags().BoolVar(&sendSMSTestFlag, "test", false, "Prefix the SMS message body with 'TEST: ' to make test sends obviously distinguishable")

	rootCmd.AddCommand(sendSMSCmd)
}
