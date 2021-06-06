package main

import (
	"log"
	"net/http"
	"os"
	"go-crud-basic/controllers"
)

func main() {
	port := os.Getenv("PORT")
	// port := "8080"
	log.Println("Server started on: http://localhost:"+port)
	http.HandleFunc("/", controllers.HomeIndex)
	http.HandleFunc("/tools/", controllers.ToolIndex)
	http.HandleFunc("/tools/show", controllers.ToolShow)
	http.HandleFunc("/tools/new", controllers.ToolNew)
	http.HandleFunc("/tools/edit", controllers.ToolEdit)
	http.HandleFunc("/tools/insert", controllers.ToolInsert)
	http.HandleFunc("/tools/update", controllers.ToolUpdate)
	http.HandleFunc("/tools/delete", controllers.ToolDelete)
	http.ListenAndServe(":"+port, nil)
}
