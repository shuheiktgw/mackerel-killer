package mkk

import "github.com/mackerelio/mackerel-client-go"

type Mkk struct {
	Client *mackerel.Client
}

type Options struct {
	DryRun bool
}

func NewMkk(token string) *Mkk {
	return &Mkk{Client: mackerel.NewClient(token)}
}

func(m *Mkk) Kill(params *mackerel.FindHostsParam, filters []Filter, options *Options) ([]*mackerel.Host, error) {
	return nil, nil
}