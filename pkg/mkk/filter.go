package mkk

import (
	"time"

	"github.com/pkg/errors"

	"github.com/mackerelio/mackerel-client-go"
)

// Filter implements Apply method which filters out given Mackerel hosts
type Filter interface {
	Apply(*mackerel.Client, []*mackerel.Host) ([]*mackerel.Host, error)
}

// GracePeriodFilter sets a grace period in second
// and filters out hosts which created within the period
type GracePeriodFilter struct {
	Seconds int64
}

// HostFilter selects the hosts with specified attribute
// HostFilter is useful when people want to utilize host attributes
// which cannot specify in mackerel.FindHostsParam
type HostFilter struct {
	Type string
}

// MetricAbsenceFilter selects hosts which does not report
// the specified metric within the given time period
type MetricAbsenceFilter struct {
	Name string
	From int64
	To   int64
}

// Apply applies GracePeriodFilter to the given hosts
func (f *GracePeriodFilter) Apply(_ *mackerel.Client, hosts []*mackerel.Host) ([]*mackerel.Host, error) {
	var filtered []*mackerel.Host

	for _, host := range hosts {
		if int64(host.CreatedAt) < time.Now().Unix()-f.Seconds {
			filtered = append(filtered, host)
		}
	}

	return filtered, nil
}

// Apply applies HostFilter to the given hosts
func (f *HostFilter) Apply(_ *mackerel.Client, hosts []*mackerel.Host) ([]*mackerel.Host, error) {
	var filtered []*mackerel.Host

	for _, host := range hosts {
		if host.Type == f.Type {
			filtered = append(filtered, host)
		}
	}

	return filtered, nil
}

// Apply applies MetricAbsenceFilter to the given hosts
func (f *MetricAbsenceFilter) Apply(m *mackerel.Client, hosts []*mackerel.Host) ([]*mackerel.Host, error) {
	if f.To == 0 {
		f.To = time.Now().Unix()
	}

	var filtered []*mackerel.Host

	for _, host := range hosts {
		time.Sleep(2 * time.Millisecond)

		values, err := m.FetchHostMetricValues(host.ID, f.Name, f.From, f.To)

		if err != nil {
			return nil, errors.Wrapf(err, "MetricAbsenceFilter.Apply fails while applying a filter: host: id: %v, name: %v, metric: %v", host.ID, host.Name, f.Name)
		}

		if len(values) == 0 {
			filtered = append(filtered, host)
		}
	}

	return filtered, nil
}
