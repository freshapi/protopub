package protopub

import "github.com/deislabs/oras/pkg/content"

// DescriptorStore is an implementation of oci provider
type DescriptorStore struct {
	*content.Memorystore
}

// NewDescriptorStore creates DescriptorStore
func NewDescriptorStore() *DescriptorStore {
	return &DescriptorStore{
		Memorystore: content.NewMemoryStore(),
	}
}
