package controllers

import (
	"database/sql"
	"log"
	"net/http"
	"text/template"

  "gopkg.in/yaml.v2"
  "io/ioutil"
	
	"go-crud-basic/models"

	_ "github.com/go-sql-driver/mysql"
)

type conf struct {
  Server struct {
    Port string `yaml:"port"`
    Host string `yaml:"host"`
  } `yaml:"server"`
  Database struct {
    DbName string `yaml:"db_name"`
    Username string `yaml:"user"`
    Password string `yaml:"pass"`
  } `yaml:"database"`
}

func (c *conf) getConf() *conf {

  yamlFile, err := ioutil.ReadFile("config.yaml")
  if err != nil {
      log.Printf("yamlFile.Get err   #%v ", err)
  }
  err = yaml.Unmarshal(yamlFile, c)
  if err != nil {
      log.Fatalf("Unmarshal: %v", err)
  }

  return c
}

func dbConn() (db *sql.DB) {
	var c conf
  c.getConf()

	dbDriver := "mysql"
	dbUser := c.Database.Username
	dbPass := c.Database.Password
	dbName := c.Database.DbName
	dbServer := c.Server.Host
	dbPort := c.Server.Port
	db, err := sql.Open(dbDriver, dbUser+":"+dbPass+"@tcp("+dbServer+":"+dbPort+")/"+dbName)
	if err != nil {
		panic(err.Error())
	}
	return db
}

var tooltmpl = template.Must(template.ParseGlob("views/tools/*"))

//Index handler
func ToolIndex(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	selDB, err := db.Query("SELECT * FROM tools ORDER BY id DESC")
	if err != nil {
		panic(err.Error())
	}

	tool := models.Tool{}
	res := []models.Tool{}

	for selDB.Next() {
		var id, rating int
		var name, category, url, notes string
		err := selDB.Scan(&id, &name, &category, &url, &rating, &notes)
		if err != nil {
			panic(err.Error())
		}
		log.Println("Listing Row: Id " + string(id) + " | name " + name + " | category " + category + " | url " + url + " | rating " + string(rating) + " | notes " + notes)

		tool.Id = id
		tool.Name = name
		tool.Category = category
		tool.URL = url
		tool.Rating = rating
		tool.Notes = notes
		res = append(res, tool)
	}
	tooltmpl.ExecuteTemplate(w, "Index", res)
	defer db.Close()
}

//Show handler
func ToolShow(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	nId := r.URL.Query().Get("id")
	selDB, err := db.Query("SELECT * FROM tools WHERE id=?", nId)
	if err != nil {
		panic(err.Error())
	}

	tool := models.Tool{}

	for selDB.Next() {
		var id, rating int
		var name, category, url, notes string
		err := selDB.Scan(&id, &name, &category, &url, &rating, &notes)
		if err != nil {
			panic(err.Error())
		}

		log.Println("Showing Row: Id " + string(id) + " | name " + name + " | category " + category + " | url " + url + " | rating " + string(rating) + " | notes " + notes)

		tool.Id = id
		tool.Name = name
		tool.Category = category
		tool.URL = url
		tool.Rating = rating
		tool.Notes = notes
	}
	tooltmpl.ExecuteTemplate(w, "Show", tool)
	defer db.Close()
}

func ToolNew(w http.ResponseWriter, r *http.Request) {
	tooltmpl.ExecuteTemplate(w, "New", nil)
}

func ToolEdit(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	nId := r.URL.Query().Get("id")
	selDB, err := db.Query("SELECT * FROM tools WHERE id=?", nId)
	if err != nil {
		panic(err.Error())
	}

	tool := models.Tool{}

	for selDB.Next() {
		var id, rating int
		var name, category, url, notes string
		err := selDB.Scan(&id, &name, &category, &url, &rating, &notes)
		if err != nil {
			panic(err.Error())
		}

		tool.Id = id
		tool.Name = name
		tool.Category = category
		tool.URL = url
		tool.Rating = rating
		tool.Notes = notes
	}

	tooltmpl.ExecuteTemplate(w, "Edit", tool)
	defer db.Close()
}

func ToolInsert(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	if r.Method == "POST" {
		name := r.FormValue("name")
		category := r.FormValue("category")
		url := r.FormValue("url")
		rating := r.FormValue("rating")
		notes := r.FormValue("notes")
		insForm, err := db.Prepare("INSERT INTO tools (name, category, url, rating, notes) VALUES (?, ?, ?, ?, ?)")
		if err != nil {
			panic(err.Error())
		}
		insForm.Exec(name, category, url, rating, notes)
		log.Println("Insert Data: name " + name + " | category " + category + " | url " + url + " | rating " + rating + " | notes " + notes)
	}
	defer db.Close()
	http.Redirect(w, r, "/tools/", 301)
}

func ToolUpdate(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	if r.Method == "POST" {
		name := r.FormValue("name")
		category := r.FormValue("category")
		url := r.FormValue("url")
		rating := r.FormValue("rating")
		notes := r.FormValue("notes")
		id := r.FormValue("uid")
		insForm, err := db.Prepare("UPDATE tools SET name=?, category=?, url=?, rating=?, notes=? WHERE id=?")
		if err != nil {
			panic(err.Error())
		}
		insForm.Exec(name, category, url, rating, notes, id)
		log.Println("UPDATE Data: name " + name + " | category " + category + " | url " + url + " | rating " + rating + " | notes " + notes)
	}
	defer db.Close()
	http.Redirect(w, r, "/tools/", 301)
}

func ToolDelete(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	tool := r.URL.Query().Get("id")
	delForm, err := db.Prepare("DELETE FROM tools WHERE id=?")
	if err != nil {
		panic(err.Error())
	}
	delForm.Exec(tool)
	log.Println("DELETE " + tool)
	defer db.Close()
	http.Redirect(w, r, "/tools/", 301)
}