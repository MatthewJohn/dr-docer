package documentgenerator

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"text/template"

	"gitlab.dockstudios.co.uk/dockstudios/dr-docer/pkg/domains/metadata"
	"go.yaml.in/yaml/v3"
)

const TemplateExtension string = ".md"

type DocumentGenerator struct {
	documentStorage   DocumentStorage
	templateDirectory string
	templates         map[string][]byte
}

func extractMetadataFromTemplate(templateData []byte) (*TemplateMetadata, error) {
	// Convert to entity
	decoder := yaml.NewDecoder(bytes.NewReader(templateData))

	for {
		var templateMetadata TemplateMetadata
		err := decoder.Decode(&templateMetadata)
		if err != nil {
			if err.Error() == "EOF" {
				break
			}
			fmt.Println(err)
			continue
		}
		return &templateMetadata, nil
	}
	return nil, fmt.Errorf("Could not find metadata in template")
}

func getTemplates(directory string) (map[string][]byte, error) {
	fileInfo, err := os.Stat(directory)
	if err != nil {
		return nil, fmt.Errorf("Error checking template directory: %s", err.Error())
	}

	if !fileInfo.IsDir() {
		return nil, fmt.Errorf("Template directory is not a directory")
	}
	matches, err := filepath.Glob(filepath.Join(directory, fmt.Sprintf("*%s", TemplateExtension)))
	if err != nil {
		return nil, fmt.Errorf("Error globbing templates: %s", err.Error())
	}
	templates := map[string][]byte{}
	for _, match := range matches {
		data, err := os.ReadFile(match)
		if err != nil {
			return nil, fmt.Errorf("Error reading template: %s: %s", match, err.Error())
		}
		metadata, err := extractMetadataFromTemplate(data)
		if err != nil {
			return nil, fmt.Errorf("Error extracting metadata from template: %s: %s", match, err.Error())
		}

		if _, ok := templates[metadata.EntityType]; ok {
			return nil, fmt.Errorf("Found duplicate template for entity Type: %s", metadata.EntityType)
		}
		templates[metadata.EntityType] = data
	}
	return templates, nil
}

func NewDocumentGenerator(documentStorage DocumentStorage, templateDirectory string) (*DocumentGenerator, error) {
	if documentStorage == nil {
		return nil, fmt.Errorf("NewDocumentGenerator: documentStoage is nil")
	}
	templates, err := getTemplates(templateDirectory)
	if err != nil {
		return nil, err
	}
	return &DocumentGenerator{
		documentStorage:   documentStorage,
		templateDirectory: templateDirectory,
		templates:         templates,
	}, nil
}

func (dg *DocumentGenerator) getTemplateForEntityType(entityType metadata.EntityType) ([]byte, error) {
	if template, ok := dg.templates[string(entityType)]; ok {
		return template, nil
	}
	return []byte{}, fmt.Errorf("Template not found for entity type: %s", entityType)
}

func (dg *DocumentGenerator) GenerateDocumentForEntity(entity metadata.Entity) error {
	templateRaw, err := dg.getTemplateForEntityType(entity.GetType())
	if err != nil {
		return err
	}
	templateRenderer := template.New(string(entity.GetName()))
	parsedTemplate, err := templateRenderer.Parse(string(templateRaw))
	if err != nil {
		return err
	}

	entityShim, err := NewTemplateEntityShim(&entity)
	if err != nil {
		return err
	}

	var b bytes.Buffer
	err = parsedTemplate.Execute(io.Writer(&b), entityShim)
	if err != nil {
		return err
	}
	return dg.documentStorage.StoreDocument(entity.GetName(), entity.GetType(), b.Bytes())
}
