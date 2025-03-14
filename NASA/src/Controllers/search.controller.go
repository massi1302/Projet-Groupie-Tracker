package controllers

import (
	temp "Nasa/assets/Templates"
	"net/http"
	"strings"
	"time"
	"unicode"
)

// SearchResult structure pour les résultats de recherche
type SearchResult struct {
	Items          []UnifiedImageData
	Query          string
	Total          int
	Filters        FilterOptions
	AppliedFilters int
}

// SearchOptions structure pour configurer la recherche
type SearchOptions struct {
	Query         string
	SearchInTitle bool
	SearchInDesc  bool
	CaseSensitive bool
	Filters       FilterOptions
}

// FilterOptions structure pour les filtres
type FilterOptions struct {
	Source           string   // NASA, APOD, Mars Rover, etc.
	MediaType        string   // image, video
	DateFrom         string   // Format: YYYY-MM-DD
	DateTo           string   // Format: YYYY-MM-DD
	AvailableSources []string // Sources disponibles pour le filtre
	AvailableTypes   []string // Types de média disponibles
}

func SearchHandler(w http.ResponseWriter, r *http.Request) {
	// Récupérer le paramètre de recherche
	query := r.URL.Query().Get("q")

	// Récupérer les paramètres de filtre
	source := r.URL.Query().Get("source")
	mediaType := r.URL.Query().Get("type")
	dateFrom := r.URL.Query().Get("from")
	dateTo := r.URL.Query().Get("to")

	// Compter le nombre de filtres appliqués
	appliedFilters := 0
	if source != "" {
		appliedFilters++
	}
	if mediaType != "" {
		appliedFilters++
	}
	if dateFrom != "" || dateTo != "" {
		appliedFilters++
	}

	// Définir les filtres disponibles (à remplir dynamiquement en production)
	availableSources := []string{"APOD", "Mars Rover", "NASA Gallery"}
	availableTypes := []string{"image", "video"}

	// Configurer les options de filtrage
	filters := FilterOptions{
		Source:           source,
		MediaType:        mediaType,
		DateFrom:         dateFrom,
		DateTo:           dateTo,
		AvailableSources: availableSources,
		AvailableTypes:   availableTypes,
	}

	// Configurer les options de recherche
	options := SearchOptions{
		Query:         query,
		SearchInTitle: true,
		SearchInDesc:  true,
		CaseSensitive: false,
		Filters:       filters,
	}

	// Effectuer la recherche
	results := performSearch(options)
	results.AppliedFilters = appliedFilters

	// Rendre le template avec les résultats
	data := struct {
		Results SearchResult
	}{
		Results: results,
	}

	if err := temp.Temp.ExecuteTemplate(w, "search", data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func performSearch(options SearchOptions) SearchResult {
	// Récupérer toutes les ressources disponibles
	allItems := getAllResources()

	var results []UnifiedImageData
	query := options.Query
	if !options.CaseSensitive {
		query = strings.ToLower(query)
	}

	// Normaliser la requête pour la recherche
	normalizedQuery := normalizeString(query)
	searchTerms := strings.Fields(normalizedQuery)

	for _, item := range allItems {
		// Appliquer les filtres
		if !applyFilters(item, options.Filters) {
			continue
		}

		// Si pas de requête de recherche mais des filtres, ajouter l'item
		if query == "" {
			results = append(results, item)
			continue
		}

		isMatch := false

		// Préparer les champs pour la recherche
		title := item.Title
		explanation := item.Explanation

		if !options.CaseSensitive {
			title = strings.ToLower(title)
			explanation = strings.ToLower(explanation)
		}

		// Normaliser les champs pour la recherche
		normalizedTitle := normalizeString(title)
		normalizedExplanation := normalizeString(explanation)

		// Rechercher tous les termes
		for _, term := range searchTerms {
			// Recherche dans le titre
			if options.SearchInTitle && strings.Contains(normalizedTitle, term) {
				isMatch = true
				break
			}

			// Recherche dans la description
			if options.SearchInDesc && strings.Contains(normalizedExplanation, term) {
				isMatch = true
				break
			}
		}

		if isMatch {
			results = append(results, item)
		}
	}

	return SearchResult{
		Items:   results,
		Query:   options.Query,
		Total:   len(results),
		Filters: options.Filters,
	}
}

// applyFilters vérifie si un élément correspond aux critères de filtrage
func applyFilters(item UnifiedImageData, filters FilterOptions) bool {
	// Filtre par source
	if filters.Source != "" && item.Source != filters.Source {
		return false
	}

	// Filtre par type de média
	if filters.MediaType != "" && item.Type != filters.MediaType {
		return false
	}

	// Filtre par date
	if filters.DateFrom != "" || filters.DateTo != "" {
		itemDate, err := time.Parse("2006-01-02", item.Date)
		if err != nil {
			// Si la date n'est pas valide, on ignore ce filtre pour cet élément
			return true
		}

		if filters.DateFrom != "" {
			fromDate, err := time.Parse("2006-01-02", filters.DateFrom)
			if err == nil && itemDate.Before(fromDate) {
				return false
			}
		}

		if filters.DateTo != "" {
			toDate, err := time.Parse("2006-01-02", filters.DateTo)
			if err == nil && itemDate.After(toDate) {
				return false
			}
		}
	}

	return true
}

// normalizeString prépare une chaîne pour la recherche
func normalizeString(s string) string {
	// Convertir en minuscules
	s = strings.ToLower(s)

	// Remplacer les caractères accentués
	replacements := map[rune]string{
		'é': "e", 'è': "e", 'ê': "e", 'ë': "e",
		'à': "a", 'â': "a", 'ä': "a",
		'î': "i", 'ï': "i",
		'ô': "o", 'ö': "o",
		'ù': "u", 'û': "u", 'ü': "u",
		'ÿ': "y",
		'ç': "c",
	}

	normalized := strings.Builder{}
	for _, r := range s {
		if replacement, ok := replacements[r]; ok {
			normalized.WriteString(replacement)
		} else if unicode.IsLetter(r) || unicode.IsNumber(r) || unicode.IsSpace(r) {
			normalized.WriteRune(r)
		}
	}

	return normalized.String()
}
