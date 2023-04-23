package conf

import (
	"errors"
	"gopkg.in/yaml.v2"
	"io"
	"os"
)

type Config struct {
	Addr          string `yaml:"Addr"`
	Database      string `yaml:"Database"`
	TokenRelation struct {
		Magnification float64 `yaml:"Magnification"`
		MinToken      int64   `yaml:"MinToken"`
	} `yaml:"TokenRelation"`
	ServerConfig struct {
		Addr string `yaml:"Addr"`
	} `yaml:"ServerConfig"`
	Auth struct {
		Secret     string `yaml:"Secret"`
		ExpireTime int64  `yaml:"ExpireTime"`
	} `yaml:"Auth"`
}

func (c *Config) SaveConfig(writer io.Writer) error {
	if c == nil {
		return errors.New("config is nil")
	}
	encoder := yaml.NewEncoder(writer)
	defer func() {
		err := encoder.Close()
		if err != nil {
			return
		}
	}()
	err := encoder.Encode(c)
	if err != nil {
		return err
	}
	return nil
}

func (c *Config) SaveConfigToFile(path string) error {
	if c == nil {
		return errors.New("config is nil")
	}
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer func() {
		err := file.Close()
		if err != nil {
			return
		}
	}()
	encoder := yaml.NewEncoder(file)
	defer func(encoder *yaml.Encoder) {
		err := encoder.Close()
		if err != nil {
			panic(err)
		}
	}(encoder)
	err = encoder.Encode(c)
	if err != nil {
		return err
	}
	return nil
}

func ParseConfig(reader io.Reader) (*Config, error) {
	decoder := yaml.NewDecoder(reader)
	var conf *Config
	err := decoder.Decode(&conf)
	if err != nil {
		return nil, err
	}
	return conf, nil
}

func ParseConfigFromFile(path string) (*Config, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, errors.New(path + "打开失败:" + err.Error())
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			panic(err)
		}
	}(file)
	conf, err := ParseConfig(file)
	if err != nil {
		return nil, err
	}
	return conf, nil
}
