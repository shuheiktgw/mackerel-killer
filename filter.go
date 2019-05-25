package mkk

import (
	"time"

	"github.com/pkg/errors"

	"github.com/mackerelio/mackerel-client-go"
)

type Filter interface {
	Apply(*mackerel.Client, []*mackerel.Host) ([]*mackerel.Host, error)
}

type MetricExistenceFilter struct {
	Name string
	From *time.Time
	To   *time.Time
}

func (f *MetricExistenceFilter) Apply(m *mackerel.Client, hosts []*mackerel.Host) ([]*mackerel.Host, error) {
	var filtered []*mackerel.Host

	for _, host := range hosts {
		time.Sleep(2 * time.Millisecond)

		values, err := m.FetchHostMetricValues(host.ID, f.Name, f.From.Unix(), f.To.Unix())

		if err != nil {
			return nil, errors.Wrap(err, "MetricExistenceFilter.Apply fails while fetching a metric")
		}

		if len(values) != 0 {
			filtered = append(filtered, host)
		}
	}

	return filtered, nil
}
