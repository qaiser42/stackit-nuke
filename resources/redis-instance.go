package resources

import (
	"context"
	"fmt"

	"github.com/ekristen/libnuke/pkg/registry"
	"github.com/ekristen/libnuke/pkg/resource"
	"github.com/ekristen/libnuke/pkg/types"

	"github.com/qaiser42/stackit-nuke/pkg/stackit"
)

const RedisInstanceResource = "RedisInstance"

func init() {
	registry.Register(&registry.Registration{
		Name:     RedisInstanceResource,
		Scope:    stackit.ProjectScope,
		Resource: &RedisInstance{},
		Lister:   &RedisInstanceLister{},
	})
}

type RedisInstance struct {
	*BaseResource `property:",inline"`
	ID            string
	Name          string
}

func (r *RedisInstance) Remove(_ context.Context) error {
	return fmt.Errorf("RedisInstance.Remove not yet implemented")
}
func (r *RedisInstance) Properties() types.Properties { return PropsFromStruct(r) }
func (r *RedisInstance) String() string               { return r.Name }

type RedisInstanceLister struct{}

func (l *RedisInstanceLister) List(_ context.Context, _ any) ([]resource.Resource, error) {
	return []resource.Resource{}, nil
}
