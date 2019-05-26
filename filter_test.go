package mkk

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/mackerelio/mackerel-client-go"
)

func TestMetricExistenceFilter_Apply(t *testing.T) {
	var cases = []struct {
		title    string
		response string
		want     int
	}{
		{
			title:    "Metric does not exist",
			response: `{"metrics": []}`,
			want:     0,
		},
		{
			title:    "Metric exists",
			response: `{"metrics": [{"time":1,"value":"100"}, {"time":2,"value":"100"}]}`,
			want:     1,
		},
	}

	for i, tc := range cases {
		t.Run(tc.title, func(t *testing.T) {
			mkk, mux, _, teardown := setup()
			defer teardown()

			id := "abcdefg"

			mux.HandleFunc(fmt.Sprintf("/api/v0/hosts/%s/metrics", id), func(w http.ResponseWriter, r *http.Request) {
				testMethod(t, r, http.MethodGet)
				testFormValues(t, r, values{"from": "0", "to": "100", "name": "test"})
				fmt.Fprint(w, tc.response)
			})

			hosts := []*mackerel.Host{{ID: id}}

			filter := MetricExistenceFilter{Name: "test", From: 0, To: 100}
			filtered, err := filter.Apply(mkk.Client, hosts)

			if err != nil {
				t.Errorf("#%d MetricExistenceFilter.Apply returned error: %v", i, err)
			}

			if got, want := len(filtered), tc.want; got != want {
				t.Errorf("#%d invalid number of hosts: got: %v, want: %v", i, got, want)
			}
		})
	}
}
