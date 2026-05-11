package resources

import (
	"context"
	"fmt"

	"github.com/ekristen/libnuke/pkg/registry"
	"github.com/ekristen/libnuke/pkg/resource"
	"github.com/ekristen/libnuke/pkg/types"

	"github.com/qaiser42/stackit-nuke/pkg/stackit"
)

const MongoDBFlexInstanceResource = "MongoDBFlexInstance"

func init() {
	registry.Register(&registry.Registration{
		Name:     MongoDBFlexInstanceResource,
		Scope:    stackit.ProjectScope,
		Resource: &MongoDBFlexInstance{},
		Lister:   &MongoDBFlexInstanceLister{},
	})
}

type MongoDBFlexInstance struct {
	*BaseResource `property:",inline"`
	ID            string
	Name          string
}

func (r *MongoDBFlexInstance) Remove(_ context.Context) error {
	return fmt.Errorf("MongoDBFlexInstance.Remove not yet implemented")
}
func (r *MongoDBFlexInstance) Properties() types.Properties { return PropsFromStruct(r) }
func (r *MongoDBFlexInstance) String() string               { return r.Name }

type MongoDBFlexInstanceLister struct{}

func (l *MongoDBFlexInstanceLister) List(_ context.Context, _ any) ([]resource.Resource, error) {
	return []resource.Resource{}, nil
}
