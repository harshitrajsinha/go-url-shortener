package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/harshitrajsinha/go-url-shortener/api/routes"
	"github.com/joho/godotenv"
	"github.com/supabase-community/supabase-go"
)

type responseMessage struct{
	Code int `json:"code"`
	Message string `json:"message"`
}

// middleware to pass Supabase db client to different route handlers
func SupabaseMiddleware(SupabaseClient *supabase.Client) mux.MiddlewareFunc {
	const SupabaseClientKey string = "SupabaseClient"
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := context.WithValue(r.Context(), SupabaseClientKey, SupabaseClient)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func handleHomeRoute(w http.ResponseWriter, r *http.Request){
	var serverMessage responseMessage = responseMessage{Code: http.StatusOK, Message: "Server is functioning"}
	json.NewEncoder(w).Encode(serverMessage)
}

func main(){

	_ = godotenv.Load()
	var port string = os.Getenv("PORT")
	
	if SupabaseClient == nil {
		log.Fatal("Supabase Client is not initialized!")
		return
	}
	
	// setup router for creating routes
	var router *mux.Router = mux.NewRouter()

	// Middleware to pass SupabaseClient to different route handlers
	router.Use(SupabaseMiddleware(SupabaseClient))

	// set routes
	router.HandleFunc("/", handleHomeRoute).Methods("GET")
	router.HandleFunc("/api/v1/shorturl", routes.HandleShortIdCreation).Methods("GET")
	
	fmt.Println("Listening at PORT ", port)
	log.Fatal(http.ListenAndServe(":"+port, router))
	
}