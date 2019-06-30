package providers

// TODO: Patreon, Social (generic social media accounts)

import (
	"github.com/pkg/errors"
)

type Providers []*Provider

// Provider specifies the minimum methods a provider needs
// to implement.
type Provider interface {
	GetName() string
}

func (pr *Providers) Add(provider *Provider) {
	*pr = append(*pr, provider)
}

func (pr *Providers) GetProvider(key string) (*Provider, error) {
	for _, provider := range *pr {
		if (*provider).GetName() == key {
			return provider, nil
		}
	}

	return nil, errors.Errorf("provider '%v' could not be found", key)
}
