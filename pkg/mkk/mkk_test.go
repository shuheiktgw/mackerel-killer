package mkk

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/shuheiktgw/mackerel-killer/test/until"

	"github.com/mackerelio/mackerel-client-go"
)

func TestMkk_FindHosts(t *testing.T) {
	hostName := "mackerel-killer-host"
	id := "abcdefg"

	var cases = []struct {
		title          string
		hostsResponse  string
		metricStatus   int
		metricResponse string
		error          bool
		want           int
	}{
		{
			title:          "Host does not exist",
			hostsResponse:  `{"hosts": []}`,
			metricStatus:   http.StatusOK,
			metricResponse: `{"metrics": [{"time":1,"value":"100"}, {"time":2,"value":"100"}]}`,
			want:           0,
		},
		{
			title:          "Invalid metric name",
			hostsResponse:  fmt.Sprintf(`{"hosts": [{"id":"%s"}]}`, id),
			metricStatus:   http.StatusNotFound,
			metricResponse: `{"message": "metric not found"}`,
			error:          true,
			want:           0,
		},
		{
			title:          "Metric does not exist",
			hostsResponse:  fmt.Sprintf(`{"hosts": [{"id":"%s"}]}`, id),
			metricStatus:   http.StatusOK,
			metricResponse: `{"metrics": []}`,
			want:           0,
		},
		{
			title:          "Metric exists",
			hostsResponse:  fmt.Sprintf(`{"hosts": [{"id":"%s"}]}`, id),
			metricStatus:   http.StatusOK,
			metricResponse: `{"metrics": [{"time":1,"value":"100"}, {"time":2,"value":"100"}]}`,
			want:           1,
		},
	}

	for i, tc := range cases {
		t.Run(tc.title, func(t *testing.T) {
			m, mux, _, teardown := setup()
			defer teardown()

			mux.HandleFunc("/api/v0/hosts", func(w http.ResponseWriter, r *http.Request) {
				util.TestMethod(t, r, http.MethodGet)
				util.TestFormValues(t, r, util.Values{"name": hostName})
				fmt.Fprint(w, tc.hostsResponse)
			})

			mux.HandleFunc(fmt.Sprintf("/api/v0/hosts/%s/metrics", id), func(w http.ResponseWriter, r *http.Request) {
				util.TestMethod(t, r, http.MethodGet)
				util.TestFormValues(t, r, util.Values{"from": "0", "to": "100", "name": "test"})
				w.WriteHeader(tc.metricStatus)
				fmt.Fprint(w, tc.metricResponse)
			})

			filters := []Filter{
				&MetricAbsenceFilter{
					Name: "test",
					From: 0,
					To:   100,
				},
			}

			param := mackerel.FindHostsParam{Name: hostName}
			hosts, err := m.FindHosts(&param, filters)

			if tc.error {
				if err == nil {
					t.Errorf("#%d error is not supposed to be nil", i)
				}
			} else {
				if err != nil {
					t.Errorf("#%d Mkk.FindHosts returned error: %v", i, err)
				}

				if got, want := len(hosts), tc.want; got != want {
					t.Errorf("#%d invalid number of hosts: got: %v, want: %v", i, got, want)
				}
			}
		})
	}
}
