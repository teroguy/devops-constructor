package generator

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/teroguy/devops-constructor/internal/core"
	"gopkg.in/yaml.v3"
)

func GenerateCompose(project core.ProjectConfig) error {
	composeService := make(map[string]ComposeService)
	composeCfgs := ComposeConfig{
		Services: make(map[string]ComposeService),
		Volumes: make(map[string]any),
	}
	for _, service := range project.Services {

		for _, vol := range service.Volumes {
			parts := strings.Split(vol, ":")
			hostPart := parts[0]

			if strings.HasPrefix(hostPart, "./") {
				dirPath := filepath.Join(project.WorkingDir, hostPart)
				_ = os.MkdirAll(dirPath, 0755)
			} else {
				composeCfgs.Volumes[hostPart] = struct{}{}
			}
		}

		composeService[service.Name] = ComposeService{
			Image:       service.Image,
			Ports:       service.Ports,
			Environment: service.Environment,
			Volumes:     service.Volumes,
		}

		composeCfgs.Services = composeService
	}

	data, err := yaml.Marshal(composeCfgs)
	if err != nil {
		return fmt.Errorf("Ошибка сбора конфигурационного файла: %w", err)
	}

	filePath := filepath.Join(project.WorkingDir, "compose.yaml")

	return os.WriteFile(filePath, data, 0644)
}
