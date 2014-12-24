//
//
package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/jinzhu/gorm"
	_ "github.com/mattn/go-sqlite3"
	"github.com/zenazn/goji"
	"github.com/zenazn/goji/web"
)

var ctx Context

type Actor struct {
	Id   int
	Name string
	Age  int
}

type Context struct {
	Db gorm.DB
}

func GetActors(res http.ResponseWriter, req *http.Request) {
	var actors []Actor
	ctx.Db.Find(&actors)
	data, err := json.Marshal(actors)
	if err != nil {
		log.Println("Error marshalling JSON")
	}
	res.Header().Set("Content-Type", "application/json")
	res.Write(data)
}

func GetActor(c web.C, res http.ResponseWriter, req *http.Request) {
	var actor Actor
	id, err := strconv.Atoi(c.URLParams["id"])
	if err != nil {
		log.Println("Error converting to integer")
	}
	ctx.Db.Where(&Actor{Id: id}).First(&actor)
	data, err := json.Marshal(actor)
	if err != nil {
		log.Println("Error marshaling JSON")
	}
	res.Header().Set("Content-Type", "application/json")
	res.Write(data)
}

func UpdateActor(c web.C, res http.ResponseWriter, req *http.Request) {
	var actor, newActor Actor
	id, err := strconv.Atoi(c.URLParams["id"])
	if err != nil {
		log.Println("Error converting to integer")
	}
	err = json.NewDecoder(req.Body).Decode(&newActor)
	if err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}
	ctx.Db.Where(&Actor{Id: id}).First(&actor)
	ctx.Db.Model(&actor).Updates(newActor)
	res.WriteHeader(http.StatusNoContent)
}

func DeleteActor(c web.C, res http.ResponseWriter, req *http.Request) {
	var actor Actor
	id, err := strconv.Atoi(c.URLParams["id"])
	if err != nil {
		log.Println("Error converting to integer")
	}
	ctx.Db.Where(&Actor{Id: id}).First(&actor)
	ctx.Db.Delete(&actor)
	res.WriteHeader(http.StatusNoContent)
}

func NewActor(res http.ResponseWriter, req *http.Request) {
	var actor Actor
	err := json.NewDecoder(req.Body).Decode(&actor)
	if err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}
	ctx.Db.Create(&actor)
	url := fmt.Sprintf("/actors/%d", actor.Id)
	http.Redirect(res, req, url, http.StatusCreated)
}

func init() {
	initDb()
}

func initDb() {
	db, err := gorm.Open("sqlite3", "/tmp/actors.db")
	if err != nil {
		log.Fatalf("Error opening database")
	}
	// db.DropTableIfExists(Actor{})
	db.CreateTable(Actor{})

	db.Create(&Actor{Id: 1, Name: "Calamere", Age: 23})
	db.Create(&Actor{Id: 2, Name: "Aslan", Age: 40})
	db.Create(&Actor{Id: 3, Name: "Shagrath", Age: 51})
	db.Create(&Actor{Id: 4, Name: "Troji", Age: 32})
	db.Create(&Actor{Id: 5, Name: "Raluk", Age: 35})
	ctx = Context{Db: db}
}

func main() {
	goji.Get("/", GetActors)
	goji.Get("/actors/:id", GetActor)
	goji.Post("/actors", NewActor)
	goji.Put("/actors/:id", UpdateActor)
	goji.Delete("/actors/:id", DeleteActor)
	goji.Serve()
}
