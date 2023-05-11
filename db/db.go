package db

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"gopkg.in/ini.v1"
)

type C_db struct {
	S_db_user   string
	S_db_pwd    string
	S_db_host   string
	S_db_name   string
	S_db_engine string
	PC_sql_db   *sql.DB
}

func (c *C_db) Load_config(_s_ini_path string) error {
	cfg, err := ini.Load(_s_ini_path)
	if err != nil && os.IsNotExist(err) {
		cfg = ini.Empty()

		sec, _ := cfg.NewSection("database")
		sec.NewKey("user", "default")
		sec.NewKey("pwd", "default")
		sec.NewKey("host_port", "host_default:port_default")
		sec.NewKey("DB_name", "defalut")
		sec.NewKey("DBMS", "defalut")

		err = cfg.SaveTo(_s_ini_path)
		if err != nil {
			return err
		}
		return err
	}

	section := cfg.Section("database")
	c.S_db_user = section.Key("user").String()
	c.S_db_pwd = section.Key("pwd").String()
	c.S_db_host = section.Key("host_port").String()
	c.S_db_name = section.Key("DB_name").String()
	c.S_db_engine = section.Key("DBMS").String()

	return nil
}

func (c *C_db) Connect_db() error {
	source := fmt.Sprintf("%s:%s@tcp(%s)/%s", c.S_db_user, c.S_db_pwd, c.S_db_host, c.S_db_name)
	db, err := sql.Open(c.S_db_engine, source)
	if err != nil {
		return err
	}

	c.PC_sql_db = db

	return nil
}
