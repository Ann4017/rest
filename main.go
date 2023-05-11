package main

import (
	"fmt"
	"net/http"
	"rest/db"

	"rest/handlers"

	"github.com/gorilla/mux"
)

func main() {
	db := &db.C_db{}
	h := &handlers.C_user_handler{
		C_db: db,
	}

	err := db.Load_config("config/config.ini")
	if err != nil {
		fmt.Println("Load_config err:", err)
	}

	err = db.Connect_db()
	if err != nil {
		fmt.Println("Connect_db err:", err)
	}

	r := mux.NewRouter()

	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "main")
	})
	r.HandleFunc("/users", h.Select_users)
	r.HandleFunc("/user/{id}", h.Select_user).Methods("GET")
	r.HandleFunc("/user", h.Insert_user).Methods("POST")
	r.HandleFunc("/user/{id}", h.Update_user).Methods("PUT")
	r.HandleFunc("/user/{id}", h.Delete_user).Methods("DELETE")

	http.ListenAndServe(":8000", r)
}
