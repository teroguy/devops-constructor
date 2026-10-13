package docker

import (
	"fmt"
	"os/exec"
)

func IsDockerRunning() bool {
	_, err := exec.LookPath("docker")
	if err != nil {
		return false
	}

	cmd := exec.Command("docker", "info")
	if err := cmd.Run(); err != nil {
		return false
	}

	return true
}

func ComposeUp(workDir string) (string, error) {
	cmd := exec.Command("docker", "compose", "up", "-d")
	cmd.Dir = workDir

	output, err := cmd.CombinedOutput()

	if err != nil {
		fmt.Println("Не удалось запустить сервис:", err)
		return string(output), err
	}

	return string(output), nil
}
