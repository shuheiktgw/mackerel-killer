// +build integration

package mkk

import (
	"fmt"
	"testing"
	"time"

	"github.com/mackerelio/mackerel-client-go"
)

func TestMetricExistenceFilter_Integration__Apply(t *testing.T) {
	now := time.Now().Unix()
	hostName := fmt.Sprintf("mackerel-killer-host-%v", now)

	// Create a host
	chp := mackerel.CreateHostParam{Name: hostName}
	id, err := integrationMkk.Client.CreateHost(&chp)
	if err != nil {
		t.Fatalf("error occurred while creating service: %v", err)
	}

	defer integrationMkk.Client.RetireHost(id)

	// Post a metric to the host
	values := []*mackerel.HostMetricValue{
		{
			HostID: id,
			MetricValue: &mackerel.MetricValue{
				Name:  "mackerel-killer-custom",
				Time:  now,
				Value: 100,
			},
		},
	}
	err = integrationMkk.Client.PostHostMetricValues(values)
	if err != nil {
		t.Fatalf("error occurred while posting a metric: %v", err)
	}

	// Wait until the metric is available via the API
	time.Sleep(20 * time.Second)

	// List hosts
	fhp := mackerel.FindHostsParam{Name: hostName}
	hosts, err := integrationMkk.Client.FindHosts(&fhp)
	if err != nil {
		t.Fatalf("error occurred while getting hosts: %v", err)
	}

	var cases = []struct {
		title string
		name  string
		from  time.Time
		to    time.Time
		error bool
		want  int
	}{
		{
			title: "Metric name does not match",
			name:  "unknown-metric",
			from:  time.Unix(now-10, 0),
			to:    time.Unix(now+10, 0),
			error: true,
		},
		{
			title: "Does not include the metric (upper bounds)",
			name:  "mackerel-killer-custom",
			from:  time.Unix(now-10, 0),
			to:    time.Unix(now-2, 0),
			want:  0,
		},
		{
			title: "Does not include the metric (lower bounds)",
			name:  "mackerel-killer-custom",
			from:  time.Unix(now+2, 0),
			to:    time.Unix(now+10, 0),
			want:  0,
		},
		{
			title: "Includes the metric",
			name:  "mackerel-killer-custom",
			from:  time.Unix(now-100, 0),
			to:    time.Unix(now+100, 0),
			want:  1,
		},
	}

	for i, tc := range cases {
		t.Run(tc.title, func(t *testing.T) {
			f := MetricExistenceFilter{Name: tc.name, From: &tc.from, To: &tc.to}
			filtered, err := f.Apply(integrationMkk.Client, hosts)

			if tc.error {
				if err == nil {
					t.Fatalf("#%v is suppoed to return error", i)
				}
			} else {
				if err != nil {
					t.Fatalf("#%v error occurred while applying a filter: %v", i, err)
				}

				if got, want := len(filtered), tc.want; got != want {
					t.Fatalf("#%v invalid number of hosts: got: %v, want: %v", i, got, want)
				}
			}
		})
	}

}
