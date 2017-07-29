package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"golang.org/x/net/context"

	"github.com/echojc/lernkartei/dict"
	"github.com/echojc/lernkartei/dict/glosbe"
)

func main() {
	port := flag.Int("port", 9000, "Listen port")
	flag.Parse()

	server := http.Server{
		Addr:    fmt.Sprintf(":%d", *port),
		Handler: setupRoutes(),
	}
	go func() {
		fmt.Println("Listening...")
		err := server.ListenAndServe()
		fmt.Println(err)
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)

	<-quit
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	server.Shutdown(ctx)
}

func setupRoutes() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", lookup)
	return mux
}

func lookup(w http.ResponseWriter, r *http.Request) {
	words := r.URL.Query()["word"]

	// limit to 10 words per request
	if len(words) > 10 {
		words = words[:10]
	}

	out := []dict.Word{}
	for _, word := range words {
		words, err := glosbe.Lookup(word)
		if err != nil {
			log.Printf("word [%+v] err [%+v]", word, err)
			continue
		}
		out = append(out, words...)
	}

	bs, err := json.Marshal(out)
	if err != nil {
		log.Printf("words [%+v] out [%+v] err [%+v]", words, out, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	h := w.Header()
	h.Add("Content-Type", "application/json")
	h.Add("Access-Control-Allow-Origin", "*")

	w.WriteHeader(http.StatusOK)
	w.Write(bs)
}
