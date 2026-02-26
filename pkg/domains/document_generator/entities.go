package documentgenerator

import (
	"fmt"

	"gitlab.dockstudios.co.uk/dockstudios/dr-docer/pkg/domains/attribute"
	"gitlab.dockstudios.co.uk/dockstudios/dr-docer/pkg/domains/metadata"
)

type TemplateMetadata struct {
	EntityType string `yaml:"entity_type"`
}

type TemplateEntityShim struct {
	Name       string
	attributes map[attribute.AttributeName]attribute.AttributeInstance
}

func NewTemplateEntityShim(entity *metadata.Entity) (*TemplateEntityShim, error) {
	if entity == nil {
		return nil, fmt.Errorf("NewTemplateEntityShim: entity is nil")
	}
	return &TemplateEntityShim{
		Name:       string(entity.GetName()),
		attributes: entity.Attributes,
	}, nil
}

func (t *TemplateEntityShim) Get(attributeName string) any {
	if attr, ok := t.attributes[attribute.AttributeName(attributeName)]; ok {
		return attr.Value
	}
	return ""
}
