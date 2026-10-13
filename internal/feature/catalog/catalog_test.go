package catalog_test

import (
	"testing"

	"github.com/teroguy/devops-constructor/internal/feature/catalog"
)

func TestCatalog_New_LoadsValidServices(t *testing.T) {
	cat, err := catalog.New()
	if err != nil {
		t.Fatalf("ожидалась успешная загрузка каталога, получена ошибка: %v", err)
	}

	services := cat.GetAll()
	if len(services) == 0 {
		t.Fatal("каталог пуст, ожидалось наличие шаблонов сервисов")
	}

	validCategories := map[string]bool{
		"backend":  true,
		"database": true,
		"tools":    true,
	}

	for _, s := range services {
		if s.ID == "" {
			t.Errorf("найден сервис с пустым ID: %+v", s)
		}
		if s.Name == "" {
			t.Errorf("сервис %s имеет пустое имя", s.ID)
		}
		if s.DefaultImage == "" {
			t.Errorf("сервис %s не имеет default_image", s.ID)
		}
		if s.DefaultPort <= 0 {
			t.Errorf("сервис %s имеет некорректный default_port: %d", s.ID, s.DefaultPort)
		}
		if !validCategories[s.Category] {
			t.Errorf("сервис %s имеет неизвестную категорию: %q", s.ID, s.Category)
		}
	}
}

func TestCatalog_GetByID(t *testing.T) {
	cat, err := catalog.New()
	if err != nil {
		t.Fatalf("ошибка инициализации каталога: %v", err)
	}

	php, found := cat.GetByID("php")
	if !found {
		t.Fatal("сервис 'php' должен существовать в каталоге")
	}
	if php.Name != "PHP 8.2 (Apache)" {
		t.Errorf("ожидалось имя 'PHP 8.2 (Apache)', получено %q", php.Name)
	}

	_, foundFake := cat.GetByID("unknown_service")
	if foundFake {
		t.Error("несуществующий сервис не должен быть найден")
	}
}

func TestCatalog_GetByCategory(t *testing.T) {
	cat, err := catalog.New()
	if err != nil {
		t.Fatalf("ошибка инициализации каталога: %v", err)
	}

	byCategory := cat.GetByCategory()

	if len(byCategory["backend"]) == 0 {
		t.Error("в категории 'backend' должен быть как минимум 1 сервис")
	}
	if len(byCategory["database"]) == 0 {
		t.Error("в категории 'database' должен быть как минимум 1 сервис")
	}
	if len(byCategory["tools"]) == 0 {
		t.Error("в категории 'tools' должен быть как минимум 1 сервис")
	}
}
