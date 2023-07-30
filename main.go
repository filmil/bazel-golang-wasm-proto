package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/grahamjenson/bazel-golang-wasm-protoc/protos/api"
	"github.com/grahamjenson/bazel-golang-wasm-protoc/server"
	"github.com/maxence-charriere/go-app/v9/pkg/app"
)

var (
	bootstrapLoc = flag.String("bootstrap-css-path", "", "path to the bootstrap.css file")
	wasmLoc      = flag.String("wasm-path", "", "path to the web app wasm file")
)

func main() {
	flag.Parse()

	if *bootstrapLoc == "" {
		log.Fatalf("The flag --bootstrap-css-path is required.")
	}
	if *wasmLoc == "" {
		log.Fatalf("The flag --bootstrap-css-path is required.")
	}

	currentDir, err := os.Getwd()
	if err != nil {
		fmt.Printf("boops, can't find current dir: %v\n", err)
	}
	fmt.Printf("current dir: %v\n", currentDir)

	mux := http.NewServeMux()

	a := &app.Handler{
		Title:  "EC2Instances",
		Author: "Graham Jenson",
		Styles: []string{"bootstrap.css"},
	}

	mux.HandleFunc("/app.wasm", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, *wasmLoc)
	})

	mux.HandleFunc("/bootstrap.css", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, *bootstrapLoc)
	})

	// Handle API
	api.RegisterApiHTTPMux(mux, &server.Server{})

	// Handle go-app
	mux.Handle("/", a)

	app.RunWhenOnBrowser()

	fmt.Println("starting local server on http://localhost:7000")
	log.Fatal(http.ListenAndServe(":7000", mux))
}
