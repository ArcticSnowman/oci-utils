package ocicomputecapacityreport

import (
	"context"

	"github.com/oracle/oci-go-sdk/v65/common"
	"github.com/oracle/oci-go-sdk/v65/core"
	"github.com/oracle/oci-go-sdk/v65/identity"
	"github.com/sanity-io/litter"
	"github.com/spf13/viper"

	log "github.com/sirupsen/logrus"
)

func initOCICOnfig() common.ConfigurationProvider {

	configProvider := common.CustomProfileConfigProvider("", viper.GetString("profile"))

	return configProvider
}

var ociConfig common.ConfigurationProvider

func GetOCIConfig() common.ConfigurationProvider {
	if ociConfig == nil {
		ociConfig = initOCICOnfig()
	}

	return ociConfig
}

func getRootCompartmentID() string {

	config := GetOCIConfig()

	id, err := config.TenancyOCID()
	if err != nil {
		log.Errorf("Failed to get Tenancy OCID: %v", err)
		return ""
	}
	return id
}

func getAvailabilityDomains() []string {

	config := GetOCIConfig()

	idClient, err := identity.NewIdentityClientWithConfigurationProvider(config)
	if err != nil {
		log.Errorf("Failed to create Identity Client: %v", err)
		return nil
	}

	request := identity.ListAvailabilityDomainsRequest{
		CompartmentId: common.String(getRootCompartmentID()),
	}

	response, err := idClient.ListAvailabilityDomains(context.Background(), request)
	if err != nil {
		log.Errorf("Failed to list availability domains: %v", err)
		return nil
	}

	var availabilityDomains []string
	for _, ad := range response.Items {
		availabilityDomains = append(availabilityDomains, *ad.Name)
	}

	return availabilityDomains

}

func getFaultDomains(availabilityDomain string) []string {

	config := GetOCIConfig()

	idClient, err := identity.NewIdentityClientWithConfigurationProvider(config)
	if err != nil {
		log.Errorf("Failed to create Identity Client: %v", err)
		return nil
	}

	request := identity.ListFaultDomainsRequest{
		CompartmentId:      common.String(getRootCompartmentID()),
		AvailabilityDomain: common.String(availabilityDomain),
	}

	log.Debugf("Fault Domain Request: %s", litter.Sdump(request))

	response, err := idClient.ListFaultDomains(context.Background(), request)
	if err != nil {
		log.Errorf("Failed to list fault domains: %v", err)
		return nil
	}

	var faultDomains []string
	for _, fd := range response.Items {
		faultDomains = append(faultDomains, *fd.Name)
	}

	return faultDomains

}

func (r Report) GetOCIData() {

	log.Infof("Generating OCI Compute Capacity Report For %s", r.String())

	rootCompartmentID := getRootCompartmentID()
	if rootCompartmentID == "" {
		log.Errorf("Failed to get root compartment ID")
		return
	}

	var availabilityDomains []string

	if r.AvailabilityDomain != "" {
		availabilityDomains = []string{r.AvailabilityDomain}
		log.Debugf("Using specified Availability Domain: %s", r.AvailabilityDomain)
	} else {
		availabilityDomains = getAvailabilityDomains()
		if availabilityDomains == nil {
			log.Errorf("Failed to get availability domains")
			return
		}
	}

	log.Debugf("Tenancy ID: %s", rootCompartmentID)
	log.Debugf("Availability Domains: %v", availabilityDomains)

	config := GetOCIConfig()

	client, err := core.NewComputeClientWithConfigurationProvider(config)
	if err != nil {
		log.Errorf("Failed to create Compute Client: %v", err)
		return
	}

	ShapeAvailability := core.CreateCapacityReportShapeAvailabilityDetails{
		InstanceShape: common.String(r.Shape),
		InstanceShapeConfig: &core.CapacityReportInstanceShapeConfig{
			Ocpus:       common.Float32(float32(r.OCPU)),
			MemoryInGBs: common.Float32(float32(r.Memory)),
		},
	}

	for _, ad := range availabilityDomains {
		for _, fd := range getFaultDomains(ad) {

			log.Debugf("Requesting Availability For Domain: %s, Fault Domain: %s", ad, fd)

			ShapeAvailability.FaultDomain = common.String(fd)

			request := core.CreateComputeCapacityReportRequest{
				CreateComputeCapacityReportDetails: core.CreateComputeCapacityReportDetails{
					CompartmentId:       common.String(rootCompartmentID),
					AvailabilityDomain:  common.String(ad),
					ShapeAvailabilities: []core.CreateCapacityReportShapeAvailabilityDetails{ShapeAvailability},
				},
			}

			resp, err := client.CreateComputeCapacityReport(context.Background(), request)
			if err != nil {
				log.Errorf("Failed to create Compute Capacity Report: %v", err)
				return
			}

			report := resp.ComputeCapacityReport
			log.Debugf("Availability: %d ", *report.ShapeAvailabilities[0].AvailableCount)
			r.SetAvailability(ad, fd, *report.ShapeAvailabilities[0].AvailableCount)

		}
	}
}
