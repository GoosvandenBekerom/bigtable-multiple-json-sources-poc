package main

import (
	"context"
	"flag"
	"github.com/GoosvandenBekerom/intaker-bigtable-poc/data"
	"github.com/GoosvandenBekerom/intaker-bigtable-poc/endpoints"
	"log"
	"net/http"
	"os"
)

func main() {
	project := flag.String("project", "fake-local-project", "gcp project for bigtable instance")
	instance := flag.String("instance", "products", "bigtable instance")
	flag.Parse()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	store := setupProductStore(ctx, *project, *instance)
	api := endpoints.NewProductAPI(store)

	mux := http.NewServeMux()
	mux.HandleFunc("/products/generate", api.GenerateProducts)
	mux.HandleFunc("/products", api.GetProducts)

	log.Println("=============================")
	log.Println("Running...")

	err := http.ListenAndServe("localhost:8080", mux)
	if err != nil {
		panic(err)
	}
}

func setupProductStore(ctx context.Context, project, instance string) data.ProductStore {
	// Make sure we connect to the dockerized emulator
	err := os.Setenv("BIGTABLE_EMULATOR_HOST", "localhost:8086")
	if err != nil {
		panic(err)
	}

	store, err := data.NewProductStore(ctx, project, instance)
	if err != nil {
		panic(err)
	}

	return store
}
