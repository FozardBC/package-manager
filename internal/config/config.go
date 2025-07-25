package config

import (
	"encoding/json"
	"log"
	"os"
	"strings"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
	"gopkg.in/yaml.v3"
)

type Config struct {
	SshHost    string `env:"PM_SSH_HOST"`
	SshUser    string `env:"PM_SSH_USER"`
	SshPort    string `env:"PM_SSH_PORT"`
	SshKeyPath string `env:"PM_SSH_KEY_PATH"`
}

type Packet struct {
	Name    string     `json:"name" yaml:"name"`
	Ver     string     `json:"ver" yaml:"ver"`
	Targets []Target   `json:"targets" yaml:"targets"`
	Packets Dependency `json:"packets,omitempty" yaml:"packets,omitempty"`
}

type Target struct {
	Path    string `json:"path" yaml:"path"`
	Exclude string `json:"exclude,omitempty" yaml:"exclude,omitempty"`
}

type Dependency struct {
	Name string `json:"name" yaml:"name"`
	Ver  string `json:"ver,omitempty" yaml:"ver,omitempty"`
}

type PackagesFile struct {
	Packages []Dependency `json:"packages" yaml:"packages"`
}

func MustLoad(path string) *Config {
	if err := godotenv.Load(); err != nil {
		panic(err)
	}

	cfg := Config{}
	if err := cleanenv.ReadEnv(&cfg); err != nil {
		help, _ := cleanenv.GetDescription(cfg, nil)
		log.Print(help)
		log.Fatal(err)
	}

	return &cfg
}

func LoadPacketConfig(filename string, out interface{}) error {

	data, err := os.ReadFile(filename)
	if err != nil {
		return err
	}

	if strings.HasSuffix(filename, ".yaml") || strings.HasSuffix(filename, ".yml") {
		return yaml.Unmarshal(data, out)
	}

	return json.Unmarshal(data, out)
}
