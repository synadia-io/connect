package cli

import (
	"fmt"
	"time"

	"github.com/synadia-io/connect/model"
)

func printLogs(logs []*model.InstanceLog) {
	if len(logs) == 0 {
		fmt.Print("No logs found\n")
		return
	}

	for _, log := range logs {
		if log == nil {
			continue
		}

		fmt.Printf("[%v %v %v] %v\n", log.ConnectorId, log.DeploymentId, log.InstanceId, log.Line)
	}
}

func printEvents(events []*model.InstanceEvent) {
	if len(events) == 0 {
		fmt.Print("No events found\n")
		return
	}

	// Check if any events have a local timezone
	var localTz *time.Location
	for _, event := range events {
		loc := event.Timestamp.Location()
		if loc != time.UTC {
			localTz = loc
			break
		}
	}
	if localTz == nil {
		localTz = time.UTC
	}

	for _, event := range events {
		if event == nil {
			continue
		}

		// Show all events in the same timezone, even if some are local and some are UTC
		ts := event.Timestamp.In(localTz)
		msg := fmt.Sprintf("[%v %v %v] %v %v", event.ConnectorId, event.DeploymentId, event.InstanceId, ts.Format(time.RFC3339), event.Type)
		if event.ExitCode != 0 {
			msg = fmt.Sprintf("%v %v %v", msg, event.ExitCode, event.Error)
		}

		fmt.Printf("%v\n", msg)
	}
}

func printMetrics(metrics []*model.InstanceMetric) {
	if len(metrics) == 0 {
		fmt.Print("No metrics found\n")
		return
	}

	for _, metric := range metrics {
		if metric == nil {
			continue
		}

		fmt.Printf("[%v %v %v] %v\n", metric.ConnectorId, metric.DeploymentId, metric.InstanceId, metric.Timestamp.Format(time.RFC3339))
		fmt.Printf("%v\n", string(metric.Data)) // print the raw prometheus metric data
	}
}
