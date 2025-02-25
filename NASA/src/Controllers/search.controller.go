package controllers

import (
	temp "Nasa/assets/Templates"
	"net/http"
	"strings"
	"unicode"
)

// SearchResult structure pour les résultats de recherche
type SearchResult struct {
	Items []UnifiedImageData
	Query string
	Total int
}

// SearchOptions structure pour configurer la recherche
type SearchOptions struct {
	Query         string
	SearchInTitle bool
	SearchInDesc  bool
	CaseSensitive bool
}

func SearchHandler(w http.ResponseWriter, r *http.Request) {
	// Récupérer le paramètre de recherche
	query := r.URL.Query().Get("q")
	if query == "" {
		http.Error(w, "Paramètre de recherche manquant", http.StatusBadRequest)
		return
	}

	// Configurer les options de recherche
	options := SearchOptions{
		Query:         query,
		SearchInTitle: true,
		SearchInDesc:  true,
		CaseSensitive: false,
	}

	// Effectuer la recherche
	results := performSearch(options)

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
		Items: results,
		Query: options.Query,
		Total: len(results),
	}
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
