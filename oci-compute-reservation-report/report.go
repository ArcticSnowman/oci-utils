package ocicomputereservationreport

import (
	"fmt"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/sanity-io/litter"
	log "github.com/sirupsen/logrus"
)

type Reserved map[string]int64

type Shapes struct {
	Name     string
	OCPU     float32
	Memory   float32
	Capacity Reserved
}

type Reservation struct {
	Id                 string
	ReservationName    string
	AvailabilityDomain string
	Shapes             []Shapes
}

type Report struct {
	CompartmentName string
	CompartmentID   string
	ReservationID   string
	Reservations    []Reservation
}

func NewReport(compartment string) Report {

	var newReport Report

	newReport.CompartmentName = compartment
	newReport.CompartmentID = getCompartmentID(compartment)

	return newReport
}

func (r *Report) PopulateReservations() {
	reservations := getComputeReservations(r.CompartmentID)
	for _, res := range reservations {
		Id := *res.Id
		ReservationName := *res.DisplayName
		AvailabilityDomain := *res.AvailabilityDomain

		instances := getComputeReservationInstances(Id)
		shapes := map[string]Shapes{}
		shapekey := ""

		for _, instance := range instances {
			shapekey = fmt.Sprintf("%s-%d-%d", *instance.Shape, *instance.ShapeConfig.Ocpus, *instance.ShapeConfig.MemoryInGBs)
			if _, ok := shapes[shapekey]; !ok {
				shapes[shapekey] = Shapes{
					Name:     *instance.Shape,
					OCPU:     *instance.ShapeConfig.Ocpus,
					Memory:   *instance.ShapeConfig.MemoryInGBs,
					Capacity: make(Reserved),
				}
			}

			fd := *instance.FaultDomain
			shapes[shapekey].Capacity[fd]++
		}

		r.Reservations = append(r.Reservations, Reservation{
			Id:                 Id,
			ReservationName:    ReservationName,
			AvailabilityDomain: AvailabilityDomain,
			Shapes: func() []Shapes {
				shapesSlice := []Shapes{}
				for _, s := range shapes {
					shapesSlice = append(shapesSlice, s)
				}
				return shapesSlice
			}(),
		})

	}
}

func (r Report) ReportReservations() {

	log.Debugf("Report %s", litter.Sdump(r))

	if len(r.Reservations) == 0 {
		log.Infof("No reservations found for compartment %s", r.CompartmentName)
		return
	}

	//Loop Reservations
	for _, res := range r.Reservations {
		t := table.NewWriter()

		mergeCols := []table.ColumnConfig{{Name: "Shape", AutoMerge: true}}
		mergeCols = append(mergeCols, table.ColumnConfig{Name: "OCPU", AutoMerge: true})
		mergeCols = append(mergeCols, table.ColumnConfig{Name: "Memory", AutoMerge: true})
		mergeCols = append(mergeCols, table.ColumnConfig{Name: "Availability Domain", AutoMerge: true})
		t.SetColumnConfigs(mergeCols)
		t.SetOutputMirror(log.StandardLogger().Out)
		t.SetTitle(fmt.Sprintf("Reservation %s for Compartment %s", res.ReservationName, r.CompartmentName))
		t.AppendHeader(table.Row{"Shape", "OCPU", "Memory", "Availability Domain", "Fault Domain", "Reservation Count"})
		for _, shape := range res.Shapes {
			for fd, count := range shape.Capacity {
				t.AppendRow(table.Row{shape.Name, shape.OCPU, shape.Memory, res.AvailabilityDomain, fd, count})
			}
		}
		t.Render()
	}
}
