package config

import (
	"errors"
	"fmt"
	"github.com/manifoldco/promptui"
	bugyoclient "github.com/tomtwinkle/bugyo-client-go"
	"gopkg.in/yaml.v3"
	"log"
	"os"
	"path/filepath"
	"strconv"
)

const configFile = "bugyoclient.yaml"

type configYaml struct {
	TenantCode string `yaml:"tenant_code"`
	OBCiD      string `yaml:"obc_id"`
	Password   string `yaml:"password"`
}

type config struct {
	ConfigPath string
}

type Config interface {
	Init() (*bugyoclient.BugyoConfig, error)
}

func NewConfig() Config {
	return &config{ConfigPath: getConfigPath()}
}

func getConfigPath() string {
	execPath, err := os.Executable()
	if err != nil {
		log.Fatal(err)
	}
	execDirPath := filepath.Dir(execPath)
	return filepath.Join(execDirPath, configFile)
}

func (c *config) Init() (*bugyoclient.BugyoConfig, error) {
	if _, err := os.Stat(c.ConfigPath); err != nil {
		if !os.IsNotExist(err) {
			return nil, err
		}
		return c.writeConfig()
	}
	return c.readConfig()
}

func (c *config) readConfig() (*bugyoclient.BugyoConfig, error) {
	file, err := os.Open(c.ConfigPath)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	d := yaml.NewDecoder(file)
	var cfg configYaml
	if err := d.Decode(&cfg); err != nil {
		return nil, err
	}

	if cfg.TenantCode == "" {
		return nil, fmt.Errorf("tenant_code is Required [%s]", c.ConfigPath)
	}
	if cfg.OBCiD == "" {
		return nil, fmt.Errorf("obc_id is Required [%s]", c.ConfigPath)
	}
	if cfg.Password == "" {
		return nil, fmt.Errorf("password is Required [%s]", c.ConfigPath)
	}
	return &bugyoclient.BugyoConfig{
		TenantCode: cfg.TenantCode,
		OBCiD:      cfg.OBCiD,
		Password:   cfg.Password,
	}, nil
}

func (c *config) writeConfig() (*bugyoclient.BugyoConfig, error) {
	tenantCode, err := c.inputTenant()
	if err != nil {
		return nil, err
	}
	obcId, err := c.inputOBCiD()
	if err != nil {
		return nil, err
	}
	password, err := c.inputPassword()
	if err != nil {
		return nil, err
	}

	file, err := os.Create(c.ConfigPath)
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

func (c *config) inputTenant() (string, error) {
	validate := func(input string) error {
		if input == "" {
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

func (c *config) inputOBCiD() (string, error) {
	validate := func(input string) error {
		_, err := strconv.ParseFloat(input, 64)
		if err != nil {
			return errors.New("OBCiDは数字です")
		}
		if input == "" {
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

func (c *config) inputPassword() (string, error) {
	validate := func(input string) error {
		if input == "" {
			return errors.New("パスワードは必須です")
		}
		return nil
	}
	prompt := promptui.Prompt{
		Label:    fmt.Sprintf("パスワードを入力してください"),
		Validate: validate,
		Mask:     '*',
	}
	result, err := prompt.Run()
	if err != nil {
		fmt.Printf("%v\n", err.Error())
		return "", err
	}
	return result, nil
}
