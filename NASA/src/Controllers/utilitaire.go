package controllers

import (
	"encoding/json"
	"fmt"
	"html"
	"log"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"sync"
	"time"
)

func getAPODData() []UnifiedImageData {
	apiKey := "UcitNGyctCYzAsLCq0lO2q9URTcol74cw31ThtML"
	endDate := time.Now()
	startDate := endDate.AddDate(0, 0, -30)

	url := fmt.Sprintf("https://api.nasa.gov/planetary/apod?api_key=%s&start_date=%s&end_date=%s",
		apiKey,
		startDate.Format("2006-01-02"),
		endDate.Format("2006-01-02"))

	var apodResponses []APODResponse
	if err := fetchAPI(url, &apodResponses); err != nil {
		log.Printf("Erreur APOD API: %v", err)
		return nil
	}

	var unifiedData []UnifiedImageData
	for _, apod := range apodResponses {
		mediaType := "image"
		if strings.Contains(apod.URL, "youtube.com") {
			mediaType = "video"
			// Convertir directement l'URL YouTube en format embed
			apod.URL = strings.Replace(apod.URL, "watch?v=", "embed/", 1)
		}

		unifiedData = append(unifiedData, UnifiedImageData{
			Title:       apod.Title,
			URL:         apod.URL,
			Date:        apod.Date,
			Explanation: apod.Explanation,
			Source:      "APOD",
			Type:        mediaType,
		})
	}

	return unifiedData
}

func getAllRoversData() []UnifiedImageData {
	rovers := []RoverConfig{
		{Name: "curiosity", Sols: []int{1000, 2000, 3000}},
		{Name: "perseverance", Sols: []int{100, 200, 300}},
	}

	var allData []UnifiedImageData
	var wg sync.WaitGroup
	var mutex sync.Mutex

	// Utiliser une map pour suivre les URLs déjà ajoutées et éviter les doublons
	urlMap := make(map[string]bool)

	for _, rover := range rovers {
		for _, sol := range rover.Sols {
			wg.Add(1)
			go func(roverName string, sol int) {
				defer wg.Done()
				data := getRoverData(roverName, sol)

				mutex.Lock()
				// Ajouter uniquement les images avec des URLs uniques
				for _, item := range data {
					if _, exists := urlMap[item.URL]; !exists {
						urlMap[item.URL] = true
						allData = append(allData, item)
					}
				}
				mutex.Unlock()
			}(rover.Name, sol)
		}
	}

	wg.Wait()
	return allData
}

func getRoverData(roverName string, sol int) []UnifiedImageData {
	apiKey := "UcitNGyctCYzAsLCq0lO2q9URTcol74cw31ThtML"
	url := fmt.Sprintf("https://api.nasa.gov/mars-photos/api/v1/rovers/%s/photos?sol=%d&page=1&api_key=%s",
		roverName, sol, apiKey)

	var roverResponse MarsRoverResponse
	if err := fetchAPI(url, &roverResponse); err != nil {
		log.Printf("Erreur Mars Rover API pour %s sol %d: %v", roverName, sol, err)
		return nil
	}

	var unifiedData []UnifiedImageData
	maxPhotos := 15 // Photos par sol et par rover
	processedPhotos := 0

	// Utiliser une map locale pour éviter les doublons dans cette requête spécifique
	localUrlMap := make(map[string]bool)

	for _, photo := range roverResponse.Photos {
		if processedPhotos >= maxPhotos {
			break
		}

		imgURL := strings.Replace(photo.ImgSrc, "http://", "https://", 1)

		// Vérifier que l'URL n'est pas déjà dans cette requête
		if _, exists := localUrlMap[imgURL]; exists {
			continue
		}

		if isValidImage(imgURL) {
			localUrlMap[imgURL] = true
			unifiedData = append(unifiedData, UnifiedImageData{
				Title: fmt.Sprintf("%s (Sol %d): %s", photo.Rover.Name, sol, photo.Camera.FullName),
				URL:   imgURL,
				Date:  photo.EarthDate,
				Explanation: fmt.Sprintf("Photo prise sur Mars par le rover %s lors du sol %d avec la caméra %s",
					photo.Rover.Name, sol, photo.Camera.FullName),
				Source: fmt.Sprintf("Mars Rover - %s", photo.Rover.Name),
				Type:   "image",
				ID:     fmt.Sprintf("%s-%d-%s", photo.Rover.Name, sol, strings.Replace(photo.ImgSrc, "/", "-", -1)),
			})
			processedPhotos++
		}
	}

	return unifiedData
}

func fetchAPI(url string, target interface{}) error {
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("API returned status code %d", resp.StatusCode)
	}

	return json.NewDecoder(resp.Body).Decode(target)
}

// Cette fonction combine les données des différentes APIs
func getAllResources() []UnifiedImageData {
	// Récupérer les données APOD
	apodData := getAPODData()

	// Récupérer les données des rovers
	marsData := getAllRoversData()

	// Combiner les deux ensembles de données
	allItems := append(apodData, marsData...)

	return allItems
}

func getNASAGallery() ([]GalleryItem, error) {
	url := "https://images-api.nasa.gov/search?media_type=image&year_start=2024&keywords=space"
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var response ImageLibraryResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, err
	}

	var gallery []GalleryItem
	for _, item := range response.Collection.Items {
		// Limiter à 6 images
		if len(gallery) >= 24 {
			break
		}

		if len(item.Data) > 0 && len(item.Links) > 0 {
			dateCreated, err := time.Parse("2006-01-02T15:04:05Z", item.Data[0].DateCreated)
			if err != nil {
				log.Printf("Erreur parsing date: %v", err)
				continue
			}

			gallery = append(gallery, GalleryItem{
				Title:       item.Data[0].Title,
				Description: item.Data[0].Description,
				ImageURL:    item.Links[0].Href,
				DateCreated: dateCreated.Format("02/01/2006"),
			})
		}
	}

	return gallery, nil
}

func getNASATopics() ([]Topic, error) {
	// Augmenter la liste des mots-clés
	keywords := []string{"Mars", "Moon", "Space Station", "Webb Telescope", "Solar System", "Artemis", "Black Holes", "Galaxies", "comets", "asteroids", "nebulae", "supernovae", "exoplanets", "cosmic rays", "dark matter", "dark energy", "gravitational waves", "neutron stars", "pulsars", "quasars", "cosmic microwave background", "cosmic inflation", "cosmic strings", "cosmic voids", "cosmic web"}
	var topics []Topic

	for _, keyword := range keywords {
		url := fmt.Sprintf("https://images-api.nasa.gov/search?q=%s&media_type=image&year_start=2023", keyword)
		resp, err := http.Get(url)
		if err != nil {
			continue
		}
		defer resp.Body.Close()

		var response TopicResponse
		if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
			continue
		}

		if len(response.Collection.Items) > 0 {
			item := response.Collection.Items[0]
			if len(item.Data) > 0 && len(item.Links) > 0 {
				// Ajouter une vérification de l'URL de l'image
				imageURL := item.Links[0].Href
				if isValidImage(imageURL) {
					topics = append(topics, Topic{
						Title:       item.Data[0].Title,
						Description: item.Data[0].Description,
						ImageURL:    imageURL,
						Keywords:    item.Data[0].Keywords,
					})
				}
			}
		}
	}

	return topics, nil
}

// Nouvelle fonction pour nettoyer le texte
func cleanText(text string) string {
	// Supprimer les caractères HTML
	text = html.UnescapeString(text)
	// Supprimer les espaces superflus
	text = strings.TrimSpace(text)
	return text
}

// Nouvelle fonction pour nettoyer la description HTML
func cleanDescription(description string) string {
	// Supprimer les balises HTML
	description = regexp.MustCompile("<[^>]*>").ReplaceAllString(description, "")
	// Décoder les entités HTML
	description = html.UnescapeString(description)
	// Supprimer les espaces multiples
	description = regexp.MustCompile(`\s+`).ReplaceAllString(description, " ")
	return strings.TrimSpace(description)
}

func getNASANews() ([]NewsItem, error) {
	// Utiliser une API RSS-to-JSON pour convertir le flux RSS de la NASA
	rssURL := "https://www.nasa.gov/rss/dyn/breaking_news.rss"
	url := fmt.Sprintf("https://api.rss2json.com/v1/api.json?rss_url=%s", url.QueryEscape(rssURL))

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("erreur lors de la récupération du flux RSS: %v", err)
	}
	defer resp.Body.Close()

	// Structure pour parser la réponse JSON
	type RSSResponse struct {
		Status string `json:"status"`
		Feed   struct {
			Title string `json:"title"`
		} `json:"feed"`
		Items []struct {
			Title       string `json:"title"`
			PubDate     string `json:"pubDate"`
			Link        string `json:"link"`
			Guid        string `json:"guid"`
			Author      string `json:"author"`
			Thumbnail   string `json:"thumbnail"`
			Description string `json:"description"`
			Content     string `json:"content"`
		} `json:"items"`
	}

	var rssResponse RSSResponse
	if err := json.NewDecoder(resp.Body).Decode(&rssResponse); err != nil {
		return nil, fmt.Errorf("erreur lors du parsing JSON: %v", err)
	}

	var news []NewsItem
	for _, item := range rssResponse.Items {
		// Chercher l'image dans différents champs possibles
		imageURL := item.Thumbnail

		if imageURL == "" {
			// Chercher dans la description/contenu si pas de thumbnail
			imgRegex := regexp.MustCompile(`<img[^>]+src=["']([^"']+)["']`)
			if matches := imgRegex.FindStringSubmatch(item.Description); len(matches) > 1 {
				imageURL = matches[1]
			} else if matches := imgRegex.FindStringSubmatch(item.Content); len(matches) > 1 {
				imageURL = matches[1]
			}
		}

		// S'assurer que l'URL est en HTTPS
		if strings.HasPrefix(imageURL, "http://") {
			imageURL = "https://" + strings.TrimPrefix(imageURL, "http://")
		}

		// Vérifier si l'image est valide
		if imageURL != "" && isValidImage(imageURL) {
			// Parser et formater la date
			pubDate, err := time.Parse("2006-01-02 15:04:05", item.PubDate)
			if err != nil {
				// Essayer un autre format si le premier échoue
				pubDate, err = time.Parse("Mon, 02 Jan 2006 15:04:05 -0700", item.PubDate)
				if err != nil {
					log.Printf("Erreur parsing date pour %s: %v", item.Title, err)
					continue
				}
			}

			news = append(news, NewsItem{
				Title:       cleanText(item.Title),
				Description: cleanDescription(item.Description),
				ImageURL:    imageURL,
				PublishDate: pubDate.Format("02/01/2006"),
				ReadMoreURL: item.Link,
			})

			log.Printf("News Item ajouté: Title=%s, ImageURL=%s", item.Title, imageURL)

			if len(news) >= 7 {
				break
			}
		}
	}

	return news, nil
}

// Amélioration de la fonction de validation des images
func isValidImage(url string) bool {
	// Vérification rapide des formats d'URL connus pour être problématiques
	if strings.Contains(url, "youtube.com") ||
		strings.Contains(url, "vimeo.com") ||
		strings.Contains(url, ".mp4") ||
		strings.Contains(url, ".mov") ||
		strings.Contains(url, ".webm") {
		return false
	}

	// S'assurer que l'URL commence par https
	if strings.HasPrefix(url, "http://") {
		url = "https://" + strings.TrimPrefix(url, "http://")
	}

	// Vérifier par extension si disponible (méthode rapide)
	if strings.HasSuffix(strings.ToLower(url), ".jpg") ||
		strings.HasSuffix(strings.ToLower(url), ".jpeg") ||
		strings.HasSuffix(strings.ToLower(url), ".png") ||
		strings.HasSuffix(strings.ToLower(url), ".gif") {

		// Vérification supplémentaire par requête HEAD
		client := &http.Client{
			Timeout: 2 * time.Second,
			CheckRedirect: func(req *http.Request, via []*http.Request) error {
				if len(via) >= 3 {
					return fmt.Errorf("trop de redirections")
				}
				return nil
			},
		}

		req, err := http.NewRequest("HEAD", url, nil)
		if err != nil {
			return false
		}

		req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36")

		resp, err := client.Do(req)
		if err != nil {
			return false
		}
		defer resp.Body.Close()

		return resp.StatusCode == http.StatusOK
	}

	// Pour les URLs sans extension, vérification complète avec content-type
	client := &http.Client{
		Timeout: 3 * time.Second,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			if len(via) >= 3 {
				return fmt.Errorf("trop de redirections")
			}
			return nil
		},
	}

	req, err := http.NewRequest("HEAD", url, nil)
	if err != nil {
		return false
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36")
	req.Header.Set("Accept", "image/webp,image/apng,image/*,*/*;q=0.8")

	resp, err := client.Do(req)
	if err != nil {
		return false
	}
	defer resp.Body.Close()

	// Vérifier le code de status et le type de contenu
	return resp.StatusCode == http.StatusOK && strings.HasPrefix(resp.Header.Get("Content-Type"), "image/")
}

// fonction pour récupérer les images Earth pour la collection
func getEarthImagesForCollection() []UnifiedImageData {
	// Predefined interesting locations with coordinates (from EarthPage)
	locations := []struct {
		Name string
		Lat  string
		Lon  string
	}{
		{"Paris", "48.8566", "2.3522"},
		{"New York", "40.7128", "-74.0060"},
		{"Sydney", "-33.8688", "151.2093"},
		{"Pyramides de Gizeh", "29.9792", "31.1342"},
		{"Las Vegas", "36.1699", "-115.1398"},
		{"Mont Everest", "27.9881", "86.9250"},
		{"Grand Canyon", "36.0544", "-112.2401"},
		{"Tokyo", "35.6762", "139.6503"},
		{"Rio de Janeiro", "-22.9068", "-43.1729"},
		{"Sahara Desert", "23.4162", "25.6628"},
		{"Amazon Rainforest", "-3.4653", "-62.2159"},
		{"Mount Fuji", "35.3606", "138.7274"},
		{"Victoria Falls", "-17.9244", "25.8567"},
		{"Himalayas", "28.3949", "84.1240"},
		{"Matterhorn", "45.9763", "7.6587"},
		{"Niagara Falls", "43.0812", "-79.0745"},
		{"Easter Island", "-27.1127", "-109.3497"},
		{"Mount Kilimanjaro", "-3.0674", "37.3556"},
		{"Angel Falls", "5.9679", "-62.5357"},
		{"Mount Vesuvius", "40.8219", "14.4265"},
		{"Taj Mahal", "27.1751", "78.0421"},
		{"Machu Picchu", "-13.1631", "-72.5450"},
		{"Great Wall of China", "40.4319", "116.5704"},
		{"Stonehenge", "51.1789", "-1.8262"},
		{"Petra", "30.3285", "35.4444"},
		{"Mount Rushmore", "43.8791", "-103.4591"},
		{"Eiffel Tower", "48.8584", "2.2945"},
		{"Colosseum", "41.8902", "12.4922"},
		{"Statue of Liberty", "40.6892", "-74.0445"},
		{"Burj Khalifa", "25.276987", "55.296249"},
		{"Golden Gate Bridge", "37.8199", "-122.4783"},
		{"Sydney Opera House", "-33.8572", "151.2152"},
		{"Mount Everest Base Camp", "28.0024", "86.8524"},
		{"Matterhorn", "45.9763", "7.6587"},
	}

	// Create a slice to hold all Earth data
	var unifiedItems []UnifiedImageData

	// Fetch data for each location
	for _, location := range locations {
		earthItems := getReliableEarthData(location.Lat, location.Lon, "")

		// Skip any empty results
		if len(earthItems) > 0 {
			for i, item := range earthItems {
				unifiedItem := UnifiedImageData{
					Title:       location.Name + " - " + item.Title,
					URL:         item.URL,
					Date:        item.Date,
					Explanation: item.Description,
					Type:        item.Type,
					Source:      "Earth View - " + item.Source,
					ID:          fmt.Sprintf("earth-%s-%d", location.Name, i), // Ajout d'un ID unique
				}
				unifiedItems = append(unifiedItems, unifiedItem)
			}
		}
	}

	log.Printf("Total Earth data items converted: %d", len(unifiedItems))
	return unifiedItems
}
