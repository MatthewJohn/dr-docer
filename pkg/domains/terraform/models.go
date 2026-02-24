package terraform

import "github.com/zclconf/go-cty/cty"

// Terraform model structures
type TerraformModel struct {
	Resources []TFResource
	Modules   []TFModule
	Variables map[string]cty.Value
	Locals    map[string]cty.Value
}

type TFResource struct {
	Type       string
	Name       string
	File       string
	Attributes map[string]cty.Value
}

type TFModule struct {
	Name   string
	File   string
	Source string
	Inputs map[string]cty.Value
}
