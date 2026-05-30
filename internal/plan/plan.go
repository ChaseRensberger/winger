package plan

import (
	"fmt"
	"os"

	"winger/internal/config"
)

const DefaultTemplate = `[%s]
Plan:

`

func Read() (string, error) {
	p, err := config.PlanPath()
	if err != nil {
		return "", err
	}

	data, err := os.ReadFile(p)
	if err != nil {
		return "", fmt.Errorf("could not read .plan file: %w", err)
	}

	return string(data), nil
}

func CreateDefault(handle string) error {
	p, err := config.PlanPath()
	if err != nil {
		return err
	}

	if _, err := os.Stat(p); err == nil {
		return nil
	}

	content := fmt.Sprintf(DefaultTemplate, handle)
	if err := os.WriteFile(p, []byte(content), 0644); err != nil {
		return fmt.Errorf("could not create .plan file: %w", err)
	}

	return nil
}
