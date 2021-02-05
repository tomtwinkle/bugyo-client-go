package config

import (
	"errors"
	"fmt"
	"github.com/manifoldco/promptui"
	bugyoclient "github.com/tomtwinkle/bugyo-client-go"
	"gopkg.in/yaml.v3"
	"os"
	"strconv"
)

const configFile = "bugyoclient.yaml"

type configYaml struct {
	TenantCode string `yaml:"tenant_code"`
	OBCiD      string `yaml:"obc_id"`
	Password   string `yaml:"password"`
}

type config struct{}

type Config interface {
	Init() (*bugyoclient.BugyoConfig, error)
}

func NewConfig() Config {
	return &config{}
}

func (s *config) Init() (*bugyoclient.BugyoConfig, error) {
	_, err := os.Stat(configFile)
	if os.IsNotExist(err) {
		return s.writeConfig()
	}
	return s.readConfig()
}

func (s *config) readConfig() (*bugyoclient.BugyoConfig, error) {
	file, err := os.Open(configFile)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	d := yaml.NewDecoder(file)
	var cfg configYaml
	if err := d.Decode(&cfg); err != nil {
		return nil, err
	}
	return &bugyoclient.BugyoConfig{
		TenantCode: cfg.TenantCode,
		OBCiD:      cfg.OBCiD,
		Password:   cfg.Password,
	}, nil
}

func (s *config) writeConfig() (*bugyoclient.BugyoConfig, error) {
	tenantCode, err := s.inputTenant()
	if err != nil {
		return nil, err
	}
	obcId, err := s.inputOBCiD()
	if err != nil {
		return nil, err
	}
	password, err := s.inputPassword()
	if err != nil {
		return nil, err
	}

	file, err := os.Create(configFile)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	e := yaml.NewEncoder(file)
	defer e.Close()
	cfg := configYaml{
		TenantCode: tenantCode,
		OBCiD:      obcId,
		Password:   password,
	}
	if err := e.Encode(&cfg); err != nil {
		return nil, err
	}
	return &bugyoclient.BugyoConfig{
		TenantCode: cfg.TenantCode,
		OBCiD:      cfg.OBCiD,
		Password:   cfg.Password,
	}, nil
}

func (s *config) inputTenant() (string, error) {
	validate := func(input string) error {
		if len(input) == 0 {
			return errors.New("テナントコードは必須です")
		}
		return nil
	}
	prompt := promptui.Prompt{
		Label:    fmt.Sprintf("テナントコードを入力してください"),
		Validate: validate,
	}
	result, err := prompt.Run()
	if err != nil {
		fmt.Printf("%v\n", err.Error())
		return "", err
	}
	return result, nil
}

func (s *config) inputOBCiD() (string, error) {
	validate := func(input string) error {
		_, err := strconv.ParseFloat(input, 64)
		if err != nil {
			return errors.New("OBCiDは数字です")
		}
		if len(input) == 0 {
			return errors.New("OBCiDは必須です")
		}
		return nil
	}
	prompt := promptui.Prompt{
		Label:    fmt.Sprintf("OBCiDを入力してください"),
		Validate: validate,
	}
	result, err := prompt.Run()
	if err != nil {
		fmt.Printf("%v\n", err.Error())
		return "", err
	}
	return result, nil
}

func (s *config) inputPassword() (string, error) {
	validate := func(input string) error {
		if len(input) == 0 {
			return errors.New("パスワードは必須です")
		}
		return nil
	}
	prompt := promptui.Prompt{
		Label:    fmt.Sprintf("パスワードを入力してください"),
		Validate: validate,
	}
	result, err := prompt.Run()
	if err != nil {
		fmt.Printf("%v\n", err.Error())
		return "", err
	}
	return result, nil
}
