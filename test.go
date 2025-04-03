package handler

// import (
// 	"context"
// 	"encoding/json"
// 	"errors"
// 	"log"
// 	"net/http"
// 	"os"
// 	"runtime/debug"

// 	"github.com/gorilla/mux"
// 	"github.com/harshitrajsinha/go-url-shortener/api/v1/routes"
// 	"github.com/harshitrajsinha/go-url-shortener/database"
// 	"github.com/joho/godotenv"
// 	"github.com/supabase-community/supabase-go"
// )

// // middleware - to pass Supabase db client to different route handlers
// func SupabaseMiddleware(SupabaseClient *supabase.Client) mux.MiddlewareFunc {
// 	const SupabaseClientKey string = "SupabaseClient"
// 	return func(next http.Handler) http.Handler {
// 		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 			ctx := context.WithValue(r.Context(), SupabaseClientKey, SupabaseClient)
// 			next.ServeHTTP(w, r.WithContext(ctx))
// 		})
// 	}
// }

// // To test server health
// func handleHomeRoute(w http.ResponseWriter, r *http.Request) {
// 	http.ServeFile(w, r, "static/index.html")
// }

// func main() {

// 	// panic recovery

// 	defer func() {
// 		if r := recover(); r != nil {
// 			log.Println("Error occured: ", r)
// 			debug.PrintStack()
// 		}
// 	}()

// 	_ = godotenv.Load()
// 	var port string = os.Getenv("PORT")

// 	if database.SupabaseClient == nil {
// 		err := errors.New("supabase Client is not initialized")
// 		panic(err)
// 	}

// 	// setup mux router for creating routes
// 	var router *mux.Router = mux.NewRouter()

// 	// Middleware - to pass SupabaseClient to different route handlers
// 	router.Use(SupabaseMiddleware(database.SupabaseClient))

// 	// Serve static files from the /static directory
// 	// router.PathPrefix("/static/").Handler(http.StripPrefix("/static/",
// 	// 	http.FileServer(http.Dir("static"))))

// 	// set routes
// 	// router.HandleFunc("/", handleHomeRoute).Methods("GET")
// 	router.HandleFunc("/api/v1/health", func(w http.ResponseWriter, r *http.Request) {
// 		w.WriteHeader(http.StatusOK)
// 		w.Header().Set("Content-Type", "application/json")
// 		json.NewEncoder(w).Encode(struct {
// 			Code    int
// 			Message string
// 		}{
// 			Code:    http.StatusOK,
// 			Message: "Server is functioning",
// 		})
// 	}).Methods("GET")
// 	router.HandleFunc("/api/v1/shorten", routes.HandleShortIdCreation).Methods("POST")
// 	router.HandleFunc("/{shortid}", routes.HandleUrlRedirection).Methods("GET")
// 	router.HandleFunc("/api/v1/urls", routes.HandleListUrls).Methods("GET")
// 	router.HandleFunc("/api/v1/update/{shortid}", routes.HandleUrlUpdate).Methods("PUT")
// 	router.HandleFunc("/api/v1/delete/{shortid}", routes.HandleUrlDelete).Methods("DELETE")

// 	// Activate server
// 	log.Println("Listening at PORT ", port)
// 	log.Fatal(http.ListenAndServe(":"+port, router))

// }
