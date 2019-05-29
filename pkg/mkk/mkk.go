package mkk

import (
	"github.com/mackerelio/mackerel-client-go"
	"github.com/pkg/errors"
)

// Mkk is a wrapper for mackerel.Client to retire the inactive Mackerel hosts
type Mkk struct {
	Client *mackerel.Client
}

// NewMkk initializes Mkk
func NewMkk(token string) *Mkk {
	return &Mkk{Client: mackerel.NewClient(token)}
}

// FindHosts finds hosts with mackerel.FindHostsParam and given filters
func (m *Mkk) FindHosts(param *mackerel.FindHostsParam, filters []Filter) ([]*mackerel.Host, error) {
	hosts, err := m.Client.FindHosts(param)
	if err != nil {
		return nil, errors.Wrap(err, "Mkk.FindHosts fails while finding hosts")
	}

	for _, f := range filters {
		hosts, err = f.Apply(m.Client, hosts)
		if err != nil {
			return nil, errors.Wrap(err, "Mkk.FindHosts fails while applying filters")
		}
	}

	return hosts, nil
}

// Kill retires specified Mackerel host
func (m *Mkk) Kill(host *mackerel.Host) error {
	return m.Client.RetireHost(host.ID)
}
