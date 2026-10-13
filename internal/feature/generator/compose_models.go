package generator

// корень файла yaml
type ComposeConfig struct {
	Services map[string]ComposeService `yaml:"services"`
	Volumes  map[string]any            `yaml:"volumes,omitempty"` // omitempty скроет секцию, если томов нет
}

// свойства контейнера
type ComposeService struct {
	Image       string            `yaml:"image"`
	Ports       []string          `yaml:"ports,omitempty"`
	Environment map[string]string `yaml:"environment,omitempty"`
	Volumes     []string          `yaml:"volumes,omitempty"`
}
