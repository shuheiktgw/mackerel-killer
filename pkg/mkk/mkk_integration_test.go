// +build integration

package mkk

import (
	"fmt"
	"testing"
	"time"

	"github.com/mackerelio/mackerel-client-go"
)

func TestMkk_Integration_Apply(t *testing.T) {
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

	var cases = []struct {
		title   string
		filters []Filter
		error   bool
		want    int
	}{
		{
			title: "Metric name does not match",
			filters: []Filter{
				&MetricExistenceFilter{
					Name: "unknown-metric",
					From: time.Unix(now-100, 0).Unix(),
					To:   time.Unix(now+100, 0).Unix(),
				},
			},
			error: true,
		},
		{
			title: "First filter excludes the metric",
			filters: []Filter{
				&MetricExistenceFilter{
					Name: "mackerel-killer-custom",
					From: time.Unix(now-1000, 0).Unix(),
					To:   time.Unix(now-900, 0).Unix(),
				},
				&MetricExistenceFilter{
					Name: "mackerel-killer-custom",
					From: time.Unix(now-100, 0).Unix(),
					To:   time.Unix(now+100, 0).Unix(),
				},
			},
			want: 0,
		},
		{
			title: "Second filter excludes the metric",
			filters: []Filter{
				&MetricExistenceFilter{
					Name: "mackerel-killer-custom",
					From: time.Unix(now-100, 0).Unix(),
					To:   time.Unix(now+100, 0).Unix(),
				},
				&MetricExistenceFilter{
					Name: "mackerel-killer-custom",
					From: time.Unix(now-1000, 0).Unix(),
					To:   time.Unix(now-900, 0).Unix(),
				},
			},
			want: 0,
		},
		{
			title: "Metric is not excluded",
			filters: []Filter{
				&MetricExistenceFilter{
					Name: "mackerel-killer-custom",
					From: time.Unix(now-100, 0).Unix(),
					To:   time.Unix(now+100, 0).Unix(),
				},
				&MetricExistenceFilter{
					Name: "mackerel-killer-custom",
					From: time.Unix(now-50, 0).Unix(),
					To:   time.Unix(now+200, 0).Unix(),
				},
			},
			want: 1,
		},
	}

	for i, tc := range cases {
		t.Run(tc.title, func(t *testing.T) {
			params := mackerel.FindHostsParam{Name: hostName}
			hosts, err := integrationMkk.Apply(&params, tc.filters)

			if tc.error {
				if err == nil {
					t.Fatalf("#%v is suppoed to return error", i)
				}
			} else {
				if err != nil {
					t.Fatalf("#%v error occurred while applying filters: %v", i, err)
				}

				if got, want := len(hosts), tc.want; got != want {
					t.Fatalf("#%v invalid number of hosts: got: %v, want: %v", i, got, want)
				}
			}
		})
	}
}

func TestMkk_Integration_Kill(t *testing.T) {
	chp := mackerel.CreateHostParam{Name: "mackerel-killer-kill"}
	_, err := integrationMkk.Client.CreateHost(&chp)
	if err != nil {
		t.Fatalf("error occurred while creating a host: %v", err)
	}

	fhp := mackerel.FindHostsParam{Name: "mackerel-killer-kill"}
	hosts, err := integrationMkk.Client.FindHosts(&fhp)
	if err != nil {
		t.Fatalf("error occurred while finding hosts: %v", err)
	}

	err = integrationMkk.Kill(hosts)
	if err != nil {
		t.Fatalf("error occurred while retiring hosts: %v", err)
	}
}
