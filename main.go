package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Result struct {
	ID    string
	Pwd   string
	Name  string
	Count int
}

type Body struct {
	ID   string
	Pwd  string
	Name string
}

func login(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var result Result
	var body Body
	var rtn string
	_ = json.NewDecoder(r.Body).Decode(&body)

	conn().Raw("SELECT count(id) FROM user_list WHERE id = ? AND pwd = ? AND delyn = 'N'", body.ID, body.Pwd).Scan(&result)
	fmt.Println(result)
	if result.Count != 0 {
		rtn = "suc"
	} else {
		rtn = "fail"
	}
	json.NewEncoder(w).Encode(rtn)
}

func register(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var result Result
	var body Body
	var rtn string
	_ = json.NewDecoder(r.Body).Decode(&body)

	conn().Raw("INSERT INTO user_list (id, pwd, name, delyn, dttm) VALUES (?, ?, ?, 'N', now()) returning id", body.ID, body.Pwd, body.Name).Scan(&result)
	if result.ID != "" {
		rtn = "suc"
	} else {
		rtn = "fail"
	}
	json.NewEncoder(w).Encode(rtn)
}

func idCheck(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	var result Result
	var rtn string

	conn().Raw("SELECT count(id) FROM user_list WHERE id = ?", params["id"]).Scan(&result)
	fmt.Println(result)
	if result.Count == 0 {
		rtn = "suc"
	} else {
		rtn = "fail"
	}
	json.NewEncoder(w).Encode(rtn)
}

func getProfile(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var result Result
	var body Body
	_ = json.NewDecoder(r.Body).Decode(&body)

	conn().Raw("SELECT id, pwd, name FROM user_list WHERE id = ? AND delyn = 'N'", body.ID).Scan(&result)
	if result.ID != "" {
		json.NewEncoder(w).Encode(result)
	} else {
		json.NewEncoder(w).Encode("fail")
	}

}

func getBoardList(w http.ResponseWriter, r *http.Request) {

}

func boardInsert(w http.ResponseWriter, r *http.Request) {

}

func boardUpdate(w http.ResponseWriter, r *http.Request) {

}

func boardDelete(w http.ResponseWriter, r *http.Request) {

}

func main() {
	r := mux.NewRouter()

	r.HandleFunc("/api/login", login).Methods("POST")
	r.HandleFunc("/api/register", register).Methods("POST")
	r.HandleFunc("/api/register/{id}", idCheck).Methods("GET")
	r.HandleFunc("/api/profile", getProfile).Methods("POST")

	r.HandleFunc("/api/board/list", getBoardList).Methods("POST")
	r.HandleFunc("/api/board/insert", boardInsert).Methods("POST")
	r.HandleFunc("/api/board/update", boardUpdate).Methods("POST")
	r.HandleFunc("/api/board/delete", boardDelete).Methods("DELETE")

	handler := cors.Default().Handler(r)

	log.Fatal(http.ListenAndServe(":8010", handler))
}

func conn() *gorm.DB {
	dsn := "host=34.64.255.101 user=postgres password=365365 dbname=postgres port=5432 sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn))

	if err != nil {
		fmt.Println(err.Error())
		panic("Failed to connect to database")
	}
	return db
}
