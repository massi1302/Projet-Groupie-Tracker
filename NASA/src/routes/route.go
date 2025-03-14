package routes

import (
	controllers "Nasa/src/Controllers"
	"net/http"
)

func homeRoutes() {
	http.HandleFunc("/", controllers.HomePage)
	http.HandleFunc("/home-search", controllers.HomeSearchHandler)
}

func collectionRoutes() {
	http.HandleFunc("/collection", controllers.CollectionPage)
}

func resourceRoutes() {
	http.HandleFunc("/resource/", controllers.ResourceDetailPage)
}

func categoryRoutes() {
	// Modifier cette ligne pour correspondre au chemin utilisé dans le contrôleur
	http.HandleFunc("/categories/", controllers.CategoryPage)
	// Ajouter également la route principale des catégories
	http.HandleFunc("/categories", controllers.CategoryPage)
}

func searchRoutes() {
	http.HandleFunc("/search", controllers.SearchHandler)
}

func aboutRoutes() {
	http.HandleFunc("/about", controllers.AboutPage)
}

func favorisRoutes() {
	http.HandleFunc("/favoris", controllers.FavorisPage)
	http.HandleFunc("/favoris/add/", controllers.AddToFavoris)
	http.HandleFunc("/favoris/remove/", controllers.RemoveFromFavoris)
}
