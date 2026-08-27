package ocicomputecapacityreport

import (
	"os/exec"

	"github.com/sanity-io/litter"
	log "github.com/sirupsen/logrus"

	"github.com/tidwall/gjson"
)

type OCIShape struct {
	Shape  string
	OCPU   int
	Memory int
}

func GetNodePoolShape(nodePoolName string) OCIShape {
	// Implementation to retrieve the shape of the node pool
	var shape OCIShape
	//kubectl get nodepools aok-ads-sfx -o json | jq '.spec.template.spec.requirements[] | select(.key=="oci.oraclecloud.com/instance-shape")| .values[]'
	cmd := exec.Command("kubectl", "get", "nodepools", nodePoolName, "-o", "json")
	output, err := cmd.Output()
	if err != nil {
		// Handle error
		return shape
	}

	// Parse the JSON output to extract the shape information
	// Example using jq: '.spec.template.spec.requirements[] | select(.key=="oci.oraclecloud.com/instance-shape")| .values[]'

	result := gjson.ParseBytes(output)
	shape.Shape = result.Get("spec.template.spec.requirements.#(key==\"oci.oraclecloud.com/instance-shape\").values.0").String()

	// Extract OCPU and Memory information based on the shape
	// kubectl get ocinodeclass aok-ads-sfx -o json | jq '.spec.shapeConfigs[]'
	cmd = exec.Command("kubectl", "get", "ocinodeclass", nodePoolName, "-o", "json")
	output, err = cmd.Output()
	if err != nil {
		// Handle error
		return shape
	}

	result = gjson.ParseBytes(output)
	shape.OCPU = int(result.Get("spec.shapeConfigs.0.ocpus").Int())
	shape.Memory = int(result.Get("spec.shapeConfigs.0.memoryInGbs").Int())

	log.Debugf("OCI Shape Retrieved: %s", litter.Sdump(shape))

	return shape
}

func newShape(shape string, ocpu int, memory int) OCIShape {
	return OCIShape{
		Shape:  shape,
		OCPU:   ocpu,
		Memory: memory,
	}
}
