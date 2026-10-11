package main

import (
	"time"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"github.com/surajklmn/hr-system/internal/database"
	"github.com/surajklmn/hr-system/internal/employee"
	"github.com/surajklmn/hr-system/internal/leave"
)

func health(w http.ResponseWriter, r *http.Request){
	data := map[string]string{"status":"ok"}
	byte,err := json.Marshal(data)
	if err!= nil{
		fmt.Println("Couldn't encode/marshal the data",err)
		http.Error(w,"",http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type","application/json")
	w.Write(byte)
}

func main(){

	//This section is for reading envs
	port := os.Getenv("PORT")
	if port == ""{
		port = "8080"
	}
	dbPath := os.Getenv("DB_PATH")
	if dbPath == ""{
		dbPath = "hr.db"
	}

	db, err := database.Open(dbPath)
	if err != nil{
		log.Fatal(err)
	}

	defer db.Close()
	log.Println("Connected to SQLite database")

	mux := http.NewServeMux()//Our router instead of default one
	empRepo := employee.NewRepository(db)
	empService := employee.NewService(empRepo)
	empHandler := employee.NewHandler(empService)
	empHandler.RegisterRoutes(mux)

// Wire leave moduel
	leaveRepo := leave.NewRepository(db)
	leaveService := leave.NewService(leaveRepo, empService)
	leaveHandler := leave.NewHandler(leaveService)
	leaveHandler.RegisterRoutes(mux)

	mux.HandleFunc("GET /health",health)

	log.Printf("Server running on port %s", port)
	server := &http.Server{
		Addr: ":" + port,
		Handler: mux,
		ReadTimeout: 5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout: time.Minute,
	}

	log.Fatal(server.ListenAndServe()) 

}
