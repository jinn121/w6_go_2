package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type Note struct {
	ID      int      `json:"id"`
	Title   string   `json:"title"`
	Content string   `json:"content"`
	Tags    []string `json:"tags"`
}

var notes []*Note

func getNextID() int {
	if len(notes) == 0 {
		return 1
	}
	return notes[len(notes)-1].ID + 1
}

func createNote(w http.ResponseWriter, r *http.Request) {
	var note Note
	err := json.NewDecoder(r.Body).Decode(&note)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	note.ID = getNextID()
	notes = append(notes, &note)
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(note)
}

func getAllNotes(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(notes)
}

func getNoteByID(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		http.Error(w, "Invalid note ID", http.StatusBadRequest)
		return
	}

	for _, note := range notes {
		if note.ID == id {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(note)
			return
		}
	}

	http.Error(w, "Note not found", http.StatusNotFound)
}

func updateNote(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		http.Error(w, "Invalid note ID", http.StatusBadRequest)
		return
	}

	for _, note := range notes {
		if note.ID == id {
			err := json.NewDecoder(r.Body).Decode(note)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(note)
			return
		}
	}

	http.Error(w, "Note not found", http.StatusNotFound)
}

func deleteNote(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		http.Error(w, "Invalid note ID", http.StatusBadRequest)
		return
	}

	for i, note := range notes {
		if note.ID == id {
			notes = append(notes[:i], notes[i+1:]...)
			w.WriteHeader(http.StatusNoContent)
			return
		}
	}

	http.Error(w, "Note not found", http.StatusNotFound)
}

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/notes", getAllNotes).Methods("GET")
	router.HandleFunc("/notes", createNote).Methods("POST")
	router.HandleFunc("/notes/{id}", getNoteByID).Methods("GET")
	router.HandleFunc("/notes/{id}", updateNote).Methods("PUT")
	router.HandleFunc("/notes/{id}", deleteNote).Methods("DELETE")
	fmt.Println("Server listening on port 8080")
	http.ListenAndServe(":8080", router)
}
