//go:build mage

package main

import (
	"os"
	"os/exec"
)

func dc(args ...string) error {
	cmd := exec.Command("docker", append([]string{"compose"}, args...)...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// Start builds and starts all services (DB + bot) in background.
func Start() error {
	return dc("up", "-d", "--build")
}

// Stop stops and removes all containers (data is preserved in volumes).
func Stop() error {
	return dc("down")
}

// Restart rebuilds and restarts only the bot container.
func Restart() error {
	return dc("up", "-d", "--build", "bot")
}

// Logs streams live logs from the bot container.
func Logs() error {
	return dc("logs", "-f", "bot")
}

// Status shows current status of all containers.
func Status() error {
	return dc("ps")
}

// DB starts only the local database (only available in dev via override).
func DB() error {
	return dc("-f", "docker-compose.yml", "-f", "docker-compose.override.yml", "up", "-d", "db")
}
