package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/joho/godotenv"

	"github.com/buraksekili/getir-task/handlers"
	"github.com/buraksekili/getir-task/persistence/inmemory"
	mongodb "github.com/buraksekili/getir-task/persistence/mongo"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type config struct {
	HTTPPort string
	MongoURI string
}

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	conf := loadConfig()
	client, collection := initMongo(conf.MongoURI)
	defer client.Disconnect(context.Background())

	mr := mongodb.NewMongo(collection)
	im := initInMemory()

	logger := log.New(os.Stdout, "getir-task ", log.LstdFlags)

	a := handlers.NewHTTPAgent(logger, mr, im)

	sm := http.NewServeMux()
	sm.HandleFunc("/mongo", a.Mongo)
	sm.HandleFunc("/in-memory", a.InMemory)

	server := http.Server{
		Addr:         fmt.Sprintf(":%s", conf.HTTPPort),
		Handler:      sm,
		ErrorLog:     logger,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	go func() {
		logger.Printf("Starting server on port %s\n", conf.HTTPPort)

		err := server.ListenAndServe()
		if err != nil {
			logger.Printf("Error starting server: %s\n", err)
			os.Exit(1)
		}
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	signal.Notify(c, os.Kill)

	sig := <-c
	log.Println("Got signal:", sig)

	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	server.Shutdown(ctx)
}

func initMongo(mongoURI string) (*mongo.Client, *mongo.Collection) {
	client, err := mongo.NewClient(options.Client().ApplyURI(mongoURI))
	if err != nil {
		log.Fatal(err)
	}
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}
	return client, client.Database("getir-case-study").Collection("records")
}

func initInMemory() inmemory.InMemory {
	return inmemory.New(map[string]string{})
}

func loadConfig() config {
	return config{
		HTTPPort: os.Getenv("GETIR_TASK_HTTP_PORT"),
		MongoURI: os.Getenv("GETIR_TASK_MONGO_URI"),
	}
}
