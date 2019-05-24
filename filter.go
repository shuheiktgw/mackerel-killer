package mkk

import (
	"github.com/mackerelio/mackerel-client-go"
	"net/http"
	"time"
)

type Filter interface {
	Apply(*Mkk, []*mackerel.Host) ([]*mackerel.Host, error)
}

type MetricExistenceFilter struct {
	Name string
	From *time.Time
	To   *time.Time
}

func (f *MetricExistenceFilter) Apply(m *Mkk, hosts []*mackerel.Host) ([]*mackerel.Host, error) {
	var filtered []*mackerel.Host

	for _, host := range hosts {
		time.Sleep(2 * time.Millisecond)

		_, err := m.Client.FetchHostMetricValues(host.ID, f.Name, f.From.Unix(), f.To.Unix())
		if err != nil {
			if e, ok := err.(*mackerel.APIError); ok && e.StatusCode == http.StatusNotFound {
				continue
			}

			return nil, err
		}

		filtered = append(filtered, host)
	}

	return filtered, nil
}
