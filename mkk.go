package mkk

import (
	"github.com/mackerelio/mackerel-client-go"
	"github.com/pkg/errors"
	"log"
	"time"
)

type Mkk struct {
	Client *mackerel.Client
}

type Options struct {
	DryRun bool
}

func NewMkk(token string) *Mkk {
	return &Mkk{Client: mackerel.NewClient(token)}
}

func (m *Mkk) Kill(params *mackerel.FindHostsParam, filters []Filter, options *Options) ([]*mackerel.Host, error) {
	hosts, err := m.Client.FindHosts(params)
	if err != nil {
		return nil, errors.Wrap(err, "Mkk.Kill fails while finding hosts")
	}

	for _, f := range filters {
		hosts, err = f.Apply(m, hosts)
		if err != nil {
			return nil, errors.Wrap(err, "Mkk.Kill fails while applying filters")
		}
	}

	if options.DryRun {
		log.Println("mackerel-killer is under Dry Run mode")

		for _, host := range hosts {
			log.Printf("ID: %v, Name: %v, Roles: %v will be retired\n", host.ID, host.Name, host.Roles)
		}

		return hosts, nil
	}

	for _, host := range hosts {
		time.Sleep(2 * time.Millisecond)

		err = m.Client.RetireHost(host.ID)
		if err != nil {
			return nil, errors.Wrap(err, "Mkk.Kill fails while retiring hosts")
		}

		log.Printf("Retired ID: %v, Name: %v, Roles: %v\n", host.ID, host.Name, host.Roles)
	}

	return hosts, nil
}
