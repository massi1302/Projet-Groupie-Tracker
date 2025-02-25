package routes

import (
	"fmt"
	"net/http"
)

func InitServe() {
	// Initialise toutes les routes de l'application
	homeRoutes()
	collectionRoutes()
	resourceRoutes()
	categoryRoutes()
	searchRoutes()
	// Affiche un message de confirmation dans la console
	fmt.Println("Le serveur est opérationel : http://localhost:8080")
	// Configure le serveur sur le port 8080 avec gestion d'erreur
	err := http.ListenAndServe("localhost:8080", nil)
	if err != nil {
		fmt.Printf("Erreur lors du démarrage du serveur: %v\n", err)
	}
}
