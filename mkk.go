package mkk

import (
	"time"

	"github.com/mackerelio/mackerel-client-go"
	"github.com/pkg/errors"
)

type Mkk struct {
	Client *mackerel.Client
}

func NewMkk(token string) *Mkk {
	return &Mkk{Client: mackerel.NewClient(token)}
}

func (m *Mkk) Apply(params *mackerel.FindHostsParam, filters []Filter) ([]*mackerel.Host, error) {
	hosts, err := m.Client.FindHosts(params)
	if err != nil {
		return nil, errors.Wrap(err, "Mkk.Apply fails while finding hosts")
	}

	for _, f := range filters {
		hosts, err = f.Apply(m.Client, hosts)
		if err != nil {
			return nil, errors.Wrap(err, "Mkk.Apply fails while applying filters")
		}
	}

	return hosts, nil
}

func (m *Mkk) Kill(hosts []*mackerel.Host) error {
	for _, host := range hosts {
		time.Sleep(2 * time.Millisecond)

		err := m.Client.RetireHost(host.ID)
		if err != nil {
			return errors.Wrap(err, "Mkk.Kill fails while retiring hosts")
		}
	}

	return nil
}
