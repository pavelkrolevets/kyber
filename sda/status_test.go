package sda

import (
	"strings"
	"testing"

	"strconv"

	"github.com/stretchr/testify/assert"
)

func TestSRStruct(t *testing.T) {
	srs := newStatusReporterStruct()
	assert.NotNil(t, srs)
	dtr := &dummyTestReporter{5}
	srs.RegisterStatusReporter("Dummy", dtr)
	assert.Equal(t, srs.ReportStatus()["Dummy"]["Connections"], "5")
	dtr.Status = 10
	assert.Equal(t, srs.ReportStatus()["Dummy"]["Connections"], "10")
}

func TestStatusHost(t *testing.T) {
	c := NewTCPConode(2000)
	defer c.Close()
	stats := c.GetStatus()
	a := ServiceFactory.RegisteredServiceNames()
	services := strings.Split(stats["Available_Services"], ",")
	assert.Equal(t, len(services), len(a))
}

type dummyTestReporter struct {
	Status int
}

func (d *dummyTestReporter) GetStatus() Status {
	return Status{
		"Connections": strconv.Itoa(d.Status),
	}
}
