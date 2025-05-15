package configz_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/wahyurudiyan/go-boilerplate/pkg/configz"
)

type mockConfig struct {
	AppName string `mapstructure:"APP_NAME"`
}

type mockConfigPanic struct {
	Port int `mapstructure:"PORT"`
}

func writeTempEnvFile(t *testing.T, dir, name, content string) string {
	t.Helper()
	path := filepath.Join(dir, name)
	err := os.WriteFile(path, []byte(content), 0644)
	require.NoError(t, err)
	return path
}

func TestLoad_WithDotEnvPresent(t *testing.T) {
	tmpDir := t.TempDir()

	// Setup .env file
	writeTempEnvFile(t, tmpDir, ".env", "APP_NAME=from-dot-env")

	// Change working directory to temp dir so config looks here
	oldDir, _ := os.Getwd()
	defer os.Chdir(oldDir)
	require.NoError(t, os.Chdir(tmpDir))

	var cfg mockConfig
	configz.MustLoadEnv(&cfg)

	require.Equal(t, "from-dot-env", cfg.AppName)
}

func TestLoad_FallbackToDotEnvExample(t *testing.T) {
	tmpDir := t.TempDir()

	// Setup only .env.example
	writeTempEnvFile(t, tmpDir, ".env.example", "APP_NAME=from-example")

	// Change working directory to temp dir
	oldDir, _ := os.Getwd()
	defer os.Chdir(oldDir)
	require.NoError(t, os.Chdir(tmpDir))

	var cfg mockConfig
	configz.MustLoadEnv(&cfg)

	require.Equal(t, "from-example", cfg.AppName)
}

func TestLoad_PanicsIfNoEnvFile(t *testing.T) {
	tmpDir := t.TempDir()

	// No .env or .env.example
	oldDir, _ := os.Getwd()
	defer os.Chdir(oldDir)
	require.NoError(t, os.Chdir(tmpDir))

	var cfg mockConfig

	require.Panics(t, func() {
		configz.MustLoadEnv(&cfg)
	})
}

func TestLoad_PanicsIfUnmarshalFails(t *testing.T) {
	tmpDir := t.TempDir()

	// Write invalid env file (.env)
	writeTempEnvFile(t, tmpDir, ".env", "PORT=not-integer")

	oldDir, _ := os.Getwd()
	defer os.Chdir(oldDir)
	require.NoError(t, os.Chdir(tmpDir))

	var cfg mockConfigPanic
	require.Panics(t, func() {
		configz.MustLoadEnv(&cfg)
	})
}
