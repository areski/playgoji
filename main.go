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

type Customer struct {
	Id   int
	Name string
	Age  int
}

// type CallData struct {
//     Id          int
//     Date        string
//     PhoneNumber string
//     HangupCause int
//     Duration    int
//     BillSec     int
// }

type Context struct {
	Db gorm.DB
}

func GetCustomers(res http.ResponseWriter, req *http.Request) {
	var customers []Customer
	ctx.Db.Find(&customers)
	data, err := json.Marshal(customers)
	if err != nil {
		log.Println("Error marshalling JSON")
	}
	res.Header().Set("Content-Type", "application/json")
	res.Write(data)
}

func GetCustomer(c web.C, res http.ResponseWriter, req *http.Request) {
	var customer Customer
	id, err := strconv.Atoi(c.URLParams["id"])
	if err != nil {
		log.Println("Error converting to integer")
	}
	ctx.Db.Where(&Customer{Id: id}).First(&customer)
	data, err := json.Marshal(customer)
	if err != nil {
		log.Println("Error marshaling JSON")
	}
	res.Header().Set("Content-Type", "application/json")
	res.Write(data)
}

func UpdateCustomer(c web.C, res http.ResponseWriter, req *http.Request) {
	var customer, newCustomer Customer
	id, err := strconv.Atoi(c.URLParams["id"])
	if err != nil {
		log.Println("Error converting to integer")
	}
	err = json.NewDecoder(req.Body).Decode(&newCustomer)
	if err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}
	ctx.Db.Where(&Customer{Id: id}).First(&customer)
	ctx.Db.Model(&customer).Updates(newCustomer)
	res.WriteHeader(http.StatusNoContent)
}

func DeleteCustomer(c web.C, res http.ResponseWriter, req *http.Request) {
	var customer Customer
	id, err := strconv.Atoi(c.URLParams["id"])
	if err != nil {
		log.Println("Error converting to integer")
	}
	ctx.Db.Where(&Customer{Id: id}).First(&customer)
	ctx.Db.Delete(&customer)
	res.WriteHeader(http.StatusNoContent)
}

func NewCustomer(res http.ResponseWriter, req *http.Request) {
	var customer Customer
	err := json.NewDecoder(req.Body).Decode(&customer)
	if err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}
	ctx.Db.Create(&customer)
	url := fmt.Sprintf("/customers/%d", customer.Id)
	http.Redirect(res, req, url, http.StatusCreated)
}

func init() {
	initDb()
}

func initDb() {
	db, err := gorm.Open("sqlite3", "/tmp/customers.db")
	if err != nil {
		log.Fatalf("Error opening database")
	}
	// db.DropTableIfExists(Customer{})
	db.CreateTable(Customer{})

	db.Create(&Customer{Id: 1, Name: "Calamere", Age: 23})
	db.Create(&Customer{Id: 2, Name: "Aslan", Age: 40})
	db.Create(&Customer{Id: 3, Name: "Shagrath", Age: 51})
	db.Create(&Customer{Id: 4, Name: "Troji", Age: 32})
	db.Create(&Customer{Id: 5, Name: "Raluk", Age: 35})
	ctx = Context{Db: db}
}

func main() {
	goji.Get("/", GetCustomers)
	goji.Get("/customers/:id", GetCustomer)
	goji.Post("/customers", NewCustomer)
	goji.Put("/customers/:id", UpdateCustomer)
	goji.Delete("/customers/:id", DeleteCustomer)
	goji.Serve()
}
