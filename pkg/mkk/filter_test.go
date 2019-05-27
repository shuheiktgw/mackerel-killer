package mkk

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/shuheiktgw/mackerel-killer/test/until"

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
			want:     1,
		},
		{
			title:    "Metric exists",
			response: `{"metrics": [{"time":1,"value":"100"}, {"time":2,"value":"100"}]}`,
			want:     0,
		},
	}

	for i, tc := range cases {
		t.Run(tc.title, func(t *testing.T) {
			m, mux, _, teardown := setup()
			defer teardown()

			id := "abcdefg"

			mux.HandleFunc(fmt.Sprintf("/api/v0/hosts/%s/metrics", id), func(w http.ResponseWriter, r *http.Request) {
				util.TestMethod(t, r, http.MethodGet)
				util.TestFormValues(t, r, util.Values{"from": "0", "to": "100", "name": "test"})
				fmt.Fprint(w, tc.response)
			})

			hosts := []*mackerel.Host{{ID: id}}

			filter := MetricAbsenceFilter{Name: "test", From: 0, To: 100}
			filtered, err := filter.Apply(m.Client, hosts)

			if err != nil {
				t.Errorf("#%d MetricAbsenceFilter.Apply returned error: %v", i, err)
			}

			if got, want := len(filtered), tc.want; got != want {
				t.Errorf("#%d invalid number of hosts: got: %v, want: %v", i, got, want)
			}
		})
	}
}
