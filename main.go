package main

import (
	"database/sql"
	"log"
	"net/http"
	"strconv"
	"time"

	game "github.com/HDIOES/hundredToOneBackend/rest/games"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"github.com/tkanos/gonfig"
)

type Configuration struct {
	DatabaseUrl        string `json:"databaseUrl"`
	MaxOpenConnections int    `json:"maxOpenConnections"`
	MaxIdleConnections int    `json:"maxIdleConnections"`
	ConnectionTimeout  int    `json:"connectionTimeout"`
	Port               int    `json:"port"`
}

func main() {

	configuration := Configuration{}
	gonfigErr := gonfig.GetConf("dbconfig.json", &configuration)
	if gonfigErr != nil {
		panic(gonfigErr)
	}

	db, err := sql.Open("postgres", configuration.DatabaseUrl)
	if err != nil {
		panic(err)
	}
	db.SetMaxIdleConns(configuration.MaxIdleConnections)
	db.SetMaxOpenConns(configuration.MaxOpenConnections)
	timeout := strconv.Itoa(configuration.ConnectionTimeout) + "s"
	timeoutDuration, durationErr := time.ParseDuration(timeout)
	if durationErr != nil {
		log.Println("Error parsing of timeout parameter")
		panic(durationErr)
	} else {
		db.SetConnMaxLifetime(timeoutDuration)
	}

	log.Println("Configuration has been loaded")

	router := mux.NewRouter()

	router.Handle("/game", game.CreateCreateGameHandler(db)).
		Methods("POST")
	router.Handle("/games", game.CreateSearchGamesHandler(db)).
		Methods("GET")

	http.Handle("/", router)
	headersOk := handlers.AllowedHeaders([]string{"X-Requested-With"})
	originsOk := handlers.AllowedOrigins([]string{"*"})
	methodsOk := handlers.AllowedMethods([]string{"GET", "HEAD", "POST", "PUT", "OPTIONS"})

	listenandserveErr := http.ListenAndServe(":"+strconv.Itoa(configuration.Port), handlers.CORS(originsOk, headersOk, methodsOk)(router))
	if listenandserveErr != nil {
		panic(err)
	}

}
