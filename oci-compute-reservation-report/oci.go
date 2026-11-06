package ocicomputereservationreport

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
func getCompartmentID(compartmentName string) string {

	config := GetOCIConfig()

	idClient, err := identity.NewIdentityClientWithConfigurationProvider(config)
	if err != nil {
		log.Errorf("Failed to create Identity Client: %v", err)
		return ""
	}

	request := identity.ListCompartmentsRequest{

		CompartmentId:          common.String(getRootCompartmentID()),
		CompartmentIdInSubtree: common.Bool(true),
		Name:                   common.String(compartmentName),
	}

	response, err := idClient.ListCompartments(context.Background(), request)
	if err != nil {
		log.Errorf("Failed to list compartments: %v", err)
		return ""
	}

	log.Debugf("%s", litter.Sdump(response))

	for _, compartment := range response.Items {
		if *compartment.Name == compartmentName {
			return *compartment.Id
		}
	}

	log.Errorf("Compartment with name %s not found", compartmentName)
	return ""
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

func getComputeReservations(compartmentID string) []core.ComputeCapacityReservationSummary {
	config := GetOCIConfig()
	computeClient, err := core.NewComputeClientWithConfigurationProvider(config)
	if err != nil {
		log.Errorf("Failed to create Compute Client: %v", err)
		return nil
	}

	reservations := make([]core.ComputeCapacityReservationSummary, 0)
	var page *string
	for {
		request := core.ListComputeCapacityReservationsRequest{
			CompartmentId: common.String(compartmentID),
			Page:          page,
		}

		response, err := computeClient.ListComputeCapacityReservations(context.Background(), request)
		if err != nil {
			log.Errorf("Failed to list compute reservations for compartment %s: %v", compartmentID, err)
			return nil
		}

		reservations = append(reservations, response.Items...)
		if response.OpcNextPage == nil {
			break
		}
		page = response.OpcNextPage
	}

	return reservations
}

func getComputeReservationInstances(ReservationID string) []core.CapacityReservationInstanceSummary {
	config := GetOCIConfig()
	computeClient, err := core.NewComputeClientWithConfigurationProvider(config)
	if err != nil {
		log.Errorf("Failed to create Compute Client: %v", err)
		return nil
	}

	request := core.ListComputeCapacityReservationInstancesRequest{
		CapacityReservationId: common.String(ReservationID),
	}
	response, err := computeClient.ListComputeCapacityReservationInstances(context.Background(), request)
	if err != nil {
		log.Errorf("Failed to list compute reservation instances for reservation %s: %v", ReservationID, err)
		return nil
	}

	return response.Items
}

// func (r Report) GetOCIData() {

// 	log.Infof("Generating OCI Compute Capacity Report For %s", r.String())

// 	rootCompartmentID := getRootCompartmentID()
// 	if rootCompartmentID == "" {
// 		log.Errorf("Failed to get root compartment ID")
// 		return
// 	}

// 	var availabilityDomains []string

// 	if r.AvailabilityDomain != "" {
// 		availabilityDomains = []string{r.AvailabilityDomain}
// 		log.Debugf("Using specified Availability Domain: %s", r.AvailabilityDomain)
// 	} else {
// 		availabilityDomains = getAvailabilityDomains()
// 		if availabilityDomains == nil {
// 			log.Errorf("Failed to get availability domains")
// 			return
// 		}
// 	}

// 	log.Debugf("Tenancy ID: %s", rootCompartmentID)
// 	log.Debugf("Availability Domains: %v", availabilityDomains)

// 	config := GetOCIConfig()

// 	client, err := core.NewComputeClientWithConfigurationProvider(config)
// 	if err != nil {
// 		log.Errorf("Failed to create Compute Client: %v", err)
// 		return
// 	}

// 	ShapeAvailability := core.CreateCapacityReportShapeAvailabilityDetails{
// 		InstanceShape: common.String(r.Shape),
// 		InstanceShapeConfig: &core.CapacityReportInstanceShapeConfig{
// 			Ocpus:       common.Float32(float32(r.OCPU)),
// 			MemoryInGBs: common.Float32(float32(r.Memory)),
// 		},
// 	}

// 	for _, ad := range availabilityDomains {
// 		for _, fd := range getFaultDomains(ad) {

// 			log.Debugf("Requesting Availability For Domain: %s, Fault Domain: %s", ad, fd)

// 			ShapeAvailability.FaultDomain = common.String(fd)

// 			request := core.CreateComputeCapacityReportRequest{
// 				CreateComputeCapacityReportDetails: core.CreateComputeCapacityReportDetails{
// 					CompartmentId:       common.String(rootCompartmentID),
// 					AvailabilityDomain:  common.String(ad),
// 					ShapeAvailabilities: []core.CreateCapacityReportShapeAvailabilityDetails{ShapeAvailability},
// 				},
// 			}
// 			log.Debugf("Submitting Request: %v", request)

// 			resp, err := client.CreateComputeCapacityReport(context.Background(), request)
// 			if err != nil {
// 				log.Errorf("Failed to create Compute Capacity Report: %v", err)
// 				return
// 			}

// 			report := resp.ComputeCapacityReport
// 			log.Debugf("Received Response: %v", report)
// 			if len(report.ShapeAvailabilities) == 0 {
// 				log.Errorf("No shape availabilities returned in the response")
// 				return
// 			}

// 			log.Debugf("Availability: %v ", report.ShapeAvailabilities[0])
// 			r.SetAvailability(ad, fd, *report.ShapeAvailabilities[0].AvailableCount)

// 		}
// 	}
// }
