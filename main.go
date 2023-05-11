package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/gorilla/mux"

	_ "github.com/go-sql-driver/mysql"
	"gopkg.in/ini.v1"
)

type C_user struct {
	I_id    int    `json:"id"`
	S_name  string `json:"name"`
	S_email string `json:"email"`
}

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

func (c *C_db) Select_user(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	user := C_user{}

	err = c.PC_sql_db.QueryRow("select * from users where id = ?", id).Scan(&user.I_id, &user.S_name, &user.S_email)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(user)
}

func (c *C_db) Select_users(w http.ResponseWriter, r *http.Request) {
	rows, err := c.PC_sql_db.Query("select * from users")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	users := []C_user{}

	for rows.Next() {
		user := C_user{}
		err := rows.Scan(&user.I_id, &user.S_name, &user.S_email)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		users = append(users, user)
	}

	json.NewEncoder(w).Encode(users)
}

func (c *C_db) Insert_user(w http.ResponseWriter, r *http.Request) {
	user := C_user{}
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	row, err := c.PC_sql_db.Exec("insert into users (name, email) value (?, ?)", user.S_name, user.S_email)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	id, err := row.LastInsertId()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	user.I_id = int(id)
	json.NewEncoder(w).Encode(user)
}

func (c *C_db) Update_user(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	user := C_user{}
	err = json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	_, err = c.PC_sql_db.Exec("update users set name = ?, email = ? where id = ?", user.S_name, user.S_email, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(user)
}

func (c *C_db) Delete_user(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	_, err = c.PC_sql_db.Exec("delete from users where id = ?", id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func main() {
	db := C_db{}
	err := db.Load_config("config.ini")
	if err != nil {
		fmt.Print("load-config err", err)
	}

	err = db.Connect_db()
	if err != nil {
		fmt.Print("connect-db err", err)
	}

	r := mux.NewRouter()

	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "main")
	})
	r.HandleFunc("/users", db.Select_users)
	r.HandleFunc("/user/{id}", db.Select_user).Methods("GET")
	r.HandleFunc("/user", db.Insert_user).Methods("POST")
	r.HandleFunc("/user/{id}", db.Update_user).Methods("PUT")
	r.HandleFunc("/user/{id}", db.Delete_user).Methods("DELETE")

	http.ListenAndServe(":8000", r)
}
