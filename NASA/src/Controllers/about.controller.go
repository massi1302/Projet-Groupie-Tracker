package controllers

import (
	temp "Nasa/assets/Templates"
	"net/http"
)

// AboutPage affiche la page "À propos"
// de l'application Web NASA
func AboutPage(w http.ResponseWriter, r *http.Request) {
	if err := temp.Temp.ExecuteTemplate(w, "about", nil); err != nil {
		http.Error(w, "Erreur lors du rendu de la page", http.StatusInternalServerError)
	}
}
