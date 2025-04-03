package routes

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"runtime/debug"

	"github.com/gorilla/mux"
	"github.com/harshitrajsinha/go-url-shortener/api/v1/common"
	"github.com/supabase-community/supabase-go"
)

func HandleUrlDelete(w http.ResponseWriter, r *http.Request) {

	// panic recovery
	defer func() {
		if r := recover(); r != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(common.Response{Code: http.StatusInternalServerError, Message: "Internal Server Error"})
			log.Println("Error occured: ", r)
			debug.PrintStack()
		}
	}()

	var queryParams map[string]string = mux.Vars(r)
	var shortID string = queryParams["shortid"]

	// Get shortid from path variable
	if shortID == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(common.Response{Code: http.StatusBadRequest, Message: "shortID is missing", Reference: "/api/v1/redirect/XXXXXXXX"})
		log.Println("Error shortid is missing")
		return
	} else if len(shortID) != 8 {
		w.WriteHeader(http.StatusBadRequest)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(common.Response{Code: http.StatusBadRequest, Message: "Invalid shortID", Reference: "/api/v1/redirect/XXXXXXXX"})
		log.Println("insufficient shortID length")
	}

	// get client ip address
	clientIpAddr := common.GetClientIpAddr(r)

	if clientIpAddr == "" {
		err := errors.New("ip address not received")
		panic(err)
	}

	// Instantiate db client
	const SupabaseClientKey string = "SupabaseClient"

	client, ok := r.Context().Value(SupabaseClientKey).(*supabase.Client)
	if !ok {
		err := errors.New("error initializing client for handler")
		panic(err)
	}

	// delete url
	var rowData []map[string]interface{}
	data, _, err := client.From("go_url_shortner").Delete("", "").Eq("client_ip_addr", clientIpAddr).Eq("short_id", shortID).Execute()
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(data, &rowData)
	if err != nil {
		panic(err)
	}

	if len(rowData) == 0 {
		// Redirect URL not found
		w.WriteHeader(http.StatusNotFound)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(common.Response{Code: http.StatusNotFound, Message: "Redirect URL not found for given shortID", Reference: "Either not found or already deleted"})
		log.Println("Redirect url not found or already deleted")
		return
	} else {
		w.WriteHeader(http.StatusNoContent)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(common.Response{Code: http.StatusNoContent})
		log.Println("Data deleted")
		return
	}
}
