package controllers

import (
	temp "Nasa/assets/Templates"
	"log"
	"net/http"
	"strconv"
	"strings"
)

type ResourceDetailData struct {
	Item         UnifiedImageData
	RelatedItems []UnifiedImageData
}

func ResourceDetailPage(w http.ResponseWriter, r *http.Request) {
	// Extraire l'ID de la ressource de l'URL
	resourceID := strings.TrimPrefix(r.URL.Path, "/resource/")

	// Convertir l'ID en entier
	id, err := strconv.Atoi(resourceID)
	if err != nil {
		http.Error(w, "ID de ressource invalide", http.StatusBadRequest)
		return
	}

	// Récupérer toutes les données
	allItems := getAllResources()

	// Vérifier si l'ID est valide
	if id < 0 || id >= len(allItems) {
		http.Error(w, "Ressource non trouvée", http.StatusNotFound)
		return
	}

	// Récupérer la ressource demandée
	resource := allItems[id]

	// Trouver des ressources similaires (même source ou correspondance de mots-clés)
	var relatedItems []UnifiedImageData
	for _, item := range allItems {
		// Éviter de s'inclure soi-même
		if item.URL == resource.URL {
			continue
		}

		// Même source
		if item.Source == resource.Source {
			relatedItems = append(relatedItems, item)
		}

		// Limiter à 4 ressources similaires
		if len(relatedItems) >= 4 {
			break
		}
	}

	data := ResourceDetailData{
		Item:         resource,
		RelatedItems: relatedItems,
	}

	if err := temp.Temp.ExecuteTemplate(w, "resource", data); err != nil {
		log.Printf("Erreur lors du rendu du template: %v", err)
		http.Error(w, "Erreur lors du rendu de la page", http.StatusInternalServerError)
	}
}
