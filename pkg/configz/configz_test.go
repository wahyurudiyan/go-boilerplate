package configz

import (
	"os"
	"path/filepath"
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/stretchr/testify/mock"
)

// Mock untuk AWS SSM Parameter Store
type MockParameterStore struct {
	mock.Mock
}

func (m *MockParameterStore) GetAllParametersByPath(path string, recursive bool) (map[string]string, error) {
	args := m.Called(path, recursive)
	return args.Get(0).(map[string]string), args.Error(1)
}

// Test struct untuk konfigurasi
type TestConfig struct {
	DatabaseURL string `mapstructure:"DATABASE_URL"`
	Port        int    `mapstructure:"PORT"`
	Debug       bool   `mapstructure:"DEBUG"`
	AppName     string `mapstructure:"APP_NAME"`
}

func TestConfigz(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Configz Suite")
}

var _ = Describe("Configz", func() {
	var tempDir string

	BeforeEach(func() {
		var err error
		tempDir, err = os.MkdirTemp("", "configz_test")
		Expect(err).ToNot(HaveOccurred())
	})

	AfterEach(func() {
		os.RemoveAll(tempDir)
	})

	Describe("LoadFromDotenv", func() {
		Context("when dotenv file exists and is valid", func() {
			It("should load configuration successfully", func() {
				// Arrange
				dotenvContent := `DATABASE_URL=postgres://localhost/test
PORT=8080
DEBUG=true
APP_NAME=test-app`

				envFile := filepath.Join(tempDir, ".env")
				err := os.WriteFile(envFile, []byte(dotenvContent), 0644)
				Expect(err).ToNot(HaveOccurred())

				var config TestConfig

				// Act
				err = LoadFromDotenv(envFile, &config)

				// Assert
				Expect(err).ToNot(HaveOccurred())
				Expect(config.DatabaseURL).To(Equal("postgres://localhost/test"))
				Expect(config.Port).To(Equal(8080))
				Expect(config.Debug).To(BeTrue())
				Expect(config.AppName).To(Equal("test-app"))
			})
		})

		Context("when dotenv file does not exist", func() {
			It("should return an error", func() {
				// Arrange
				nonExistentFile := filepath.Join(tempDir, "non-existent.env")
				var config TestConfig

				// Act
				err := LoadFromDotenv(nonExistentFile, &config)

				// Assert
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("no such file or directory"))
			})
		})

		Context("when dotenv file has invalid format", func() {
			It("should return error", func() {
				// Arrange
				dotenvContent := `DATABASE_URL=postgres://localhost/test
INVALID_LINE_WITHOUT_EQUALS
PORT=8080
ANOTHER_INVALID=
DEBUG=true`

				envFile := filepath.Join(tempDir, ".env")
				err := os.WriteFile(envFile, []byte(dotenvContent), 0644)
				Expect(err).ToNot(HaveOccurred())

				var config TestConfig

				// Act
				err = LoadFromDotenv(envFile, &config)

				// Assert
				Expect(err).To(HaveOccurred())
			})
		})

		Context("when output struct is nil", func() {
			It("should return an error", func() {
				// Arrange
				dotenvContent := `DATABASE_URL=postgres://localhost/test`
				envFile := filepath.Join(tempDir, ".env")
				err := os.WriteFile(envFile, []byte(dotenvContent), 0644)
				Expect(err).ToNot(HaveOccurred())

				// Act
				err = LoadFromDotenv(envFile, nil)

				// Assert
				Expect(err).To(HaveOccurred())
			})
		})

		Context("when environment variables override dotenv values", func() {
			It("should prioritize environment variables", func() {
				// Arrange
				dotenvContent := `PORT=8080
DEBUG=false`

				envFile := filepath.Join(tempDir, ".env")
				err := os.WriteFile(envFile, []byte(dotenvContent), 0644)
				Expect(err).ToNot(HaveOccurred())

				// Set environment variable that should override
				os.Setenv("PORT", "9090")
				defer os.Unsetenv("PORT")

				var config TestConfig

				// Act
				err = LoadFromDotenv(envFile, &config)

				// Assert
				Expect(err).ToNot(HaveOccurred())
				Expect(config.Port).To(Equal(9090)) // Should be from env var, not dotenv
				Expect(config.Debug).To(BeFalse())  // Should be from dotenv
			})
		})
	})

	Describe("LoadFromAWSParameterStore", func() {
		Context("when AWS Parameter Store returns valid parameters", func() {
			It("should load configuration successfully", func() {
				// Arrange
				// setup := ParameterStoreSetup{
				// 	Path:   "/go-boilerplate/",
				// 	Type:   "json",
				// 	Prefix: "",
				// }

				// // Mock parameters dari AWS Parameter Store
				// mockParams := map[string]string{
				// 	"DATABASE_URL": "postgres://aws-rds/prod",
				// 	"PORT":         "3000",
				// 	"DEBUG":        "false",
				// 	"APP_NAME":     "production-app",
				// }

				// // Note: Dalam implementasi nyata, Anda perlu mock awsssm.NewParameterStore()
				// // Untuk test ini, kita asumsikan ada cara untuk inject mock

				// var config TestConfig

				// Act & Assert
				// Karena kode asli menggunakan awsssm.NewParameterStore() secara langsung,
				// kita perlu refactor kode untuk dependency injection agar bisa di-test
				// Untuk sekarang, kita skip test ini atau mock secara berbeda

				Skip("Skipping AWS Parameter Store test - requires refactoring for dependency injection")
			})
		})

		Context("when AWS Parameter Store setup has prefix", func() {
			It("should use the prefix correctly", func() {
				// Arrange
				// setup := ParameterStoreSetup{
				// 	Path:   "/go-boilerplate/",
				// 	Type:   "json",
				// 	Prefix: "APP",
				// }

				Skip("Skipping AWS Parameter Store test - requires refactoring for dependency injection")
			})
		})

		Context("when AWS Parameter Store connection fails", func() {
			It("should return an error", func() {
				// Arrange
				// setup := ParameterStoreSetup{
				// 	Path:   "/go-boilerplate/",
				// 	Type:   "json",
				// 	Prefix: "",
				// }

				// var config TestConfig

				Skip("Skipping AWS Parameter Store test - requires refactoring for dependency injection")
			})
		})

		Context("when GetAllParametersByPath fails", func() {
			It("should return an error", func() {
				// Arrange
				// setup := ParameterStoreSetup{
				// 	Path:   "/go-boilerplate/",
				// 	Type:   "json",
				// 	Prefix: "",
				// }

				// var config TestConfig

				Skip("Skipping AWS Parameter Store test - requires refactoring for dependency injection")
			})
		})

		Context("when viper unmarshal fails", func() {
			It("should return an error", func() {
				// Arrange
				// setup := ParameterStoreSetup{
				// 	Path:   "/go-boilerplate/",
				// 	Type:   "json",
				// 	Prefix: "",
				// }

				// // Invalid struct for unmarshaling
				// var invalidConfig string

				Skip("Skipping AWS Parameter Store test - requires refactoring for dependency injection")
			})
		})
	})
})

// Contoh refactor yang disarankan untuk membuat kode lebih testable
var _ = Describe("Configz with Dependency Injection (Suggested Refactor)", func() {
	Describe("LoadFromAWSParameterStore with injected ParameterStore", func() {
		Context("when using mock parameter store", func() {
			It("should demonstrate how to test with dependency injection", func() {
				// Contoh bagaimana kode bisa di-refactor untuk testing:
				//
				// type ParameterStore interface {
				//     GetAllParametersByPath(path string, recursive bool) (map[string]string, error)
				// }
				//
				// func LoadFromAWSParameterStoreWithDI(ps ParameterStore, setup ParameterStoreSetup, out any) error {
				//     params, err := ps.GetAllParametersByPath(setup.Path, true)
				//     if err != nil {
				//         return err
				//     }
				//     // ... rest of the logic
				// }

				Skip("This is just an example of how to refactor for better testability")
			})
		})
	})
})
