package controllers

import (
	temp "Nasa/assets/Templates"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
)

func CollectionPage(w http.ResponseWriter, r *http.Request) {
	log.Println("Début du traitement de la requête CollectionPage")

	// Récupérer le paramètre de page depuis l'URL
	pageStr := r.URL.Query().Get("page")
	page := 1
	if pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
			page = p
		}
	}

	// Nombre d'éléments par page
	itemsPerPage := 12

	// Récupérer toutes les données
	allItems := make([]UnifiedImageData, 0)

	// Récupérer les données APOD
	apodData := getAPODData()
	for _, item := range apodData {
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
	log.Printf("Nombre total d'éléments APOD filtrés et valides: %d", len(allItems))

	// Récupérer les données de la galerie NASA
	gallery, err := getNASAGallery()
	if err != nil {
		log.Printf("Erreur lors de la récupération de la galerie NASA: %v", err)
	} else {
		// Convertir les éléments de la galerie au format UnifiedImageData
		for i, item := range gallery {
			if isValidImage(item.ImageURL) {
				allItems = append(allItems, UnifiedImageData{
					Title:       item.Title,
					URL:         item.ImageURL,
					Date:        item.DateCreated,
					Explanation: item.Description,
					Source:      "NASA Gallery",
					Type:        "image",
					ID:          fmt.Sprintf("gallery-%d", i), // Ajout d'un ID unique
				})
			}
		}
	}
	log.Printf("Nombre total d'éléments après ajout de la galerie: %d", len(allItems))

	// Ajouter les données Earth (intégration de la page Earth)
	earthData := getEarthImagesForCollection()
	allItems = append(allItems, earthData...)
	log.Printf("Nombre total d'éléments après ajout des images Earth: %d", len(allItems))

	// Calculer la pagination
	totalItems := len(allItems)
	totalPages := (totalItems + itemsPerPage - 1) / itemsPerPage

	// Calculer les pages précédente et suivante
	previousPage := page - 1
	if previousPage < 1 {
		previousPage = 1
	}

	nextPage := page + 1
	if nextPage > totalPages {
		nextPage = totalPages
	}

	// Ajuster la page si nécessaire
	if page > totalPages && totalPages > 0 {
		page = totalPages
	}

	// Calculer les indices de début et de fin pour la page actuelle
	startIndex := (page - 1) * itemsPerPage
	endIndex := startIndex + itemsPerPage
	if endIndex > totalItems {
		endIndex = totalItems
	}

	// Sélectionner les éléments pour la page actuelle
	var pageItems []UnifiedImageData
	if startIndex < totalItems {
		pageItems = allItems[startIndex:endIndex]
	}

	// Créer une plage pour les liens de pagination (montrer 5 pages autour de la page actuelle)
	paginationRange := make([]int, 0)
	startPage := max(1, page-2)
	endPage := min(totalPages, page+2)

	for i := startPage; i <= endPage; i++ {
		paginationRange = append(paginationRange, i)
	}

	data := struct {
		Items           []UnifiedImageData
		CurrentPage     int
		TotalPages      int
		PaginationRange []int
		PreviousPage    int
		NextPage        int
		ItemsPerPage    int
	}{
		Items:           pageItems,
		CurrentPage:     page,
		TotalPages:      totalPages,
		PaginationRange: paginationRange,
		PreviousPage:    previousPage,
		NextPage:        nextPage,
		ItemsPerPage:    itemsPerPage,
	}

	if err := temp.Temp.ExecuteTemplate(w, "collection", data); err != nil {
		log.Printf("Erreur lors du rendu du template: %v", err)
		http.Error(w, "Erreur lors du rendu de la page", http.StatusInternalServerError)
		return
	}
}
