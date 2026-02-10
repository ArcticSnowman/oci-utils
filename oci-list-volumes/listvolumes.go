package ocilistvolumes

import (
	"context"
	"fmt"
	"os"

	log "github.com/sirupsen/logrus"

	"github.com/oracle/oci-go-sdk/v65/common"
	"github.com/oracle/oci-go-sdk/v65/core"
	"github.com/oracle/oci-go-sdk/v65/identity"
	"github.com/spf13/viper"

	"github.com/jedib0t/go-pretty/v6/table"
)

func initOCICOnfig() common.ConfigurationProvider {

	configProvider := common.DefaultConfigProvider()

	return configProvider
}

var ociConfig common.ConfigurationProvider

func GetOCIConfig() common.ConfigurationProvider {
	if ociConfig == nil {
		ociConfig = initOCICOnfig()
	}

	return ociConfig
}

func GetCompartment(compartment string) string {

	client, err := identity.NewIdentityClientWithConfigurationProvider(GetOCIConfig())
	if err != nil {
		log.Fatalf("Error creating OCI Identity Client: %v", err)
	}

	log.Printf("Looking up OCI for compartment %s", viper.GetString("compartment"))

	tenancyOCID, err := ociConfig.TenancyOCID()
	if err != nil {
		log.Fatalf("Error getting Tenancy OCID: %v", err)
	}

	if viper.GetBool("debug") {
		log.Printf("Tenancy OCID: %s\n", tenancyOCID)
	}

	listResponse, err := client.ListCompartments(context.Background(), identity.ListCompartmentsRequest{
		CompartmentId:          common.String(tenancyOCID),
		CompartmentIdInSubtree: common.Bool(true),
	})
	if err != nil {
		log.Fatalf("Error listing compartments: %v", err)
	}

	for _, compartmentItem := range listResponse.Items {
		// if viper.GetBool("debug") {
		// 	log.Printf("Found Compartment OCID: %s for Name: %s\n", *compartmentItem.Id, *compartmentItem.Name)
		// }
		if *compartmentItem.Name == compartment {
			if viper.GetBool("debug") {
				log.Printf("Found Compartment OCID: %s for Name: %s\n", *compartmentItem.Id, *compartmentItem.Name)
			}
			return *compartmentItem.Id
		}
	}

	return ""
}

func ListVolumes() error {
	// Implement volume listing logic here

	CompartmentOCID := GetCompartment(viper.GetString("compartment"))
	if CompartmentOCID == "" {
		return fmt.Errorf("Compartment %s not found in tenancy.", viper.GetString("compartment"))
	}

	log.Infof("Listing volumes in compartment: %s", viper.GetString("compartment"))

	if viper.GetBool("boot") {
		bootVolumes, err := listBootVolumes(CompartmentOCID)
		if err != nil {
			return fmt.Errorf("Error listing boot volumes: %v", err)
		}
		attachments, err := listVolumeAttachments(CompartmentOCID)
		if err != nil {
			return fmt.Errorf("Error listing volume attachments: %v", err)
		}
		attachedList := genAttachedList(attachments)
		outputBootVolumes(bootVolumes, attachedList)

	} else {
		volumes, err := listBlockVolumes(CompartmentOCID)
		if err != nil {
			return fmt.Errorf("Error listing block volumes: %v", err)
		}
		attachments, err := listVolumeAttachments(CompartmentOCID)
		if err != nil {
			return fmt.Errorf("Error listing volume attachments: %v", err)
		}
		attachedList := genAttachedList(attachments)

		outputVolumes(volumes, attachedList)
	}
	return nil
}

// Prepare Block Storage Client
func prepareBlockStorageClient() (*core.BlockstorageClient, error) {
	client, err := core.NewBlockstorageClientWithConfigurationProvider(GetOCIConfig())
	if err != nil {
		return nil, fmt.Errorf("Error creating OCI Block Storage Client: %v", err)
	}
	return &client, nil
}

// Prepare Compute Client
func prepareComputeClient() (*core.ComputeClient, error) {
	client, err := core.NewComputeClientWithConfigurationProvider(GetOCIConfig())
	if err != nil {
		return nil, fmt.Errorf("Error creating OCI Compute Client: %v", err)
	}
	return &client, nil
}

// List Block Volumes
func listBlockVolumes(compartmentOCID string) ([]core.Volume, error) {

	var Volumes []core.Volume = make([]core.Volume, 0)
	client, err := prepareBlockStorageClient()
	if err != nil {
		return Volumes, err
	}

	var page string
	for {
		listVolResponse, err := client.ListVolumes(context.Background(), core.ListVolumesRequest{
			CompartmentId: common.String(compartmentOCID),
			Page:          common.String(page),
		})
		if err != nil {
			return Volumes, fmt.Errorf("Error listing block volumes: %v", err)
		}
		Volumes = append(Volumes, listVolResponse.Items...)

		// Next Page
		if listVolResponse.OpcNextPage != nil {
			page = *listVolResponse.OpcNextPage
		} else {
			break
		}
	}

	log.Debugf("%d block volumes found in compartment %s", len(Volumes), viper.GetString("compartment"))

	// Implement logic to list block volumes using the client
	return Volumes, nil
}

func listVolumeAttachments(compartmentOCID string) ([]core.VolumeAttachment, error) {

	var VolumeAttachments []core.VolumeAttachment = make([]core.VolumeAttachment, 0)
	client, err := prepareComputeClient()
	if err != nil {
		return VolumeAttachments, err
	}

	var page string
	for {
		listVolAttResponse, err := client.ListVolumeAttachments(context.Background(), core.ListVolumeAttachmentsRequest{
			CompartmentId: common.String(compartmentOCID),
			Page:          common.String(page),
		})
		if err != nil {
			return VolumeAttachments, fmt.Errorf("Error listing volume attachments: %v", err)
		}
		VolumeAttachments = append(VolumeAttachments, listVolAttResponse.Items...)

		// Next Page
		if listVolAttResponse.OpcNextPage != nil {
			page = *listVolAttResponse.OpcNextPage
		} else {
			break
		}
	}

	log.Debugf("%d volume attachments found in compartment %s", len(VolumeAttachments), viper.GetString("compartment"))

	return VolumeAttachments, nil
}

func genAttachedList(attachments []core.VolumeAttachment) map[string]bool {
	attachedList := make(map[string]bool)
	for _, attachment := range attachments {
		attachedList[*attachment.GetVolumeId()] = true
	}
	return attachedList
}

func outputVolumes(volumes []core.Volume, volAttached map[string]bool) {
	writer := table.NewWriter()
	writer.SetOutputMirror(os.Stdout)
	writer.AppendHeader(table.Row{"Name", "Id", "Size", "State", "Attached"})

	totalSize := 0
	for _, volume := range volumes {
		if viper.GetBool("unattached") && volAttached[*volume.Id] {
			continue
		}

		if viper.GetBool("available") && volume.LifecycleState != core.VolumeLifecycleStateAvailable {
			continue
		}

		totalSize += int(*volume.SizeInGBs)
		if !volAttached[*volume.Id] {
			writer.AppendRow(table.Row{*volume.DisplayName, *volume.Id, *volume.SizeInGBs, volume.LifecycleState, "No"})
		} else {
			writer.AppendRow(table.Row{*volume.DisplayName, *volume.Id, *volume.SizeInGBs, volume.LifecycleState, "Yes"})
		}
	}

	writer.Render()
	log.Infof("Total Volumes: %d, Total Size: %d GB", len(volumes), totalSize)
}

// List Boot Volumes
func listBootVolumes(compartmentOCID string) ([]core.BootVolume, error) {

	var BootVolumes []core.BootVolume = make([]core.BootVolume, 0)
	client, err := prepareBlockStorageClient()
	if err != nil {
		return BootVolumes, err
	}

	var page string
	for {
		listBootVolResponse, err := client.ListBootVolumes(context.Background(), core.ListBootVolumesRequest{
			CompartmentId: common.String(compartmentOCID),
			Page:          common.String(page),
		})
		if err != nil {
			return BootVolumes, fmt.Errorf("Error listing boot volumes: %v", err)
		}
		BootVolumes = append(BootVolumes, listBootVolResponse.Items...)

		// Next Page
		if listBootVolResponse.OpcNextPage != nil {
			page = *listBootVolResponse.OpcNextPage
		} else {
			break
		}
	}

	log.Debugf("%d boot volumes found in compartment %s", len(BootVolumes), viper.GetString("compartment"))

	// Implement logic to list boot volumes using the client
	return BootVolumes, nil
}

func outputBootVolumes(volumes []core.BootVolume, volAttached map[string]bool) {
	writer := table.NewWriter()
	writer.SetOutputMirror(os.Stdout)
	writer.AppendHeader(table.Row{"Name", "Id", "Size", "State", "Attached"})

	totalSize := 0
	for _, volume := range volumes {
		if viper.GetBool("unattached") && volAttached[*volume.Id] {
			continue
		}

		if viper.GetBool("available") && volume.LifecycleState != core.BootVolumeLifecycleStateAvailable {
			continue
		}

		totalSize += int(*volume.SizeInGBs)
		if !volAttached[*volume.Id] {
			writer.AppendRow(table.Row{*volume.DisplayName, *volume.Id, *volume.SizeInGBs, volume.LifecycleState, "No"})
		} else {
			writer.AppendRow(table.Row{*volume.DisplayName, *volume.Id, *volume.SizeInGBs, volume.LifecycleState, "Yes"})
		}
	}
	writer.Render()
	log.Infof("Total Boot Volumes: %d, Total Size: %d GB", len(volumes), totalSize)
}
