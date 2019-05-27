package mkk

import (
	"github.com/mackerelio/mackerel-client-go"
	"github.com/pkg/errors"
)

type Mkk struct {
	Client *mackerel.Client
}

func NewMkk(token string) *Mkk {
	return &Mkk{Client: mackerel.NewClient(token)}
}

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

func (m *Mkk) Kill(host *mackerel.Host) error {
	return m.Client.RetireHost(host.ID)
}
