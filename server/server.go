package server

import (
	"api/config"
	"html/template"
	"net/http"
)

var templates *template.Template

func addOne(i int) int {
	return i + 1
}

func init() {
	templates = template.Must(template.New("hangman").Funcs(template.FuncMap{"addOne": addOne}).ParseGlob(config.App.Server.StaticWeb.Template.Dir))
}

func indexHandler(responseWriter http.ResponseWriter, request *http.Request) {
	if err := templates.ExecuteTemplate(responseWriter, "index", nil); err != nil {
		http.Error(responseWriter, err.Error(), http.StatusInternalServerError)
	}
}

func ServeMux() *http.ServeMux {
	serveMux := http.NewServeMux()

	fileServer := http.FileServer(http.Dir(config.App.Server.StaticWeb.Assets.Dir))
	serveMux.Handle("/assets/", http.StripPrefix("/assets/", fileServer))

	serveMux.HandleFunc("/", indexHandler)
	serveMux.HandleFunc("/index", indexHandler)

	return serveMux
}
