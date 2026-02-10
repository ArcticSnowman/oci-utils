package main

import (
	ocilistvolumes "oci-list-volumes"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	log "github.com/sirupsen/logrus"
)

// Cobra based CLI tool to list OCI volumes in the current tenancy and region.

var Cmd = &cobra.Command{
	Use:     "oci-list-volumes",
	Short:   "List OCI volumes in selected compartment for current tenancy and region.",
	Version: ocilistvolumes.VERSION,
	RunE: func(cmd *cobra.Command, args []string) error {

		log.SetLevel(log.InfoLevel)
		if viper.GetBool("debug") {
			log.SetLevel(log.DebugLevel)
			log.Debug("Debug logging enabled.")
			log.Debugf("Compartment: %s", viper.GetString("compartment"))
			log.Debugf("Boot Volumes: %t", viper.GetBool("boot"))
			log.Debugf("Unattached Volumes Only: %t", viper.GetBool("unattached"))
		}

		return ocilistvolumes.ListVolumes()
	},
}

func init() {
	// Flag to select Compartment to list volumes in.
	Cmd.Flags().StringP("compartment", "c", "", "Compartment Name to list volumes in.")
	Cmd.MarkFlagRequired("compartment")
	viper.BindPFlag("compartment", Cmd.Flags().Lookup("compartment"))

	//Flag to List Boot Volumes
	Cmd.Flags().BoolP("boot", "b", false, "List Boot Volumes instead of Block Volumes.")
	viper.BindPFlag("boot", Cmd.Flags().Lookup("boot"))

	// Flag to list unattached volumes only.
	Cmd.Flags().BoolP("unattached", "u", false, "List unattached volumes only.")
	viper.BindPFlag("unattached", Cmd.Flags().Lookup("unattached"))

	// flag to list available volumes only.
	Cmd.Flags().BoolP("available", "a", false, "List available volumes only.")
	viper.BindPFlag("available", Cmd.Flags().Lookup("available"))

	// flag to enable debug logging
	Cmd.Flags().BoolP("debug", "d", false, "Enable debug logging.")
	viper.BindPFlag("debug", Cmd.Flags().Lookup("debug"))

}

func main() {

	if err := Cmd.Execute(); err != nil {
		panic(err)
	}
}
