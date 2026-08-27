package main

import (
	ocireport "ocicomputecapacityreport"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	log "github.com/sirupsen/logrus"
)

var ociCapReport = cobra.Command{
	Use:     "oci-compute-capacity-report",
	Version: ocireport.VERSION,
	Short:   "OCI Compute Capacity Report",
	Long:    "OCI Compute Capacity Report is a tool to generate reports on OCI compute capacity.",

	Run: func(cmd *cobra.Command, args []string) {

		log.SetLevel(log.InfoLevel)
		if viper.GetBool("debug") {
			log.SetLevel(log.DebugLevel)
			log.Debug("Debug logging enabled.")
			viper.Debug()
		}
		log.Debugf("Shape: %s", viper.GetString("shape"))
		log.Debugf("OCPU: %d", viper.GetInt("ocpu"))
		log.Debugf("Memory: %d", viper.GetInt("memory"))
		log.Debugf("NodeClass: %s", viper.GetString("nodeclass"))
		log.Debugf("Profile: %s", viper.GetString("profile"))
		log.Debugf("Availability Domain: %s", viper.GetString("availabilitydomain"))

		// Require Shape, OCPU and Memory, or NodeClass to be specified
		if viper.GetString("nodeclass") == "" && (viper.GetString("shape") == "" || viper.GetInt("ocpu") == 0 || viper.GetInt("memory") == 0) {
			log.Error("Either NodeClass must be specified or Shape, OCPU and Memory must all be specified.")
			cmd.Usage()
			return
		}

		if viper.GetString("nodeclass") != "" {
			shape := ocireport.GetNodePoolShape(viper.GetString("nodeclass"))
			viper.Set("shape", shape.Shape)
			viper.Set("ocpu", shape.OCPU)
			viper.Set("memory", shape.Memory)
		}

		report := ocireport.NewReport(
			viper.GetString("shape"),
			viper.GetInt("ocpu"),
			viper.GetInt("memory"),
			viper.GetString("nodeclass"),
			viper.GetString("profile"),
			viper.GetString("availabilitydomain"),
		)

		report.GetData()

		report.DisplayAvailability()

	},
}

func init() {

	// Shape
	ociCapReport.Flags().String("shape", "", "OCI compute shape ")
	viper.BindPFlag("shape", ociCapReport.Flags().Lookup("shape"))

	// OCPU
	ociCapReport.Flags().Int("ocpu", 0, "OCI compute OCPU count to filter by")
	viper.BindPFlag("ocpu", ociCapReport.Flags().Lookup("ocpu"))

	// Memory
	ociCapReport.Flags().Int("memory", 0, "OCI compute memory (in GB) to filter by")
	viper.BindPFlag("memory", ociCapReport.Flags().Lookup("memory"))

	// NodeClass
	ociCapReport.Flags().String("nodeclass", "", "OCI compute node class to filter by")
	viper.BindPFlag("nodeclass", ociCapReport.Flags().Lookup("nodeclass"))

	// Availability Domain
	ociCapReport.Flags().String("availabilitydomain", "", "OCI availability domain to filter by")
	viper.BindPFlag("availabilitydomain", ociCapReport.Flags().Lookup("availabilitydomain"))

	// Command Argument: oci_profile
	ociCapReport.Flags().String("profile", "DEFAULT", "OCI profile to use")
	viper.BindPFlag("profile", ociCapReport.Flags().Lookup("profile"))

	// flag to enable debug logging
	ociCapReport.Flags().BoolP("debug", "d", false, "Enable debug logging.")
	viper.BindPFlag("debug", ociCapReport.Flags().Lookup("debug"))

}

func main() {
	ociCapReport.Execute()
}
