package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

var db *gorm.DB
var err error

type Result struct {
	Data interface{} `json:"data"`
	Message string `json:"message"`
}

type Konsumen struct {
    ID      int    `json:"id" gorm:"primaryKey"`
    Nama    string `json:"nama"`
    TglLahir string `json:"tgl_lahir"`
    Umur    int    `json:"umur"`
    Alamat  string `json:"alamat"`
}

func main() {
	db, err = gorm.Open("mysql", "root:@/db_konsumen?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		log.Println("Connection failed", err)
	} else {
		log.Println("Connection established")
	}

	db.AutoMigrate(&Konsumen{})
	handleRequests()
}

func handleRequests()  {
	log.Println("Start the development server at http://127.0.0.1:8080")

	myRouter := mux.NewRouter().StrictSlash(true)

	myRouter.HandleFunc("/", homePage)
	myRouter.HandleFunc("/konsumen/create", createKonsumen).Methods("POST")
	myRouter.HandleFunc("/konsumen/create", getKonsumen).Methods("GET")
	myRouter.HandleFunc("/konsumen/create/{id}", getKonsumens).Methods("GET")
	myRouter.HandleFunc("/konsumen/create/{id}", updateKonsumens).Methods("PUT")
	myRouter.HandleFunc("/konsumen/create/{id}", deleteKonsumen).Methods("DELETE")

	log.Fatal(http.ListenAndServe(":8080", myRouter))
}

func homePage(w http.ResponseWriter, r *http.Request)  {
	fmt.Fprintf(w, "Welcome! This is my first golang crud")
}

func createKonsumen(w http.ResponseWriter, r *http.Request) {
    payloads, _ := ioutil.ReadAll(r.Body)

    var konsumen Konsumen
    if err := json.Unmarshal(payloads, &konsumen); err != nil {
        log.Println("Error unmarshalling JSON:", err)
        http.Error(w, "Failed to process JSON payload", http.StatusBadRequest)
        return
    }

    if err := db.Create(&konsumen).Error; err != nil {
        log.Println("Error creating Konsumen:", err)
        http.Error(w, "Failed to create Konsumen", http.StatusInternalServerError)
        return
    }

    res := Result{Data: konsumen, Message: "Success create data Konsumen"}
    result, err := json.Marshal(res)
    if err != nil {
        log.Println("Error marshalling JSON:", err)
        http.Error(w, "Failed to process response", http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    w.Write(result)
}

func getKonsumen(w http.ResponseWriter, r *http.Request)  {
	konsumens := []Konsumen{}

	db.Find(&konsumens)
	
	res := Result{Data: konsumens, Message: "Succes get data konsumen"}
	results, err := json.Marshal(res)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    w.Write(results)
}

func getKonsumens(w http.ResponseWriter, r *http.Request)  {
	vars := mux.Vars(r)
	konsID := vars["id"]

	var konsumen Konsumen
	db.First(&konsumen, konsID)
	
	res := Result{Data: konsumen, Message: "Succes get data konsumen"}
	result, err := json.Marshal(res)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    w.Write(result)
}

func updateKonsumens(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	konsID := vars["id"]

    payloads, _ := ioutil.ReadAll(r.Body)

    var konsumenUpdate Konsumen
    json.Unmarshal(payloads, &konsumenUpdate)

	var konsumen Konsumen
	db.First(&konsumen, konsID)
	db.Model(&konsumen).Updates(konsumenUpdate)

    res := Result{Data: konsumen, Message: "Success update data Konsumen"}
    result, err := json.Marshal(res)

    if err != nil {
        log.Println("Error marshalling JSON:", err)
        http.Error(w, "Failed to process response", http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    w.Write(result)
}

func deleteKonsumen(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	konsID := vars["id"]

	var konsumen Konsumen
	db.First(&konsumen, konsID)
	db.Delete(&konsumen)

	res := Result{Message: "Success delete data Konsumen"}
    result, err := json.Marshal(res)

    if err != nil {
        log.Println("Error marshalling JSON:", err)
        http.Error(w, "Failed to process response", http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    w.Write(result)
}