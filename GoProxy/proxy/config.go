package proxy

import (
	"errors"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"strings"
)

type Config struct {
	Mode       string `yaml:"mode"`
	Auth       bool   `yaml:"auth"`
	Server     string `yaml:"server"`
	Timeout    int    `yaml:"timeout"`
	Listen     string `yaml:"listen"`
	CertFile   string `yaml:"cert_file"`
	KeyFile    string `yaml:"key_file"`
	XorKey     int    `yaml:"xor_key"`
	UserFile   string `yaml:"user_file"`
	RuntimeLog string `yaml:"runtime_log"`
	UserLog    string `yaml:"user_log"`
	AccessLog  string `yaml:"access_log"`
	HTTPProxy  string `yaml:"http_proxy"`
	Username   string `yaml:"username"`
	Password   string `yaml:"password"`
}

type UserItem struct {
	Mac      string `yaml:"mac"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
}

type UserConfig struct {
	Mac  bool       `yaml:"mac"`
	User []UserItem `yaml:"user"`
}

var _config *Config
var _user *UserConfig

func LoadConfig(filename string) (config *Config, err error) {
	config = new(Config)
	in, er := ioutil.ReadFile(filename)
	if er != nil {
		err = er
		return
	}
	if er := yaml.Unmarshal(in, &config); er != nil {
		err = er
	}
	_config = config
	return
}

func LoadUserConfig(filename string) (config *UserConfig, err error) {
	config = new(UserConfig)
	in, er := ioutil.ReadFile(filename)
	if er != nil {
		err = er
		return
	}
	if er := yaml.Unmarshal(in, &config); er != nil {
		err = er
	}
	_user = config
	return
}

func GetConfig() *Config {
	return _config
}

func GetUserConfig() *UserConfig {
	return _user
}

func AuthUser(username, password string, macList []string) error {
	if _user == nil {
		log.Println("user is null")
		return nil
	}
	for _, user := range _user.User {
		if strings.Compare(user.Username, username) == 0 {
			if strings.Compare(user.Password, password) != 0 {
				return errors.New("auth: password")
			}
			if len(user.Mac) > 0 {
				for _, mac := range macList {
					if strings.Contains(strings.ToUpper(user.Mac), strings.ToUpper(mac)) == true {
						return nil
					}
				}
				return errors.New("auth: mac")
			}
			return nil
		}
	}
	return errors.New("auth: username")
}
