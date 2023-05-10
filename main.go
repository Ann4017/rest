package main

import (
	"net/http"
	"rest/config"
	"rest/db"
	"rest/handlers"
)

func Init() error {
	cfg, err := config.Load_config("config.ini")
	if err != nil {
		return err
	}

	sql_db, err := db.Connect_db("mysql", cfg)
	if err != nil {
		return err
	}

	mux := http.NewServeMux()

	mux.HandleFunc("/create", handlers.Insert_handler)
	mux.HandleFunc("/read", handlers.Select_handler)
	mux.HandleFunc("/update", handlers.Update_handler)
	mux.HandleFunc("/delete", handlers.Delete_handler)

	http.ListenAndServe(":3000", mux)
}

func main() {

}
