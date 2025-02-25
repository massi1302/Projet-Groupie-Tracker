package main

import (
	temp "Nasa/assets/Templates"
	_ "Nasa/src/Controllers"
	"Nasa/src/routes"
)

// main.go
func main() {
	// Initialise les templates HTML pour le rendu des pages
	temp.InitTemplates()
	// Démarre le serveur web et configure les routes
	routes.InitServe()
}
