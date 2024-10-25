package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
)

type Film struct {
	ID         int     `json:"id"`
	Judul      string  `json:"judul"`
	TahunRilis int     `json:"tahun_rilis"`
	Genre      string  `json:"genre"`
	Rating     float64 `json:"rating"`
}

var films []Film

func main() {
	loadData()
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		respond(w, "Selamat datang di API Film! Gunakan /films untuk mengakses data film.", http.StatusOK)
	})
	http.HandleFunc("/films", handleFilms)
	http.HandleFunc("/films/", handleFilmByID)

	fmt.Println("Server berjalan di port 8081")
	log.Fatal(http.ListenAndServe(":8081", nil))
}

func loadData() {
	file, err := os.ReadFile("films.json")
	if err == nil {
		json.Unmarshal(file, &films)
	}
}

func saveData() {
	data, _ := json.MarshalIndent(films, "", "  ")
	os.WriteFile("films.json", data, 0644)
}

func handleFilms(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		respondJSON(w, films, http.StatusOK)
	} else if r.Method == http.MethodPost {
		var newFilm Film
		if err := json.NewDecoder(r.Body).Decode(&newFilm); err == nil {
			newFilm.ID = len(films) + 1
			films = append(films, newFilm)
			saveData()
			respondJSON(w, newFilm, http.StatusCreated)
		} else {
			respond(w, "Format data tidak valid", http.StatusBadRequest)
		}
	} else {
		respond(w, "Metode tidak didukung", http.StatusMethodNotAllowed)
	}
}

func handleFilmByID(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Path[len("/films/"):])
	if err != nil {
		respond(w, "ID tidak valid", http.StatusBadRequest)
		return
	}

	for i, film := range films {
		if film.ID == id {
			switch r.Method {
			case http.MethodGet:
				respondJSON(w, film, http.StatusOK)
			case http.MethodPut:
				var updatedFilm Film
				if json.NewDecoder(r.Body).Decode(&updatedFilm) == nil {
					updatedFilm.ID = id
					films[i] = updatedFilm
					saveData()
					respondJSON(w, updatedFilm, http.StatusOK)
				} else {
					respond(w, "Format data tidak valid", http.StatusBadRequest)
				}
			case http.MethodDelete:
				films = append(films[:i], films[i+1:]...)
				saveData()
				respond(w, "data berhasil dihapus", http.StatusOK)
			default:
				respond(w, "Metode tidak didukung", http.StatusMethodNotAllowed)
			}
			return
		}
	}
	respond(w, "Film tidak ditemukan", http.StatusNotFound)
}

func respond(w http.ResponseWriter, message string, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(map[string]string{"message": message})
}

func respondJSON(w http.ResponseWriter, data interface{}, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(data)
}
