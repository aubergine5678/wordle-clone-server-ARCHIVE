package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"
	"wordle-clone-server/db_client"

	"github.com/gorilla/mux"
)

type Response struct {
	Message string `json:"message"`
	Id      string `json:"id"`
}

type User struct {
	Id       int64     `json:"id"`
	Username string    `json:"username"`
	Forename string    `json:"forename"`
	Surname  string    `json:"surname"`
	Dob      time.Time `json:"dob"`
}

type Game struct {
	Id       int64  `json:"id"`
	Username string `json:"username"`
	Attempts int64  `json:"attempts"`
	Time     int64  `json:"time"`
	GameMode string `json:"gamemode"`
}

func main() {
	db_client.InitialiseDBConnection()

	log.Println("Starting server...")
	router := mux.NewRouter()
	router.Use(mux.CORSMethodMiddleware(router))
	// Test ping route
	router.HandleFunc("/ping", ping).Methods(http.MethodGet, http.MethodOptions)

	// Games routes
	router.HandleFunc("/games", gameHandler).Methods(http.MethodGet, http.MethodPost, http.MethodOptions)
	router.HandleFunc("/games/{id}", getGameById).Methods(http.MethodGet, http.MethodOptions)

	// Users routes
	router.HandleFunc("/users", getAllUsers).Methods(http.MethodGet, http.MethodOptions)
	router.HandleFunc("/users/{id}", getUserById).Methods(http.MethodGet, http.MethodOptions)

	log.Println("Server started on port :8080")
	err := http.ListenAndServe(":8080", router)
	if err != nil {
		log.Fatal(err)
	}
}

func ping(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	log.Println("Request for /ping received")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(Response{Message: "Server is up and running"})
}

func gameHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		createGame(w, r)
	} else if r.Method == http.MethodGet {
		getAllGames(w, r)
	}
}

func getAllGames(w http.ResponseWriter, r *http.Request) {
	var games []Game
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	log.Println("Request for /games received")
	rows, err := db_client.DBClient.Query("SELECT id, username, game_type, attempts, time FROM games")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(Response{Message: "Internal server error"})
		return
	}

	for rows.Next() {
		var game Game
		err := rows.Scan(&game.Id, &game.Username, &game.GameMode, &game.Attempts, &game.Time)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(Response{Message: "Internal server error"})
			return
		}
		games = append(games, game)
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(games)
}

func getGameById(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	log.Println("Request for /games/{id} received")
	req := mux.Vars(r)
	id, err := strconv.Atoi(req["id"])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(Response{Message: "Invalid ID"})
		return
	}

	var games []Game
	rows, err := db_client.DBClient.Query("SELECT id, username, game_type, attempts, time FROM games")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(Response{Message: "Internal server error"})
		return
	}

	for rows.Next() {
		var game Game
		err := rows.Scan(&game.Id, &game.Username, &game.GameMode, &game.Attempts, &game.Time)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(Response{Message: "Internal server error"})
			return
		}
		games = append(games, game)
	}

	for _, game := range games {
		if game.Id == int64(id) {
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(game)
			return
		}
	}

	w.WriteHeader(http.StatusNotFound)
	json.NewEncoder(w).Encode(Response{Message: "Game record with specified id not found"})
}

func createGame(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	log.Println("POST request for /games received")

	var reqBody Game
	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		json.NewEncoder(w).Encode(Response{Message: "Invalid request body"})
		return
	}

	if reqBody.Username == "" || reqBody.Attempts <= 0 || reqBody.Time <= 0 || reqBody.GameMode == "" {
		w.WriteHeader(http.StatusUnprocessableEntity)
		json.NewEncoder(w).Encode(Response{Message: "Invalid request body"})
		return
	}

	if reqBody.GameMode != "game5" && reqBody.GameMode != "game7" {
		w.WriteHeader(http.StatusUnprocessableEntity)
		json.NewEncoder(w).Encode(Response{Message: "Invalid game mode"})
	}

	res, err := db_client.DBClient.Exec("INSERT INTO games (username, attempts, time, game_type) VALUES (?, ?, ?, ?)", reqBody.Username, reqBody.Attempts, reqBody.Time, reqBody.GameMode)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(Response{Message: "Internal server error"})
		return
	}

	id, err := res.LastInsertId()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(Response{Message: "Internal server error"})
		return
	}

	w.WriteHeader(http.StatusCreated)
	reqBody.Id = id
	json.NewEncoder(w).Encode(reqBody)
}

func getAllUsers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	log.Println("Request for /users received")
	var users []User
	rows, err := db_client.DBClient.Query("SELECT id, username, firstname, surname, dob FROM users")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(Response{Message: "Internal server error"})
		return
	}

	for rows.Next() {
		var user User
		err := rows.Scan(&user.Id, &user.Username, &user.Forename, &user.Surname, &user.Dob)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(Response{Message: "Internal server error"})
			return
		}
		users = append(users, user)
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(users)
}

func getUserById(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	log.Println("Request for /users/{id} received")

	req := mux.Vars(r)
	id, err := strconv.Atoi(req["id"])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(Response{Message: "Invalid ID"})
		return
	}

	var users []User
	rows, err := db_client.DBClient.Query("SELECT id, username, firstname, surname, dob FROM users")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(Response{Message: "Internal server error"})
		return
	}

	for rows.Next() {
		var user User
		err := rows.Scan(&user.Id, &user.Username, &user.Forename, &user.Surname, &user.Dob)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(Response{Message: "Internal server error"})
			return
		}
		users = append(users, user)
	}

	for _, user := range users {
		if user.Id == int64(id) {
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(user)
			return
		}
	}

	w.WriteHeader(http.StatusNotFound)
	json.NewEncoder(w).Encode(Response{Message: "User with specified id not found"})
}
