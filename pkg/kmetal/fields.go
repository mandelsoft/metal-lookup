package kmetal

import (
	"fmt"

	"github.com/gardener/controller-manager-library/pkg/fieldpath"
	"github.com/metal-stack/metal-go/api/models"
)

type Field struct {
	source fieldpath.Field
	target fieldpath.Node
}

type Fields map[string]Field

var base = &models.V1MachineResponse{}

func (this Fields) Add(target, path string) error {
	f, err := fieldpath.NewField(base, path)
	if err != nil {
		return fmt.Errorf("invalid field %q: %s", path, err)
	}
	n, err := fieldpath.Compile(target)
	if err != nil {
		return fmt.Errorf("invalid name %q; %s", target, err)
	}
	this[target] = Field{
		source: f,
		target: n,
	}

	return nil
}
