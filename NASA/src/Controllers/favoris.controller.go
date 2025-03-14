package controllers

import (
	temp "Nasa/assets/Templates"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

// Structure pour représenter un favori
type Favorite struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	URL         string `json:"url"`
	Type        string `json:"type"`
	Source      string `json:"source"`
	Date        string `json:"date"`
	Explanation string `json:"explanation"`
}

// Chemin du fichier JSON pour stocker les favoris
const favorisFilePath = "data/favoris.json"

// Initialiser le fichier de favoris s'il n'existe pas
func init() {
	// Créer le répertoire data s'il n'existe pas
	if err := os.MkdirAll(filepath.Dir(favorisFilePath), 0755); err != nil {
		log.Printf("Erreur lors de la création du répertoire data: %v", err)
	}

	// Vérifier si le fichier existe
	if _, err := os.Stat(favorisFilePath); os.IsNotExist(err) {
		// Créer un fichier vide avec un tableau JSON vide
		emptyFavorites := []Favorite{}
		data, err := json.MarshalIndent(emptyFavorites, "", "  ")
		if err != nil {
			log.Printf("Erreur lors de la création du fichier de favoris: %v", err)
			return
		}

		if err := os.WriteFile(favorisFilePath, data, 0644); err != nil {
			log.Printf("Erreur lors de l'écriture du fichier de favoris: %v", err)
		}
	}
}

// Récupérer tous les favoris depuis le fichier JSON
func getFavorites() ([]Favorite, error) {
	// Lire le fichier
	data, err := os.ReadFile(favorisFilePath)
	if err != nil {
		return nil, fmt.Errorf("erreur lors de la lecture du fichier de favoris: %v", err)
	}

	// Désérialiser les données
	var favorites []Favorite
	if err := json.Unmarshal(data, &favorites); err != nil {
		return nil, fmt.Errorf("erreur lors de la désérialisation des favoris: %v", err)
	}

	return favorites, nil
}

// Sauvegarder les favoris dans le fichier JSON
func saveFavorites(favorites []Favorite) error {
	// Sérialiser les données avec une indentation pour la lisibilité
	data, err := json.MarshalIndent(favorites, "", "  ")
	if err != nil {
		return fmt.Errorf("erreur lors de la sérialisation des favoris: %v", err)
	}

	// Écrire dans le fichier
	if err := os.WriteFile(favorisFilePath, data, 0644); err != nil {
		return fmt.Errorf("erreur lors de l'écriture du fichier de favoris: %v", err)
	}

	return nil
}

// Contrôleur pour afficher la page des favoris
func FavorisPage(w http.ResponseWriter, r *http.Request) {
	log.Println("Affichage de la page des favoris")

	// Récupérer tous les favoris
	favorites, err := getFavorites()
	if err != nil {
		log.Printf("Erreur lors de la récupération des favoris: %v", err)
		http.Error(w, "Erreur lors de la récupération des favoris", http.StatusInternalServerError)
		return
	}

	// Convertir les favoris en UnifiedImageData pour réutiliser les templates
	items := make([]UnifiedImageData, len(favorites))
	for i, fav := range favorites {
		items[i] = UnifiedImageData{
			ID:          fav.ID,
			Title:       fav.Title,
			URL:         fav.URL,
			Type:        fav.Type,
			Source:      fav.Source,
			Date:        fav.Date,
			Explanation: fav.Explanation,
		}
	}

	// Préparer les données pour le template
	data := struct {
		Items []UnifiedImageData
	}{
		Items: items,
	}

	// Exécuter le template
	if err := temp.Temp.ExecuteTemplate(w, "favoris", data); err != nil {
		log.Printf("Erreur lors du rendu du template: %v", err)
		http.Error(w, "Erreur lors du rendu de la page", http.StatusInternalServerError)
	}
}

// Contrôleur pour ajouter un élément aux favoris
func AddToFavoris(w http.ResponseWriter, r *http.Request) {
	log.Println("Ajout d'un élément aux favoris")

	// Vérifier la méthode HTTP
	if r.Method != http.MethodPost {
		http.Error(w, "Méthode non autorisée", http.StatusMethodNotAllowed)
		return
	}

	// Récupérer l'ID de la ressource depuis l'URL
	resourceID := strings.TrimPrefix(r.URL.Path, "/favoris/add/")

	log.Printf("Ajout aux favoris de la ressource avec ID: %s", resourceID)

	// Récupérer toutes les données EXACTEMENT comme dans ResourceDetailPage
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

	// Récupérer les données de la galerie NASA
	gallery, err := getNASAGallery()
	if err == nil {
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

	// Ajouter les données Earth
	earthData := getEarthImagesForCollection()
	allItems = append(allItems, earthData...)

	// Recherche de la ressource par ID
	var foundResource UnifiedImageData
	var resourceFound bool = false

	// 1. D'abord, chercher par ID exact
	for _, item := range allItems {
		if item.ID == resourceID {
			foundResource = item
			resourceFound = true
			break
		}
	}

	// 2. Si non trouvé et si l'ID ressemble à un nombre, essayer par index
	if !resourceFound {
		var index int
		if _, err := fmt.Sscanf(resourceID, "%d", &index); err == nil && index >= 0 && index < len(allItems) {
			foundResource = allItems[index]
			resourceFound = true
		}
	}

	if !resourceFound {
		log.Printf("Ressource non trouvée avec ID: %s", resourceID)
		http.Error(w, "Ressource non trouvée", http.StatusNotFound)
		return
	}

	// Récupérer les favoris existants
	favorites, err := getFavorites()
	if err != nil {
		log.Printf("Erreur lors de la récupération des favoris: %v", err)
		http.Error(w, "Erreur interne du serveur", http.StatusInternalServerError)
		return
	}

	// Vérifier si la ressource est déjà dans les favoris
	for _, fav := range favorites {
		if fav.ID == foundResource.ID {
			log.Printf("La ressource est déjà dans les favoris: %s", foundResource.ID)
			// Rediriger vers la page de détail
			http.Redirect(w, r, "/favoris?exists=true", http.StatusSeeOther)
			return
		}
	}

	// Créer un nouveau favori
	newFavorite := Favorite{
		ID:          foundResource.ID,
		Title:       foundResource.Title,
		URL:         foundResource.URL,
		Type:        foundResource.Type,
		Source:      foundResource.Source,
		Date:        foundResource.Date,
		Explanation: foundResource.Explanation,
	}

	// Ajouter le favori à la liste
	favorites = append(favorites, newFavorite)

	// Enregistrer la liste mise à jour
	if err := saveFavorites(favorites); err != nil {
		log.Printf("Erreur lors de l'enregistrement des favoris: %v", err)
		http.Error(w, "Erreur interne du serveur", http.StatusInternalServerError)
		return
	}

	log.Printf("Ressource ajoutée aux favoris avec succès: %s", foundResource.Title)

	// Rediriger vers la page de détail avec un message de succès
	http.Redirect(w, r, "/favoris?added="+resourceID, http.StatusSeeOther)
}

// Contrôleur pour supprimer un élément des favoris
func RemoveFromFavoris(w http.ResponseWriter, r *http.Request) {
	log.Println("Suppression d'un élément des favoris")

	// Vérifier la méthode HTTP
	if r.Method != http.MethodPost {
		http.Error(w, "Méthode non autorisée", http.StatusMethodNotAllowed)
		return
	}

	// Récupérer l'ID de la ressource depuis l'URL
	resourceID := strings.TrimPrefix(r.URL.Path, "/favoris/remove/")

	log.Printf("Suppression des favoris de la ressource avec ID: %s", resourceID)

	// Récupérer les favoris existants
	favorites, err := getFavorites()
	if err != nil {
		log.Printf("Erreur lors de la récupération des favoris: %v", err)
		http.Error(w, "Erreur interne du serveur", http.StatusInternalServerError)
		return
	}

	// Chercher l'index du favori à supprimer
	index := -1
	for i, fav := range favorites {
		if fav.ID == resourceID {
			index = i
			break
		}
	}

	// Si le favori n'a pas été trouvé
	if index == -1 {
		log.Printf("Favori non trouvé avec ID: %s", resourceID)
		http.Error(w, "Favori non trouvé", http.StatusNotFound)
		return
	}

	// Supprimer le favori de la liste
	favorites = append(favorites[:index], favorites[index+1:]...)

	// Enregistrer la liste mise à jour
	if err := saveFavorites(favorites); err != nil {
		log.Printf("Erreur lors de l'enregistrement des favoris: %v", err)
		http.Error(w, "Erreur interne du serveur", http.StatusInternalServerError)
		return
	}

	log.Printf("Ressource supprimée des favoris avec succès: %s", resourceID)

	// Rediriger vers la page des favoris
	http.Redirect(w, r, "/favoris", http.StatusSeeOther)
}
