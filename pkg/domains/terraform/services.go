package terraform

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclparse"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/zclconf/go-cty/cty"
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing/object"
)

type TerraformParser struct {
}

func NewTerraformParser() (*TerraformParser, error) {
	return &TerraformParser{}, nil
}

func (t *TerraformParser) ParseTerraform(repo *git.Repository, dir string) (*TerraformModel, error) {
	files, err := loadTerraformFiles(repo, dir)
	if err != nil {
		return nil, err
	}
	model, err := parseTerraformHCLFiles(files)
	if err != nil {
		return nil, err
	}
	return model, err
}

func loadTerraformFiles(repo *git.Repository, dir string) (map[string][]byte, error) {
	files := map[string][]byte{}

	ref, err := repo.Head()
	if err != nil {
		return nil, err
	}

	commit, err := repo.CommitObject(ref.Hash())
	if err != nil {
		return nil, err
	}

	tree, err := commit.Tree()
	if err != nil {
		return nil, err
	}

	err = tree.Files().ForEach(func(f *object.File) error {
		if strings.HasPrefix(f.Name, dir) && strings.HasSuffix(f.Name, ".tf") {

			content, err := f.Contents()
			if err != nil {
				return err
			}

			files[f.Name] = []byte(content)
		}
		return nil
	})

	return files, err
}

func parseTerraformHCLFiles(files map[string][]byte) (*TerraformModel, error) {
	tmpDir, err := os.MkdirTemp("", "tf-hcl")
	if err != nil {
		return nil, err
	}
	defer os.RemoveAll(tmpDir)

	for name, content := range files {
		fullPath := filepath.Join(tmpDir, name)
		if err := os.MkdirAll(filepath.Dir(fullPath), 0o755); err != nil {
			return nil, err
		}
		if err := os.WriteFile(fullPath, content, 0o644); err != nil {
			return nil, err
		}
	}

	parser := hclparse.NewParser()
	model := &TerraformModel{
		Variables: map[string]cty.Value{},
		Locals:    map[string]cty.Value{},
	}

	// parse all files
	for name, content := range files {
		file, diags := parser.ParseHCL(content, name)
		if diags.HasErrors() {
			return nil, fmt.Errorf("parsing %s failed: %s", name, diags.Error())
		}
		// extract top-level blocks (resource, module, variable, locals)
		processFile(file, name, model)
	}

	return model, nil
}
func processFile(file *hcl.File, filename string, model *TerraformModel) {
	body := file.Body.(*hclsyntax.Body)
	for _, block := range body.Blocks {
		switch block.Type {
		case "variable":
			// capture default if present
			if defAttr, ok := block.Body.Attributes["default"]; ok {
				val, _ := defAttr.Expr.Value(nil) // literal evaluation
				model.Variables[block.Labels[0]] = val
			}
		case "locals":
			for name, attr := range block.Body.Attributes {
				val, _ := attr.Expr.Value(nil)
				model.Locals[name] = val
			}
		case "resource":
			processResourceBlock(block, filename, model)
		case "module":
			processModuleBlock(block, filename, model)
		}
	}
}

func processResourceBlock(block *hclsyntax.Block, filename string, model *TerraformModel) {
	// for simplicity, only handles literal values and for_each map literals
	forEachMap := map[string]cty.Value{}
	if feAttr, ok := block.Body.Attributes["for_each"]; ok {
		val, _ := feAttr.Expr.Value(nil)
		if val.Type().IsObjectType() || val.Type().IsMapType() {
			it := val.ElementIterator()
			for it.Next() {
				k, v := it.Element()
				forEachMap[k.AsString()] = v
			}
		}
	}

	if len(forEachMap) == 0 {
		// single resource
		attrMap := map[string]cty.Value{}
		for name, attr := range block.Body.Attributes {
			val, _ := attr.Expr.Value(nil)
			attrMap[name] = val
		}
		model.Resources = append(model.Resources, TFResource{
			Type:       block.Labels[0],
			Name:       block.Labels[1],
			File:       filename,
			Attributes: attrMap,
		})
	} else {
		// expand for_each
		for key, val := range forEachMap {
			attrMap := map[string]cty.Value{}
			for name, attr := range block.Body.Attributes {
				v, _ := attr.Expr.Value(&hcl.EvalContext{
					Variables: map[string]cty.Value{
						"each": cty.ObjectVal(map[string]cty.Value{
							"key":   cty.StringVal(key),
							"value": val,
						}),
					},
				})
				attrMap[name] = v
			}
			model.Resources = append(model.Resources, TFResource{
				Type:       block.Labels[0],
				Name:       fmt.Sprintf("%s[%s]", block.Labels[1], key),
				File:       filename,
				Attributes: attrMap,
			})
		}
	}
}

func processModuleBlock(block *hclsyntax.Block, filename string, model *TerraformModel) {
	// for_each similar to resources
	forEachMap := map[string]cty.Value{}
	if feAttr, ok := block.Body.Attributes["for_each"]; ok {
		val, _ := feAttr.Expr.Value(nil)
		if val.Type().IsObjectType() || val.Type().IsMapType() {
			it := val.ElementIterator()
			for it.Next() {
				k, v := it.Element()
				forEachMap[k.AsString()] = v
			}
		}
	}

	if len(forEachMap) == 0 {
		inputs := map[string]cty.Value{}
		for name, attr := range block.Body.Attributes {
			val, _ := attr.Expr.Value(nil)
			inputs[name] = val
		}
		model.Modules = append(model.Modules, TFModule{
			Name:   block.Labels[0],
			File:   filename,
			Source: inputs["source"].AsString(),
			Inputs: inputs,
		})
	} else {
		for key, val := range forEachMap {
			inputs := map[string]cty.Value{}
			for name, attr := range block.Body.Attributes {
				v, _ := attr.Expr.Value(&hcl.EvalContext{
					Variables: map[string]cty.Value{
						"each": cty.ObjectVal(map[string]cty.Value{
							"key":   cty.StringVal(key),
							"value": val,
						}),
					},
				})
				inputs[name] = v
			}
			model.Modules = append(model.Modules, TFModule{
				Name:   fmt.Sprintf("%s[%s]", block.Labels[0], key),
				File:   filename,
				Source: inputs["source"].AsString(),
				Inputs: inputs,
			})
		}
	}
}
