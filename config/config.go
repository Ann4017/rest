package config

import (
	"os"

	"gopkg.in/ini.v1"
)

type Config struct {
	S_db_user string
	S_db_pwd  string
	S_db_host string
	S_db_name string
}

func Load_config(_s_ini_path string) (*Config, error) {
	cfg, err := ini.Load(_s_ini_path)
	if err != nil && os.IsNotExist(err) {
		cfg = ini.Empty()

		sec, _ := cfg.NewSection("database")
		sec.NewKey("user", "default")
		sec.NewKey("pwd", "default")
		sec.NewKey("host_port", "host_default:port_default")
		sec.NewKey("DB_name", "defalut")

		err = cfg.SaveTo(_s_ini_path)
		if err != nil {
			return nil, err
		}
		return nil, err
	}

	config := &Config{}
	section := cfg.Section("database")
	config.S_db_user = section.Key("user").String()
	config.S_db_pwd = section.Key("pwd").String()
	config.S_db_host = section.Key("host_port").String()
	config.S_db_name = section.Key("DB_name").String()

	return config, nil
}
