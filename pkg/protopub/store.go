package protopub

import "github.com/deislabs/oras/pkg/content"

type DescriptorStore struct {
	*content.Memorystore
}

func NewDescriptorStore() *DescriptorStore {
	return &DescriptorStore{
		Memorystore: content.NewMemoryStore(),
	}
}
