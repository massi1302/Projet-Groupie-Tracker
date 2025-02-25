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
		{Name: "opportunity", Sols: []int{100, 500, 1000}},
		{Name: "spirit", Sols: []int{100, 500, 1000}},
		{Name: "perseverance", Sols: []int{100, 200, 300}},
	}

	var allData []UnifiedImageData
	var wg sync.WaitGroup
	var mutex sync.Mutex

	for _, rover := range rovers {
		for _, sol := range rover.Sols {
			wg.Add(1)
			go func(roverName string, sol int) {
				defer wg.Done()
				data := getRoverData(roverName, sol)
				mutex.Lock()
				allData = append(allData, data...)
				mutex.Unlock()
			}(rover.Name, sol)
		}
	}

	wg.Wait()
	return allData
}

func getRoverData(roverName string, sol int) []UnifiedImageData {
	apiKey := "UcitNGyctCYzAsLCq0lO2q9URTcol74cw31ThtML"
	url := fmt.Sprintf("https://api.nasa.gov/mars-photos/api/v1/rovers/%s/photos?sol=%d&api_key=%s",
		roverName, sol, apiKey)

	var roverResponse MarsRoverResponse
	if err := fetchAPI(url, &roverResponse); err != nil {
		log.Printf("Erreur Mars Rover API pour %s sol %d: %v", roverName, sol, err)
		return nil
	}

	var unifiedData []UnifiedImageData
	maxPhotos := 15 // Photos par sol et par rover
	processedPhotos := 0

	for _, photo := range roverResponse.Photos {
		if processedPhotos >= maxPhotos {
			break
		}

		imgURL := strings.Replace(photo.ImgSrc, "http://", "https://", 1)

		if isValidImage(imgURL) {
			unifiedData = append(unifiedData, UnifiedImageData{
				Title: fmt.Sprintf("%s (Sol %d): %s", photo.Rover.Name, sol, photo.Camera.FullName),
				URL:   imgURL,
				Date:  photo.EarthDate,
				Explanation: fmt.Sprintf("Photo prise sur Mars par le rover %s lors du sol %d avec la caméra %s",
					photo.Rover.Name, sol, photo.Camera.FullName),
				Source: fmt.Sprintf("Mars Rover - %s", photo.Rover.Name),
				Type:   "image",
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
		if len(gallery) >= 6 {
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
