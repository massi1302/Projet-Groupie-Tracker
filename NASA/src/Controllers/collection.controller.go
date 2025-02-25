package controllers

import (
	temp "Nasa/assets/Templates"
	"log"
	"net/http"
	"strings"
)

func CollectionPage(w http.ResponseWriter, r *http.Request) {
	log.Println("Début du traitement de la requête CollectionPage")

	data := CollectionData{
		Items: make([]UnifiedImageData, 0),
	}

	// Récupérer les données APOD
	apodData := getAPODData()
	for _, item := range apodData {
		if item.Type == "video" {
			// Convertir les URLs YouTube en format embed
			if strings.Contains(item.URL, "youtube.com/watch?v=") {
				item.URL = strings.Replace(item.URL, "watch?v=", "embed/", 1)
			}
			data.Items = append(data.Items, item)
		} else if item.Type == "image" && isValidImage(item.URL) {
			data.Items = append(data.Items, item)
		}
	}

	// Récupérer les données de tous les rovers
	marsData := getAllRoversData()
	for _, item := range marsData {
		if isValidImage(item.URL) {
			data.Items = append(data.Items, item)
		}
	}

	log.Printf("Nombre total d'éléments filtrés et valides: %d", len(data.Items))

	if err := temp.Temp.ExecuteTemplate(w, "collection", data); err != nil {
		log.Printf("Erreur lors du rendu du template: %v", err)
		http.Error(w, "Erreur lors du rendu de la page", http.StatusInternalServerError)
		return
	}
}
