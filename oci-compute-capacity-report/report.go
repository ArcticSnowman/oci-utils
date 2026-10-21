package ocicomputecapacityreport

import (
	"fmt"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/sanity-io/litter"
	log "github.com/sirupsen/logrus"
)

type Availability map[string]map[string]int64

type Report struct {
	Shape              string
	OCPU               int
	Memory             int
	NodeClass          string
	Profile            string
	AvailabilityDomain string
	Available          Availability
}

func NewReport(shape string, ocpu int, memory int, nodeClass string, profile string, availabilityDomain string) Report {
	return Report{
		Shape:              shape,
		OCPU:               ocpu,
		Memory:             memory,
		NodeClass:          nodeClass,
		Profile:            profile,
		AvailabilityDomain: availabilityDomain,
		Available:          make(Availability),
	}
}

func (r Report) SetAvailability(ad string, fd string, count int64) {
	if _, ok := r.Available[ad]; !ok {
		r.Available[ad] = make(map[string]int64)
	}
	r.Available[ad][fd] = count
}

func (r Report) String() string {
	if r.NodeClass != "" {
		return fmt.Sprintf("NodeClass %s in Profile: %s  Availability Domain: %s", r.NodeClass, r.Profile, r.AvailabilityDomain)
	} else {
		return fmt.Sprintf("Shape  %s (OCPU: %d, Memory: %d)  Profile: %s  Availability Domain: %s", r.Shape, r.OCPU, r.Memory, r.Profile, r.AvailabilityDomain)
	}
}

func (r Report) DisplayAvailability() {
	t := table.NewWriter()

	log.Debugf("Report %s", litter.Sdump(r))

	mergeCols := []table.ColumnConfig{{Name: "Shape", AutoMerge: true}}
	mergeCols = append(mergeCols, table.ColumnConfig{Name: "OCPU", AutoMerge: true})
	mergeCols = append(mergeCols, table.ColumnConfig{Name: "Memory", AutoMerge: true})
	mergeCols = append(mergeCols, table.ColumnConfig{Name: "Availability Domain", AutoMerge: true})
	t.SetColumnConfigs(mergeCols)
	t.SetOutputMirror(log.StandardLogger().Out)
	t.AppendHeader(table.Row{"Shape", "OCPU", "Memory", "Availability Domain", "Fault Domain", "Available Count"})
	for ad, fds := range r.Available {
		for fd, count := range fds {
			t.AppendRow(table.Row{r.Shape, r.OCPU, r.Memory, ad, fd, count})
		}
	}
	t.Render()
}
