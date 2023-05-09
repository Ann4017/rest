package db

import (
	"database/sql"
	"fmt"
	"rest/config"
)

func Connect_db(_s_DBMS string, _cfg *config.Config) (*sql.DB, error) {
	source := fmt.Sprintf("%s:%s@tcp(%s)/%s", _cfg.S_db_user, _cfg.S_db_pwd, _cfg.S_db_host, _cfg.S_db_name)
	db, err := sql.Open(_s_DBMS, source)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}
