package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"rest/db"
	"rest/models"

	"github.com/gorilla/mux"
)

type C_user_handler struct {
	C_db *db.C_db
}

func (h *C_user_handler) Select_user(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	user := models.C_user{}

	err = h.C_db.PC_sql_db.QueryRow("select * from users where id = ?", id).Scan(&user.I_id, &user.S_name, &user.S_email)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(user)
}

func (h *C_user_handler) Select_users(w http.ResponseWriter, r *http.Request) {
	rows, err := h.C_db.PC_sql_db.Query("select * from users")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	users := []models.C_user{}

	for rows.Next() {
		user := models.C_user{}
		err := rows.Scan(&user.I_id, &user.S_name, &user.S_email)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		users = append(users, user)
	}

	json.NewEncoder(w).Encode(users)
}

func (h *C_user_handler) Insert_user(w http.ResponseWriter, r *http.Request) {
	user := models.C_user{}
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	row, err := h.C_db.PC_sql_db.Exec("insert into users (name, email) value (?, ?)", user.S_name, user.S_email)
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

func (h *C_user_handler) Update_user(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	user := models.C_user{}
	err = json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	_, err = h.C_db.PC_sql_db.Exec("update users set name = ?, email = ? where id = ?", user.S_name, user.S_email, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(user)
}

func (h *C_user_handler) Delete_user(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	_, err = h.C_db.PC_sql_db.Exec("delete from users where id = ?", id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
