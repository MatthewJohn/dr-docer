package attribute

import "reflect"

// defaultPriority Default priority for new Attribute Instances.
// Priority goes 0 (least) -> inf (highest).
// Therefore, default 0 of unset is best
// There MAY be a case for an attribute to be overriden by
// a attribute instance priority of negative (e.g. -5).
// This _could_ be useful for overriding attributes where
// a default is unknown or optionally set
const defaultPriority int = 0

// AttributeName Name of attribute
type AttributeName string

// Attribute Structure for holding information about a type
// of attribute that will exist on an entity.
type Attribute struct {
	Name         AttributeName
	Type         reflect.Type
	DefaultValue any
}

// CreateInstance Returns a new instance of the Attribute
// to be assigned to an Entity
func (a *Attribute) CeateInstance() AttributeInstance {
	return AttributeInstance{
		Attribute: a,
		Priority:  defaultPriority,
		Value:     a.DefaultValue,
	}
}

// AttributeInstance An instance of an attribute, to be assigned to an Entity
type AttributeInstance struct {
	Attribute *Attribute
	Priority  int
	Value     any
}

func (ai *AttributeInstance) MergeAttribute(new *AttributeInstance) {
	if new == nil {
		return
	}

	if ai.Attribute != new.Attribute {
		return
	}

	if ai.Priority >= new.Priority {
		return
	}

	// Merge attribute
	ai.Value = new.Value
	ai.Priority = new.Priority
}
