package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"sync"
)

type FormData struct {
	Prenom string  `json:"prenom"`
	Poids  float64 `json:"poids"`
	Taille float64 `json:"taille"`
}

var (
	dataFile = "data.json"
	mu       sync.Mutex
)

func saveData(entry FormData) error {
	mu.Lock()
	defer mu.Unlock()

	// Lire les données existantes
	var data = make(map[string]FormData)
	file, err := os.Open(dataFile)
	if err == nil {
		defer file.Close()
		json.NewDecoder(file).Decode(&data)
	} else {
		data = make(map[string]FormData)
	}

	// Ajouter ou mettre à jour l'entrée
	data[entry.Prenom] = entry

	// Réécrire le fichier
	file, err = os.Create(dataFile)
	if err != nil {
		return err
	}
	defer file.Close()

	return json.NewEncoder(file).Encode(data)
}

func formHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Méthode non autorisée", http.StatusMethodNotAllowed)
		return
	}

	var entry FormData
	body, _ := io.ReadAll(r.Body)
	err := json.Unmarshal(body, &entry)
	if err != nil {
		http.Error(w, "Données invalides", http.StatusBadRequest)
		return
	}

	err = saveData(entry)
	if err != nil {
		http.Error(w, "Erreur enregistrement", http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "Données enregistrées avec succès.")
}

func main() {
	http.Handle("/", http.FileServer(http.Dir("./static")))
	http.HandleFunc("/submit", formHandler)

	fmt.Println("Serveur démarré sur http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
