package mkk

type Filter func(*Mkk, []*mackerel.Host) ([]*mackerel.Host, error)
