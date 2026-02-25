package attribute

import "fmt"

type AttributeFactory struct {
	Attributes map[AttributeName]Attribute
}

func NewAttributeFactory() (*AttributeFactory, error) {
	return &AttributeFactory{
		Attributes: map[AttributeName]Attribute{},
	}, nil
}

func (af *AttributeFactory) RegisterAttribute(attribute *Attribute) error {
	if attribute == nil {
		return fmt.Errorf("RegisterAttribute: Attribute is an nil pointer")
	}
	attributeName := attribute.Name
	if attributeName == "" {
		return fmt.Errorf("RegisterAttribute: Attribute name empty")
	}
	if _, ok := af.Attributes[attributeName]; ok {
		return fmt.Errorf("RegisterAttribute: Attribute %s already registered", attributeName)
	}

	af.Attributes[attributeName] = *attribute
	return nil
}

func (af *AttributeFactory) GetAttributeByName(name AttributeName) *Attribute {
	if attribute, ok := af.Attributes[name]; ok {
		return &attribute
	}
	return nil
}
