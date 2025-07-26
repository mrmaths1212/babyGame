package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

// Structure pour les données utilisateur
type FormData struct {
	VotrePrenom string  `json:"votrePrenom"`
	VotreNom    string  `json:"votreNom"`
	Prenom      string  `json:"prenom"`
	Poids       float64 `json:"poids"`
	Taille      float64 `json:"taille"`
	Naissance   string  `json:"naissance"`
}

var (
	dataFile = "data/results.json"
	dataLock sync.Mutex
)

// GET / → affiche la page HTML
func indexHandler(w http.ResponseWriter, r *http.Request) {
	tmplPath := filepath.Join("templates", "index.html")
	tmpl, err := template.ParseFiles(tmplPath)
	if err != nil {
		http.Error(w, "Erreur de template", 500)
		return
	}
	tmpl.Execute(w, nil)
}

// POST /submit → enregistre les données
func formHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Méthode non autorisée", http.StatusMethodNotAllowed)
		return
	}

	var entry FormData
	err := json.NewDecoder(r.Body).Decode(&entry)
	if err != nil {
		http.Error(w, "Données invalides", http.StatusBadRequest)
		return
	}

	// Sauvegarder dans le fichier
	err = saveData(entry)
	if err != nil {
		http.Error(w, "Erreur d'enregistrement", http.StatusInternalServerError)
		return
	}
	log.Printf("Données reçues: %+v\n", entry)

	w.Write([]byte("Données enregistrées."))
}

// Enregistre les données dans data/results.json
func saveData(entry FormData) error {
	dataLock.Lock()
	defer dataLock.Unlock()

	// Charger les anciennes données
	var all map[string]FormData = make(map[string]FormData)

	file, err := os.ReadFile(dataFile)
	if err == nil && len(file) > 0 {
		_ = json.Unmarshal(file, &all)
	}
	key := fmt.Sprintf("%s.%s", entry.VotrePrenom, entry.VotreNom)
	key = strings.ToLower(strings.ReplaceAll(key, " ", ""))

	// Mettre à jour ou ajouter
	all[key] = entry

	// Réécrire le fichier
	output, err := json.MarshalIndent(all, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(dataFile, output, 0644)
}

func main() {
	// Router
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/submit", formHandler)

	// Fichiers statiques (JS/CSS)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	log.Println("Serveur démarré sur http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
