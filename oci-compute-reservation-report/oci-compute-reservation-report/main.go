package main

import (
	"fmt"
	"ocicomputereservationreport"
	"os"

	"github.com/sanity-io/litter"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func newRootCommand() *cobra.Command {
	// config := viper.New()

	cmd := &cobra.Command{
		Use:           "oci-compute-reservation-report",
		Short:         "Generate OCI compute reservation reports",
		Long:          "OCI Compute Reservation Report is a tool to generate reports on OCI compute reservations.",
		Args:          cobra.NoArgs,
		Version:       ocicomputereservationreport.VERSION,
		SilenceUsage:  true,
		SilenceErrors: true,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			return viper.BindPFlags(cmd.Flags())
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			logger := log.New()
			logger.SetOutput(cmd.ErrOrStderr())
			if viper.GetBool("debug") {
				logger.SetLevel(log.DebugLevel)
				viper.Debug()
			}
			logger.Debug("Debug logging enabled")
			logger.Debugf("Compartment: %s", viper.GetString("compartment"))
			logger.Debugf("OCI profile: %s", viper.GetString("profile"))

			report := ocicomputereservationreport.NewReport(viper.GetString("compartment"))

			report.PopulateReservations()
			logger.Debugf("Report: %v", litter.Sdump(report))

			report.ReportReservations()
			return nil
		},
	}

	cmd.PersistentFlags().BoolP("debug", "d", false, "Enable debug logging")
	cmd.Flags().StringP("compartment", "c", "", "OCI compartment name")
	cmd.Flags().String("profile", "DEFAULT", "OCI profile to use")
	_ = cmd.MarkFlagRequired("compartment")
	return cmd
}

func main() {
	if err := newRootCommand().Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
