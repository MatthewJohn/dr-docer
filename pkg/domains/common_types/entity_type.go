package commontypes

import (
	"reflect"

	"gitlab.dockstudios.co.uk/dockstudios/dr-docer/pkg/domains/attribute"
	"gitlab.dockstudios.co.uk/dockstudios/dr-docer/pkg/domains/metadata"
)

var (
	EntityServer  metadata.EntityType = "server"
	EntityService metadata.EntityType = "service"
)

var AttributeUrl = attribute.Attribute{
	Name:         "url",
	Type:         reflect.TypeOf(""),
	DefaultValue: "",
}

var AttributeIpAddress attribute.Attribute = attribute.Attribute{
	Name:         "ip_address",
	Type:         reflect.TypeOf(""),
	DefaultValue: "",
}
