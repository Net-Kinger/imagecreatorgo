package conf

import (
	"errors"
	"gopkg.in/yaml.v2"
	"io"
	"os"
)

var Conf Config

type Config struct {
	Addr          string `yaml:"Addr"`
	Database      string `yaml:"Database"`
	TokenRelation struct {
		Magnification float64 `yaml:"Magnification"`
		MinToken      int64   `yaml:"MinToken"`
	} `yaml:"TokenRelation"`
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
	defer encoder.Close()
	encoder.Encode(c)
	return nil
}

func ParseConfig(reader io.Reader) error {
	decoder := yaml.NewDecoder(reader)
	err := decoder.Decode(&Conf)
	if err != nil {
		return err
	}
	return nil
}

func ParseConfigFromFile(path string) error {
	file, err := os.Open(path)
	if err != nil {
		return errors.New(path + "打开失败:" + err.Error())
	}
	defer file.Close()
	err = ParseConfig(file)
	if err != nil {
		return err
	}
	return nil
}
