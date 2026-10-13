package core

// конфигурация одного контейнера
type Service struct {
	Name        string            // Имя сервиса
	Image       string            // образ с тегом
	Ports       []string          // хз как правильно назвать, но пример ["8080:80"]
	Environment map[string]string // Переменные окружения
	Volumes     []string          // Монтируемые директории
}

// описывает запускаемый стек проекта
type ProjectConfig struct {
	Name       string    // Название проекта/лабы
	WorkingDir string    // Путь к папке проекта
	Services   []Service // Список сервисов
}
