package main

import (
	"database/sql"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"

	"serv/internal/api"
	"serv/internal/db"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/gorilla/mux"
	"gopkg.in/yaml.v3"
)

type Config struct {
	DsnHTTP  string `yaml:"dsnHTTP"`
	DSN      string `yaml:"dsn"`
	LogLevel int    `yaml:"log_level"`
}

func main() {

	yamlContent, err := os.ReadFile("./../../config.yml")
	if err != nil {
		fmt.Println("problems")
	}
	var systemConfig Config
	err = yaml.Unmarshal(yamlContent, &systemConfig)
	if err != nil {
		fmt.Println("problems")
	}

	// создаём миграции
	rawSQLConn, err := sql.Open("postgres", systemConfig.DSN) // открываем доступ к базе данных postgreSQL
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	m, err := migrate.New( // сохраняем в переменную m данные из папки migrate
		"file://../../migrate", systemConfig.DSN)
	if err != nil {

		log.Fatal(err)
	}
	if err := m.Up(); err != nil && err != migrate.ErrNoChange { // отправляем запрос из переменной m в базу данных sql из папки migrate
		fmt.Println("миграция не прошла")
		log.Fatal(err)
	}
	// _________________________________________________________________________________________________________________

	repo := db.NewRepository(rawSQLConn) // передаём доступ слою db  к базе данных

	var store api.Store = repo
	//store = repo

	ourServer := api.Server{
		//Database: repo,
		Database: store,
	}

	log := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.Level(systemConfig.LogLevel),
	}))

	r := mux.NewRouter()

	r.Use(api.Logg(log, ourServer))

	r.HandleFunc("/book", ourServer.AddBook).Methods(http.MethodPost)
	r.HandleFunc("/book", ourServer.GetBook).Methods(http.MethodGet)
	r.HandleFunc("/book", ourServer.DeleteBook).Methods(http.MethodDelete)
	r.HandleFunc("/book", ourServer.UpdateBook).Methods(http.MethodPut)
	r.HandleFunc("/books", ourServer.GetAllBooks).Methods(http.MethodGet)
	r.HandleFunc("/user", ourServer.AddUser).Methods(http.MethodPost)
	r.HandleFunc("/user", ourServer.LoginUser).Methods(http.MethodGet)
	log.Debug("сервер запущен")
	err = http.ListenAndServe(systemConfig.DsnHTTP, r)
	if err != nil {
		fmt.Println("problem ListenAndServe")
	}
}
