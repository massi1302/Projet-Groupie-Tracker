package controllers

import (
	temp "Nasa/assets/Templates"
	"fmt"
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

	log.Printf("Recherche de la ressource avec ID: %s", resourceID)

	// Récupérer toutes les données EXACTEMENT comme dans CollectionPage
	allItems := make([]UnifiedImageData, 0)

	// Récupérer les données APOD
	apodData := getAPODData()
	for i, item := range apodData {
		// Assigner un ID unique basé sur l'index si non présent
		if item.ID == "" {
			item.ID = fmt.Sprintf("apod-%d", i)
		}

		if item.Type == "video" {
			// Convertir les URLs YouTube en format embed
			if strings.Contains(item.URL, "youtube.com/watch?v=") {
				item.URL = strings.Replace(item.URL, "watch?v=", "embed/", 1)
			}
			allItems = append(allItems, item)
		} else if item.Type == "image" && isValidImage(item.URL) {
			allItems = append(allItems, item)
		}
	}
	log.Printf("Nombre d'éléments APOD chargés: %d", len(allItems))

	// Récupérer les données de la galerie NASA
	gallery, err := getNASAGallery()
	if err != nil {
		log.Printf("Erreur lors de la récupération de la galerie NASA: %v", err)
	} else {
		// Convertir les éléments de la galerie au format UnifiedImageData
		for i, item := range gallery {
			if isValidImage(item.ImageURL) {
				galleryItem := UnifiedImageData{
					Title:       item.Title,
					URL:         item.ImageURL,
					Date:        item.DateCreated,
					Explanation: item.Description,
					Source:      "NASA Gallery",
					Type:        "image",
					ID:          fmt.Sprintf("gallery-%d", i),
				}
				allItems = append(allItems, galleryItem)
			}
		}
	}
	log.Printf("Nombre d'éléments après ajout de la galerie: %d", len(allItems))

	// Ajouter les données Earth
	earthData := getEarthImagesForCollection()
	allItems = append(allItems, earthData...)
	log.Printf("Nombre total d'éléments chargés: %d", len(allItems))

	// Recherche de la ressource par ID
	var resource UnifiedImageData
	var foundResource bool = false

	// Afficher les IDs disponibles pour le débogage
	log.Printf("Recherche de l'ID: %s parmi les ressources disponibles:", resourceID)
	for i, item := range allItems {
		if i < 5 || item.ID == resourceID { // Afficher les 5 premiers et celui recherché
			log.Printf("  Item %d: ID=%s, Titre=%s", i, item.ID, item.Title)
		}
	}

	// 1. D'abord, chercher par ID exact
	for _, item := range allItems {
		if item.ID == resourceID {
			resource = item
			foundResource = true
			log.Printf("Ressource trouvée par ID exact: %s", resourceID)
			break
		}
	}

	// 2. Si non trouvé et si l'ID ressemble à un nombre, essayer par index
	if !foundResource {
		if id, err := strconv.Atoi(resourceID); err == nil && id >= 0 && id < len(allItems) {
			resource = allItems[id]
			foundResource = true
			log.Printf("Ressource trouvée par index numérique: %d", id)
		}
	}

	if !foundResource {
		log.Printf("ERREUR: Ressource non trouvée avec ID: %s", resourceID)
		http.Error(w, "Ressource non trouvée", http.StatusNotFound)
		return
	}

	log.Printf("Ressource trouvée: %s (Source: %s, ID: %s)", resource.Title, resource.Source, resource.ID)

	// Trouver des ressources similaires (même source)
	var relatedItems []UnifiedImageData
	for _, item := range allItems {
		// Éviter de s'inclure soi-même
		if item.ID == resource.ID {
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
