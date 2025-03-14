package controllers

import (
	temp "Nasa/assets/Templates"
	"fmt"
	"log"
	"net/http"
	"strings"
)

// Fonction pour récupérer les catégories avec des images personnalisées
func getCategoriesWithCustomImages() []CategoryInfo {
	return []CategoryInfo{
		{
			Name:     "APOD",
			ImageURL: "/assets/Images/custom/galaxies.jpg",
		},
		{
			Name:     "Rover",
			ImageURL: "/assets/Images/custom/rovers.jpg",
		},
		{
			Name:     "Earth",
			ImageURL: "/assets/Images/custom/nebulae.jpg",
		},
		{
			Name:     "Gallery",
			ImageURL: "/assets/Images/custom/moon.jpg",
		},
	}
}

// Fonction pour récupérer toutes les noms de catégories
func getAllCategories() []string {
	categories := getCategoriesWithCustomImages()
	var names []string
	for _, cat := range categories {
		names = append(names, cat.Name)
	}
	return names
}

// Fonction pour récupérer les ressources selon la catégorie
func getResourcesByCategory(category string) []interface{} {
	switch category {
	case "APOD":
		// Récupère les 4 premières ressources APOD
		resources := getAPODData()
		if len(resources) > 12 {
			return interfaceSlice(resources[:12])
		}
		return interfaceSlice(resources)
	case "Rover":
		// Récupère les 4 premières ressources Rover
		resources := getAllRoversData()
		if len(resources) > 12 {
			return interfaceSlice(resources[:12])
		}
		return interfaceSlice(resources)
	case "Earth":
		// Récupère les 4 premières ressources Earth
		resources := getEarthImagesForCollection()
		if len(resources) > 12 {
			return interfaceSlice(resources[:12])
		}
		return interfaceSlice(resources)
	case "Gallery":
		// Récupère les 4 premières ressources Gallery
		resources, err := getNASAGallery()
		if err != nil {
			log.Printf("Erreur lors de la récupération de la galerie: %v", err)
			return nil
		}

		// Ajout d'ID pour chaque élément de la galerie
		for i := range resources {
			resources[i].ID = fmt.Sprintf("gallery-%d", i)
		}

		if len(resources) > 12 {
			return interfaceSliceGallery(resources[:12])
		}
		return interfaceSliceGallery(resources)
	default:
		return nil
	}
}

// Fonction utilitaire pour convertir une slice typée en slice d'interface{}
func interfaceSlice(slice interface{}) []interface{} {
	s := make([]interface{}, 0)

	switch v := slice.(type) {
	case []UnifiedImageData:
		for _, item := range v {
			s = append(s, item)
		}
	case []GalleryItem:
		for _, item := range v {
			s = append(s, item)
		}
	}

	return s
}

// Fonction spécifique pour convertir des éléments de galerie en interface{}
func interfaceSliceGallery(items []GalleryItem) []interface{} {
	result := make([]interface{}, len(items))
	for i, item := range items {
		// Conversion vers UnifiedImageData pour assurer la compatibilité avec le template
		result[i] = UnifiedImageData{
			Title:       item.Title,
			URL:         item.ImageURL, // Mapper ImageURL à URL
			Date:        item.DateCreated,
			Explanation: item.Description, // Mapper Description à Explanation
			Source:      "Gallery",
			Type:        "image",
			ID:          item.ID, // Utiliser l'ID déjà défini
		}
	}
	return result
}

// Fonction CategoryPage mise à jour
func CategoryPage(w http.ResponseWriter, r *http.Request) {
	log.Printf("URL de la requête: %s", r.URL.Path)

	// Extraire le nom de la catégorie de l'URL
	requestedCategory := strings.TrimPrefix(r.URL.Path, "/categories/")

	// Si aucune catégorie n'est spécifiée ou si l'URL est exactement "/categories", afficher la liste
	if requestedCategory == "" || r.URL.Path == "/categories" {
		categories := getCategoriesWithCustomImages()

		// Données pour la page des catégories
		previewData := struct {
			Categories []CategoryInfo
		}{
			Categories: categories,
		}

		log.Printf("Affichage de toutes les catégories avec images personnalisées")

		if err := temp.Temp.ExecuteTemplate(w, "categories", previewData); err != nil {
			log.Printf("Erreur lors du rendu du template 'categories': %v", err)
			http.Error(w, "Erreur lors du rendu de la page", http.StatusInternalServerError)
		}
		return
	}

	// Récupérer la liste de toutes les catégories pour le menu
	allCategories := getAllCategories()

	// Récupérer les ressources pour la catégorie demandée
	resources := getResourcesByCategory(requestedCategory)

	// Logging pour déboguer
	log.Printf("Catégorie demandée: %s, nombre de ressources: %d", requestedCategory, len(resources))

	// Préparer les données pour le template
	data := CategoryData{
		CategoryName: requestedCategory,
		Categories:   allCategories,
		Resources:    resources,
	}

	if err := temp.Temp.ExecuteTemplate(w, "category", data); err != nil {
		log.Printf("Erreur lors du rendu du template 'category': %v", err)
		http.Error(w, "Erreur lors du rendu de la page", http.StatusInternalServerError)
	}
}
