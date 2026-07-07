package conf

import (
	"fmt"
	"os"
	"path"

	"github.com/pelletier/go-toml/v2"
)

type Config struct {
	DB *DB
}

func LoadConfig(dir string, name string) (*Config, error) {
	filePath := path.Join(dir, fmt.Sprintf("%s.toml", name))
	fileContent, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	var c Config
	err = toml.Unmarshal(fileContent, &c)
	if err != nil {
		return nil, err
	}

	return &c, nil
}
