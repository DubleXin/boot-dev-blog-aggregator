package config

import "fmt"

const FILENAME = ".gatorconfig.json"

type Config struct {
	DBUrl       string `json:"db_url"`
	CurrentUser string `json:"current_user_name"`
}

func (c *Config) SetUser(current_user_name string) error {
	if len(current_user_name) == 0 {
		return fmt.Errorf("Cannot set empty user name. Expected not 0 length user name, but got \"%s\"", current_user_name)
	}
	c.CurrentUser = current_user_name
	return WriteConfig(c)
}
