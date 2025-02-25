package controllers

import (
	temp "Nasa/assets/Templates"
	"log"
	"net/http"
	"strings"
)

// Modifions la fonction CategoryPage pour envoyer des objets CategoryInfo au template
func CategoryPage(w http.ResponseWriter, r *http.Request) {
	log.Printf("URL de la requête: %s", r.URL.Path)

	// Extraire le nom de la catégorie de l'URL
	requestedCategory := strings.TrimPrefix(r.URL.Path, "/categories/")

	// Si aucune catégorie n'est spécifiée ou si l'URL est exactement "/categories", afficher la liste
	if requestedCategory == "" || r.URL.Path == "/categories" {
		categories := getAllCategoriesWithImages()

		// Récupérer quelques éléments de chaque catégorie pour l'aperçu
		previewData := struct {
			Categories []CategoryInfo
		}{
			Categories: categories,
		}

		log.Printf("Affichage de toutes les catégories avec images")

		if err := temp.Temp.ExecuteTemplate(w, "categories", previewData); err != nil {
			log.Printf("Erreur lors du rendu du template 'categories': %v", err)
			http.Error(w, "Erreur lors du rendu de la page", http.StatusInternalServerError)
		}
		return
	}

	// Récupérer les éléments de la catégorie demandée
	items := getItemsByCategory(requestedCategory)

	log.Printf("Catégorie demandée: %s, Nombre d'éléments trouvés: %d", requestedCategory, len(items))

	if len(items) == 0 {
		// Au lieu de renvoyer une erreur 404, afficher une page avec un message
		allCategories := getAllCategories()
		data := CategoryData{
			CategoryName: requestedCategory,
			Items:        []UnifiedImageData{},
			Categories:   allCategories,
		}

		if err := temp.Temp.ExecuteTemplate(w, "category", data); err != nil {
			log.Printf("Erreur lors du rendu du template 'category': %v", err)
			http.Error(w, "Erreur lors du rendu de la page", http.StatusInternalServerError)
		}
		return
	}

	// Récupérer la liste de toutes les catégories pour le menu
	allCategories := getAllCategories()

	data := CategoryData{
		CategoryName: requestedCategory,
		Items:        items,
		Categories:   allCategories,
	}

	if err := temp.Temp.ExecuteTemplate(w, "category", data); err != nil {
		log.Printf("Erreur lors du rendu du template 'category': %v", err)
		http.Error(w, "Erreur lors du rendu de la page", http.StatusInternalServerError)
	}
}

// Récupère la liste des catégories disponibles
func getAllCategories() []string {
	return []string{
		"Mars",
		"Moon",
		"Space Station",
		"Webb Telescope",
		"Solar System",
		"Galaxies",
		"Exoplanets",
		"Nebulae",
		"APOD",
		"Rovers",
	}
}

// Filtre les éléments par catégorie
func getItemsByCategory(category string) []UnifiedImageData {
	// Récupérer toutes les données
	allItems := getAllResourcess()
	var filteredItems []UnifiedImageData

	// Cas spécial pour APOD
	if category == "APOD" {
		for _, item := range allItems {
			if item.Source == "APOD" {
				filteredItems = append(filteredItems, item)
			}
		}
		return filteredItems
	}

	// Cas spécial pour Rovers
	if category == "Rovers" {
		for _, item := range allItems {
			if strings.Contains(item.Source, "Mars Rover") {
				filteredItems = append(filteredItems, item)
			}
		}
		return filteredItems
	}

	// Pour les autres catégories, chercher dans les titres et descriptions
	for _, item := range allItems {
		titleMatch := strings.Contains(strings.ToLower(item.Title), strings.ToLower(category))
		descMatch := strings.Contains(strings.ToLower(item.Explanation), strings.ToLower(category))

		if titleMatch || descMatch {
			filteredItems = append(filteredItems, item)
		}
	}

	return filteredItems
}

// Fonction getAllResources qui manquait
func getAllResourcess() []UnifiedImageData {
	// Récupérer toutes les données comme dans getAllResources() original
	apodData := getAPODData()
	marsData := getAllRoversData()

	// Combiner les données
	allItems := append(apodData, marsData...)
	return allItems
}

func getAllCategoriesWithImages() []CategoryInfo {
	// Récupérer toutes les données
	allItems := getAllResourcess()
	categories := getAllCategories()

	var categoriesWithImages []CategoryInfo

	// Pour chaque catégorie, trouvons une image représentative
	for _, category := range categories {
		var imageURL string

		// Collecter plusieurs candidats d'images pour cette catégorie
		var candidateImages []string
		for _, item := range allItems {
			if item.Type != "image" {
				continue // Ignorer les vidéos et autres médias non-image
			}

			// Déterminer si l'élément correspond à la catégorie
			isMatch := false
			if category == "APOD" && item.Source == "APOD" {
				isMatch = true
			} else if category == "Rovers" && strings.Contains(item.Source, "Mars Rover") {
				isMatch = true
			} else if strings.Contains(strings.ToLower(item.Title), strings.ToLower(category)) ||
				strings.Contains(strings.ToLower(item.Explanation), strings.ToLower(category)) {
				isMatch = true
			}

			if isMatch {
				candidateImages = append(candidateImages, item.URL)
				if len(candidateImages) >= 5 {
					break // Limiter à 5 candidats par catégorie
				}
			}
		}

		// Tester les candidats jusqu'à trouver une image valide
		for _, imgURL := range candidateImages {
			if isValidImage(imgURL) {
				imageURL = imgURL
				log.Printf("Image valide trouvée pour %s: %s", category, imageURL)
				break
			} else {
				log.Printf("Image candidate invalide pour %s: %s", category, imgURL)
			}
		}

		// Si aucune image valide n'a été trouvée, utiliser l'image par défaut
		if imageURL == "" {
			imageURL = getDefaultImageForCategory(category)
			log.Printf("Utilisation de l'image par défaut pour %s: %s", category, imageURL)

			// Vérifions même l'image par défaut
			if !isValidImage(imageURL) {
				log.Printf("AVERTISSEMENT: Même l'image par défaut pour %s est invalide!", category)
				// Utiliser une image de secours absolue
				imageURL = "https://apod.nasa.gov/apod/image/2202/NGC7380_Soulie_1080.jpg"
			}
		}

		categoriesWithImages = append(categoriesWithImages, CategoryInfo{
			Name:     category,
			ImageURL: imageURL,
		})
	}

	return categoriesWithImages
}

// Fonction pour fournir une image par défaut selon la catégorie
func getDefaultImageForCategory(category string) string {
	defaultImages := map[string]string{
		"Mars":           "https://apod.nasa.gov/apod/image/2010/PIA24377mars1024c.jpg",
		"Moon":           "https://apod.nasa.gov/apod/image/2010/MoonIllusionGraphic.jpg",
		"Space Station":  "https://apod.nasa.gov/apod/image/2009/ISS_Nespoli_960.jpg",
		"Webb Telescope": "https://apod.nasa.gov/apod/image/2307/WR124_Webb_960.jpg",
		"Solar System":   "https://apod.nasa.gov/apod/image/1912/SolarSystem_Positions_1080.jpg",
		"Galaxies":       "https://apod.nasa.gov/apod/image/2105/M87_EHT_1080.jpg",
		"Exoplanets":     "https://apod.nasa.gov/apod/image/2307/ExoplanetOrbits_KeelerWebb_960.jpg",
		"Nebulae":        "https://apod.nasa.gov/apod/image/2306/M17Web_1024.jpg",
		"APOD":           "https://apod.nasa.gov/apod/image/2302/cometE3_Bartlett_960.jpg",
		"Rovers":         "https://mars.nasa.gov/system/resources/detail_files/25962_PIA25026-web.jpg",
	}

	if url, exists := defaultImages[category]; exists {
		return url
	}

	// Image par défaut si aucune correspondance
	return "https://apod.nasa.gov/apod/image/2202/NGC7380_Soulie_1080.jpg"
}
