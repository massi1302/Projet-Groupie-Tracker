package controllers

import (
	temp "Nasa/assets/Templates"
	"fmt"
	"log"
	"net/http"
	"strings"
	"sync"
)

func HomePage(w http.ResponseWriter, r *http.Request) {
	log.Println("Début du traitement de la requête HomePage")

	var pageData HomePageData
	var wg sync.WaitGroup
	var mu sync.Mutex
	errChan := make(chan error, 4) // Mis à jour pour inclure les topics

	wg.Add(4)

	// APOD
	go func() {
		defer wg.Done()
		apodData := getAPODData()
		if len(apodData) == 0 {
			errChan <- fmt.Errorf("APOD error: no data")
			return
		}
		mu.Lock()
		pageData.APOD = APODResponse{
			Title:       apodData[0].Title,
			URL:         apodData[0].URL,
			Explanation: apodData[0].Explanation,
			Date:        apodData[0].Date,
		}
		mu.Unlock()
	}()

	// NASA News RSS
	go func() {
		defer wg.Done()
		news, err := getNASANews()
		if err != nil {
			errChan <- fmt.Errorf("news error: %v", err)
			return
		}
		mu.Lock()
		pageData.News = news
		mu.Unlock()
	}()

	// NASA Image Library
	go func() {
		defer wg.Done()
		gallery, err := getNASAGallery()
		if err != nil {
			errChan <- fmt.Errorf("gallery error: %v", err)
			return
		}
		mu.Lock()
		pageData.Gallery = gallery
		mu.Unlock()
	}()

	// NASA Topics
	go func() {
		defer wg.Done()
		topics, err := getNASATopics()
		if err != nil {
			errChan <- fmt.Errorf("topics error: %v", err)
			return
		}
		mu.Lock()
		pageData.Topics = topics
		mu.Unlock()
	}()

	wg.Wait()
	close(errChan)

	for err := range errChan {
		log.Printf("Error fetching data: %v", err)
	}
	for i, news := range pageData.News {
		log.Printf("Article %d - Title: %s, ImageURL: %s", i, news.Title, news.ImageURL)
	}
	if err := temp.Temp.ExecuteTemplate(w, "home", pageData); err != nil {
		log.Printf("Erreur lors du rendu du template: %v", err)
		http.Error(w, "Erreur lors du rendu de la page", http.StatusInternalServerError)
		return
	}
}

func HomeSearchHandler(w http.ResponseWriter, r *http.Request) {
	query := r.FormValue("q")

	// Récupérer d'éventuels filtres initiaux (facultatif)
	source := r.FormValue("source")
	mediaType := r.FormValue("type")

	// Construire l'URL de redirection
	redirectURL := "/search"

	// Ajouter les paramètres
	if query != "" || source != "" || mediaType != "" {
		redirectURL += "?"

		params := []string{}
		if query != "" {
			params = append(params, "q="+query)
		}
		if source != "" {
			params = append(params, "source="+source)
		}
		if mediaType != "" {
			params = append(params, "type="+mediaType)
		}

		redirectURL += strings.Join(params, "&")
	} else {
		// Si pas de recherche ou filtres, rediriger vers la page d'accueil
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	// Rediriger vers la page de recherche avec les paramètres
	http.Redirect(w, r, redirectURL, http.StatusSeeOther)
}
