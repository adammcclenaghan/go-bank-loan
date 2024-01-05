/*
When I originally worked on this project in 2022, it was for an interview.

The HTTP "Bank API" was supplied by them as a Docker image. However, since
it was part of their interview process, I cannot include it in my docker
compose anymor. 

As I write this, it is 2024, and I would like to make this project public as
it showcases a little of the Go code I have written. However I need to remove
references to the company I interviewed with.

So, over my lunch break I have lazily asked ChatGPT to write all of this logic
... Magic :) 

*/
package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/mux"
)

// Application represents the application data structure.
type Application struct {
	ID        string `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Status    string `json:"status"`
}

var (
	applications     = make(map[string]Application)
	applicationsLock sync.Mutex
)

func main() {
	rand.Seed(time.Now().UnixNano())
	r := mux.NewRouter()

	r.HandleFunc("/api/applications", CreateApplication).Methods("POST")
	r.HandleFunc("/api/jobs", GetApplicationStatus).Methods("GET").Queries("application_id", "{id}")

	http.Handle("/", r)

	port := 8000
	fmt.Printf("Server is running on port %d...\n", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), nil))
	log.Printf("The server is shutting down")
}

// CreateApplication handles the creation of a new application.
func CreateApplication(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("Creating new application")
	var newApp Application
	err := json.NewDecoder(r.Body).Decode(&newApp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	applicationsLock.Lock()
	defer applicationsLock.Unlock()

	// Check if the application ID is already used
	if _, exists := applications[newApp.ID]; exists {
		http.Error(w, "the application ID is already used.", http.StatusBadRequest)
		return
	}

	// Set the default status to 'pending' for new applications
	newApp.Status = "pending"

	// Add the application to the map
	applications[newApp.ID] = newApp
	fmt.Printf("Added app %s to map", newApp.ID)

	// Simulate processing time between 5-20 seconds
	go func(id string) {
		time.Sleep(time.Duration((5 + rand.Intn(16))) * time.Second)
		status := getRandomStatus()

		applicationsLock.Lock()
		defer applicationsLock.Unlock()

		// Update the status of the existing application
		if app, ok := applications[id]; ok {
			app.Status = status
			applications[id] = app
		}
	}(newApp.ID)

	// Return the response with a 201 Created status code
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newApp)
}

// GetApplicationStatus handles the retrieval of an application's status based on the application ID.
func GetApplicationStatus(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	applicationID := params["id"]
	fmt.Printf("Getting app %s from map", applicationID)

	applicationsLock.Lock()
	defer applicationsLock.Unlock()

	// Retrieve the application from the map
	app, ok := applications[applicationID]
	if !ok {
		http.Error(w, "application not found", http.StatusNotFound)
		return
	}

	// Return the response with only necessary fields
	response := struct {
		ID     string `json:"id"`
		Status string `json:"status"`
	}{
		ID:     app.ID,
		Status: app.Status,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// getRandomStatus generates a random status for an application.
func getRandomStatus() string {
	statuses := []string{"completed", "rejected"}
	return statuses[rand.Intn(len(statuses))]
}

