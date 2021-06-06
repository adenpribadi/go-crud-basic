package controllers

import (

	"net/http"
	
	"text/template"

)


var hometmpl = template.Must(template.ParseGlob("views/home/*"))

//Index handler
func HomeIndex(w http.ResponseWriter, r *http.Request) {
	
	hometmpl.ExecuteTemplate(w, "Index", nil)
	
}
