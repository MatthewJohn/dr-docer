package metadata

import (
	"fmt"
	"reflect"

	"gitlab.dockstudios.co.uk/dockstudios/dr-docer/pkg/domains/attribute"
)

type EntityName string
type EntityType string

type EntityId struct {
	Name EntityName
	Type EntityType
}

type Entity struct {
	Name            EntityName
	Type            EntityType
	DefaultPriority int
	Attributes      map[attribute.AttributeName]attribute.AttributeInstance
}

func NewEntity(name EntityName, entityType EntityType, defaultPriority int) (*Entity, error) {
	return &Entity{
		Name:            name,
		Type:            entityType,
		DefaultPriority: defaultPriority,
		Attributes:      map[attribute.AttributeName]attribute.AttributeInstance{},
	}, nil
}

// SetAttribute: Set attribute of entity
func (e *Entity) SetAttribute(attribute *attribute.Attribute, value any) error {
	return e.SetAttributeWithPriority(attribute, value, 0)
}

// SetAttribute: Set attribute of entity with overriden priority
func (e *Entity) SetAttributeWithPriority(attribute *attribute.Attribute, value any, overridePriority int) error {
	if attribute == nil {
		return fmt.Errorf("SetAttribute: attribute is nil")
	}
	if attributeInstance := e.GetAttributeByName(attribute.Name); attributeInstance != nil {
		return fmt.Errorf("Attribute %s already set on instance", attribute.Name)
	}

	if reflect.TypeOf(value) != attribute.Type {
		return fmt.Errorf("Value for attribute %s is not type %s", attribute.Name, attribute.Type)
	}

	attributeInstance := attribute.CeateInstance()
	attributeInstance.Value = value
	if overridePriority == 0 {
		overridePriority = e.DefaultPriority
	}

	// Assign attribute to entity
	e.registerAttributeInstance(attributeInstance)

	return nil
}

// registerAttributeInstance: Register an attribute instance with entity
func (e *Entity) registerAttributeInstance(attributeInstance attribute.AttributeInstance) {
	e.Attributes[attributeInstance.Attribute.Name] = attributeInstance
}

func (e *Entity) GetName() EntityName {
	return e.Name
}

func (e *Entity) GetType() EntityType {
	return e.Type
}

func (e *Entity) GetAttributes() map[attribute.AttributeName]attribute.AttributeInstance {
	return e.Attributes
}

func (e *Entity) GetAttributeByName(attributeName attribute.AttributeName) *attribute.AttributeInstance {
	if attribute, ok := e.Attributes[attributeName]; ok {
		return &attribute
	}
	return nil
}

func (e *Entity) MergeAttributes(new *Entity) {
	for _, newAttribute := range new.GetAttributes() {
		if existingAttribute := e.GetAttributeByName(newAttribute.Attribute.Name); existingAttribute != nil {
			existingAttribute.MergeAttribute(&newAttribute)
		}
	}
}
