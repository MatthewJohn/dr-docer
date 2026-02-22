package metadata

type EntityType string

const (
	EntityTypeSerer   EntityType = "server"
	EntityTypeService EntityType = "service"
)

type StorageMetata struct {
}

type Metadata struct {
	Type         EntityType    `yaml:"type"`
	Name         string        `yaml:"name"`
	Criticality  string        `yaml:"criticality"`
	Host         string        `yaml:"host"`
	Storage      StorageMetata `yaml:"storage"`
	Dependencies []string      `yaml:"dependencies"`
	Terraform    []string      `yaml:"terraform"`
}
