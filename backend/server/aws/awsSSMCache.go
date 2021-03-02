package aws

import (
	"context"
)

// SSMCache uses aws paramater store to cache tls certs
type SSMCache struct {
	*SSM
	encrypted bool
	root string
}

// NewSSMCache returns a new ssm cache
func NewSSMCache(encrypted bool, root string) *SSMCache {
	return &SSMCache{
		encrypted: encrypted,
		root: root,
	}
}

// Get implements autocert Get method
func (s *SSMCache) Get(ctx context.Context, key string) ([]byte, error) {
	param, err := s.GetParams(ctx, s.encrypted, s.root, []string{key})	

	return []byte(param[key]), err
}

// Put implements autocert Put method
func (s *SSMCache) Put(ctx context.Context, key string, data []byte) error {
	return s.PutParam(ctx, s.encrypted, s.root, key, string(data))
}

// Delete implements autocert Delete method
func (s *SSMCache) Delete(ctx context.Context, key string) error {
	return s.DeleteParam(ctx, s.root, key)
}
