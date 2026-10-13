package main

import (
	"errors"
	"fmt"
	"os"

	"github.com/charmbracelet/huh"
	"github.com/teroguy/devops-constructor/internal/core"
	"github.com/teroguy/devops-constructor/internal/feature/catalog"
	"github.com/teroguy/devops-constructor/internal/feature/docker"
	"github.com/teroguy/devops-constructor/internal/feature/generator"
	portcheck "github.com/teroguy/devops-constructor/internal/feature/port_check"
)

func main() {
	projectCfg := core.ProjectConfig{}

	projectDir, err := os.Getwd()
	if err != nil {
		fmt.Println("Не удалось получить путь до текущей директории")
		os.Exit(1)
	}
	projectCfg.WorkingDir = projectDir

	cat, err := catalog.New()
	if err != nil {
		fmt.Println("Ошибка при работе с каталогами") // TODO: надо будет сделать Newmust, чтобы тут не обрабатывать ошибку
		os.Exit(1)
	}

	var options []huh.Option[string]

	for _, serviceDef := range cat.GetAll() {
		label := fmt.Sprintf("%s - %s", serviceDef.Name, serviceDef.Description)
		options = append(options, huh.NewOption(label, serviceDef.ID))
	}

	var selectedIDs []string

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("Введите название проекта: ").
				Description("Это будет название вышего compose проекта").
				Value(&projectCfg.Name),

			huh.NewMultiSelect[string]().
				Title("Выберите сервисы для окружения").
				Options(options...).
				Value(&selectedIDs),
		),
	)

	if err := form.Run(); err != nil {
		if errors.Is(err, huh.ErrUserAborted) {
			fmt.Println("Выход из программы")
			return
		}

		fmt.Println("ошибка работы формы", err)
		os.Exit(1)
	}

	if len(selectedIDs) == 0 {
		fmt.Println("Не выбрано ни одного сервиса. Завершение работы.")
		return
	}

	projectCfg.Services = make([]core.Service, 0, len(selectedIDs))

	for _, selectedID := range selectedIDs {
		if serviceDef, ok := cat.GetByID(selectedID); ok {
			portResult := fmt.Sprintf("%d:%d", serviceDef.DefaultPort, serviceDef.TargetPort)
			// надо подумать не инкапсулировать ли IsPortFree, так будет только один вызов и одна сборка
			// но в варианте как сейчас может получаться так, что вызывается только один метод,
			// порт может быть свободен изначально и искать его не придется
			if portFree := portcheck.IsPortFree(serviceDef.DefaultPort); !portFree {
				newPort := portcheck.FindFreePort(serviceDef.DefaultPort)
				portResult = fmt.Sprintf("%d:%d", newPort, serviceDef.TargetPort)
			}

			projectCfg.Services = append(projectCfg.Services, core.Service{
				Name:        serviceDef.ID,
				Image:       serviceDef.DefaultImage,
				Environment: serviceDef.Env,
				Volumes:     serviceDef.Volumes,
				Ports:       []string{portResult},
			})
		}
	}

	fmt.Println("\n=== Шаги 1 и 2 успешно выполнены! ===")
	fmt.Printf("Проект: %s\n", projectCfg.Name)
	fmt.Printf("Рабочая папка: %s\n", projectCfg.WorkingDir)
	fmt.Println("Выбранные сервисы:")
	for _, s := range projectCfg.Services {
		fmt.Printf("  • %s (образ: %s, порты: %v)\n", s.Name, s.Image, s.Ports)
	}

	err = generator.GenerateCompose(projectCfg)
	if err != nil {
		fmt.Printf("Ошибка генерации compose.yaml: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("\n Файл compose.yaml и папки успешно созданы!")

	if !docker.IsDockerRunning() {
		fmt.Println("Перед запуском сервиса, нужно запустить docker")
		os.Exit(1)
	}

	output, err := docker.ComposeUp(projectCfg.WorkingDir)
	if err != nil {
		fmt.Println("Ну удалось запустить сервисы docker", output)
		os.Exit(1)
	}

	fmt.Println("\n Окружение успешно запущено в Docker!")
	fmt.Println("Ссылки для работы:")
	for _, s := range projectCfg.Services {
		for _, p := range s.Ports {
			fmt.Printf("  • Сервис %s доступен на порту: %s\n", s.Name, p)
		}
	}
}
