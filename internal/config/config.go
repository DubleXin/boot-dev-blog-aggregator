package config

const FILENAME = ".gatorconfig.json"

type Config struct {
	DBUrl       string `json:"db_url"`
	CurrentUser string `json:"current_user_name"`
}

func (c *Config) SetUser(current_user_name string) error {
	c.CurrentUser = current_user_name
	return WriteConfig(c)
}
