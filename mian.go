package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type User struct {
	Id    int    `json:"id"`
	Name  string `json:"username"`
	Age   int    `json:"age"`
	Email string `json:"email"`
}

var users = []User{
	{
		Id:    1,
		Name:  "Mila",
		Age:   25,
		Email: "kamal3@gmail.com",
	},
	{
		Id:    2,
		Name:  "Koilla",
		Age:   28,
		Email: "millat@gmail.com",
	},
}

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("/", rootHandler)
	mux.HandleFunc("/health", healthHandler)
	mux.HandleFunc("POST /createUser", createUserHandler)
	mux.HandleFunc("GET /users", getUserHandler)
	mux.HandleFunc("GET /users/{id}", getUserByIDHandler)
	mux.HandleFunc("PUT /users/{id}", updateUserHandler)
	mux.HandleFunc("DELETE /users/{id}", deleteUserHandler)

	fmt.Println("Server is running on 5000")

	err := http.ListenAndServe(":5000", mux)
	if err != nil {
		fmt.Println("Error starting server:", err)
	}
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, World!")
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "OK! server is healthy")
}

func createUserHandler(w http.ResponseWriter, r *http.Request) {
	var newUser User

	err := json.NewDecoder(r.Body).Decode(&newUser)
	if err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}
	fmt.Println(newUser)
    newUser.Id = len(users) + 1
	users = append(users, newUser)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	json.NewEncoder(w).Encode(newUser)
}

func getUserHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(users)
}

func getUserByIDHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Path[len("/users/"):]

	for _, user := range users {
		if fmt.Sprintf("%d", user.Id) == id {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(user)
			return
		}
	}

	http.Error(w, "User not found", http.StatusNotFound)
}

func updateUserHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Path[len("/users/"):]

	for i, user := range users {
		if fmt.Sprintf("%d", user.Id) == id {
			var updatedUser User
			err := json.NewDecoder(r.Body).Decode(&updatedUser)
			if err != nil {
				http.Error(w, "Invalid JSON", http.StatusBadRequest)
				return
			}
			updatedUser.Id = user.Id
			users[i] = updatedUser
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(updatedUser)
			return
		}
	}

	http.Error(w, "User not found", http.StatusNotFound)
}

func deleteUserHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Path[len("/users/"):]

	for i, user := range users {
		if fmt.Sprintf("%d", user.Id) == id {
			users = append(users[:i], users[i+1:]...)
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(user)
			return
		}
	}

	http.Error(w, "User not found", http.StatusNotFound)
}