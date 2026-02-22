package metadata

type EntityType string

const (
	EntityTypeSerer   EntityType = "server"
	EntityTypeService EntityType = "service"
)

type Entity struct {
	Name string
	Type EntityType
}
