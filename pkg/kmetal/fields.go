package kmetal

import (
	"github.com/gardener/controller-manager-library/pkg/fieldpath"
	"github.com/metal-stack/metal-go/api/models"
)

type Fields map[string]fieldpath.Field

var base = &models.V1MachineResponse{}

func (this Fields) Add(name, path string) error {
	f, err := fieldpath.NewField(base, path)
	if err != nil {
		return err
	}
	this[name] = f

	return nil
}
