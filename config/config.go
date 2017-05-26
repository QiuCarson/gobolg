package config

import (
	"encoding/json"
	"io/ioutil"
	"log"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Title       string
	Keywords    string
	Description string
	Blogname    string
}
type Mysqlstruct struct {
	Username string
	Password string
	Mysql    string
	Prefix   string
}

func Conf() Config {
	var c Config

	/*c.Title = models.GetOption("blogname")
	c.Keywords = models.GetOption("d_keywords")
	c.Description = models.GetOption("blogdescription")
	c.Blogname = models.GetOption("blogname")*/
	return c
}

/**

 */
func GetYamlConfig(path string) map[interface{}]interface{} {
	data, err := ioutil.ReadFile(path)
	m := make(map[interface{}]interface{})
	if err != nil {
		log.Fatal(err.Error())
	}
	err = yaml.Unmarshal([]byte(data), &m)
	return m
}

/**
json 配置文件解析
*/
func MysqlConfig() Mysqlstruct {
	var m Mysqlstruct

	configFile, err := ioutil.ReadFile("config/config.json")
	if err != nil {
		log.Fatal(err.Error())
	}
	json.Unmarshal([]byte(configFile), &m)
	return m
}
