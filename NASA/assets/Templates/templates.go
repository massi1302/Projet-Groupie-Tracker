package temp

import (
	"fmt"
	"html/template"
	"net/http"
	"os"
)

var Temp *template.Template

// Fonctions helper pour les templates
var templateFuncs = template.FuncMap{
	"iterate": func(count int) []int {
		var items []int
		for i := 0; i < count; i++ {
			items = append(items, i)
		}
		return items
	},
	"add": func(a, b int) int {
		return a + b
	},
	"subtract": func(a, b int) int {
		return a - b
	},
	"min": func(a, b int) int {
		if a < b {
			return a
		}
		return b
	},
	"max": func(a, b int) int {
		if a > b {
			return a
		}
		return b
	},
}

// InitTemplates charge les templates avec les fonctions helper
func InitTemplates() {
	fileServer := http.FileServer(http.Dir("./assets"))
	http.Handle("/assets/", http.StripPrefix("/assets/", fileServer))

	temp, tempErr := template.New("").Funcs(templateFuncs).ParseGlob("./assets/Templates/*.html")
	if tempErr != nil {
		fmt.Printf("Erreur Template - Une erreur lors du chargement des templates \n message d'erreur : %v\n", tempErr.Error())
		os.Exit(1)
	}
	Temp = temp
}
