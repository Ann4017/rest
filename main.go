package main

import (
	"fmt"
	"log"
	"net/http"
	"rest/db"

	"rest/handlers"

	"github.com/gorilla/mux"
)

type Server struct {
	PC_db *db.C_db
	PC_h  *handlers.C_user_handler
}

func (s *Server) Init(_s_ini_path string) error {
	s.PC_db = &db.C_db{}
	s.PC_h = &handlers.C_user_handler{
		C_db: s.PC_db,
	}

	err := s.PC_db.Load_config(_s_ini_path)
	if err != nil {
		fmt.Println("Load_config err:", err)
		return err
	}

	err = s.PC_db.Connect_db()
	if err != nil {
		fmt.Println("Connect_db err:", err)
		return err
	}

	return nil
}

func (s *Server) Run() {
	r := mux.NewRouter()

	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "main")
	})
	r.HandleFunc("/users", s.PC_h.Select_users)
	r.HandleFunc("/user/{id}", s.PC_h.Select_user).Methods("GET")
	r.HandleFunc("/user", s.PC_h.Insert_user).Methods("POST")
	r.HandleFunc("/user/{id}", s.PC_h.Update_user).Methods("PUT")
	r.HandleFunc("/user/{id}", s.PC_h.Delete_user).Methods("DELETE")

	http.ListenAndServe(":8000", r)
}

func main() {
	s := Server{}
	err := s.Init("config/config.ini")
	if err != nil {
		log.Fatal(err)
	}

	s.Run()
}
