package aws

import (
	"context"
	"log"
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
	var paramToGet []string
	paramToGet = append(paramToGet, key)

	param, err := s.GetParams(ctx, s.encrypted, s.root, paramToGet)	

	if err != nil {
		log.Println(err)
	} else {
		log.Println("No error fetching: ", key, ": ", param[key])
	}

	return []byte(param[key]), err
}

// Put implements autocert Put method
func (s *SSMCache) Put(ctx context.Context, key string, data []byte) error {
	err := s.PutParam(ctx, s.encrypted, s.root, key, string(data))
	
	if err != nil {
		log.Println(err)
	}else {
		log.Println("No error putting: ", key, ": ", string(data))
	}

	return err
}

// Delete implements autocert Delete method
func (s *SSMCache) Delete(ctx context.Context, key string) error {
	err := s.DeleteParam(ctx, s.root, key)

	if err != nil {
		log.Println(err)
	}else {
		log.Println("No error deleting: ", key)
	}

	return err
}
