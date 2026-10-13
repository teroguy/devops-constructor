package catalog

import (
	"fmt"
	"io/fs"
	"path"
	"strings"

	"github.com/teroguy/devops-constructor/templates"
	"gopkg.in/yaml.v3"
)

const servicesDir = "services"

// один сервис из шаблонов
type ServiceDefinition struct {
	ID           string            `yaml:"id"`
	Name         string            `yaml:"name"`
	Category     string            `yaml:"category"`
	Description  string            `yaml:"description"`
	DefaultImage string            `yaml:"default_image"`
	DefaultPort  int               `yaml:"default_port"`
	TargetPort   int               `yaml:"target_port"`
	Env          map[string]string `yaml:"env"`
	Volumes      []string          `yaml:"volumes"`
}

// загруженные шаблоны сервисов
type Catalog struct {
	services map[string]ServiceDefinition
}

func New() (*Catalog, error) {
	
	entries, err := fs.ReadDir(templates.ServicesFS, servicesDir)
	if err != nil {
		return nil, fmt.Errorf("не удалось прочитать директорию шаблонов: %w", err)
	}
	
	services := make(map[string]ServiceDefinition)

	for _, entry := range entries {

		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".yaml") {
			continue
		}

		filePath := path.Join(servicesDir, entry.Name())
		data, err := templates.ServicesFS.ReadFile(filePath)
		if err != nil {
			return nil, fmt.Errorf("ошибка чтения файла %s: %w", filePath, err)
		}

		var def ServiceDefinition
		if err := yaml.Unmarshal(data, &def); err != nil {
			return nil, fmt.Errorf("ошибка парсинга YAML в %s: %w", filePath, err)
		}

		services[def.ID] = def
	}

	return &Catalog{
		services: services,
	}, nil
}

func (c *Catalog) GetAll() []ServiceDefinition {
	result := make([]ServiceDefinition, 0, len(c.services))
	for _, def := range c.services {
		result = append(result, def)
	}
	return result
}

func (c *Catalog) GetByCategory() map[string][]ServiceDefinition {
	result := make(map[string][]ServiceDefinition)
	for _, def := range c.services {
		result[def.Category] = append(result[def.Category], def)
	}
	return result
}

func (c *Catalog) GetByID(id string) (ServiceDefinition, bool) {
	def, ok := c.services[id]
	return def, ok
}
