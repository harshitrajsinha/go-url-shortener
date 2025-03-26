package routes

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/harshitrajsinha/go-url-shortener/api/v1/common"
	"github.com/supabase-community/supabase-go"
)

// utils -> Resposne struct, generateShortId func


func HandleUrlRedirection(w http.ResponseWriter, r *http.Request) {

	var queryParams map[string]string = mux.Vars(r)
	var shortID string = queryParams["shortid"]

	// Get shortid from path variable
	if shortID == ""{
		w.WriteHeader(http.StatusBadRequest)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(common.Response{Code: http.StatusBadRequest, Message: "shortID is missing", Reference: "/api/v1/redirect/XXXXXXXX"})
		log.Println("Error shortid is missing")
		return
	}else if len(shortID) != 8{
		w.WriteHeader(http.StatusBadRequest)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(common.Response{Code: http.StatusBadRequest, Message: "Invalid shortID", Reference: "/api/v1/redirect/XXXXXXXX"})
		log.Println("insufficient shortID length")
		return
	}

	// get client ip address 
	clientIpAddr := common.GetClientIpAddr(r)
	
	if clientIpAddr == ""{
		w.WriteHeader(http.StatusInternalServerError)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(common.Response{Code: http.StatusInternalServerError, Message: "Error fetching request"})
		log.Println("IP address not received")
		return
	}

	
	if clientIpAddr == ""{
		w.WriteHeader(http.StatusInternalServerError)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(common.Response{Code: http.StatusInternalServerError, Message: "Error fetching request"})
		log.Println("IP address not received")
		return
	}

	// Instantiate db client
	const SupabaseClientKey string = "SupabaseClient"
	client, ok := r.Context().Value(SupabaseClientKey).(*supabase.Client)
	if !ok {
		w.WriteHeader(http.StatusInternalServerError)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(common.Response{Code: http.StatusInternalServerError, Message: "Internal Server Error"})
		log.Fatal("Error initializing client for handler")
		return
	}

	// get redirect url
	var rowData []map[string]interface{}
	data, _, err := client.From("go_url_shortner").Select("*", "exact", false).Eq("client_ip_addr", clientIpAddr).Eq("short_id", shortID).Execute()
	if err != nil{
		w.WriteHeader(http.StatusInternalServerError)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(common.Response{Code: http.StatusInternalServerError, Message: "Internal Server Error"})
		log.Fatal("Error fetching redirect url from database")
		return
	}
	json.Unmarshal(data, &rowData)
	
	if len(rowData) == 0{
		// Redirect URL not found
		w.WriteHeader(http.StatusNotFound)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(common.Response{Code: http.StatusNotFound, Message: "Redirect URL not found", Reference: "Length of shortid must be 8"})
		log.Println("Redirect url not found")
		return
	}else{
		// Redirect to URL
		w.Header().Set("Location", rowData[0]["redirect_url"].(string))
		w.WriteHeader(http.StatusPermanentRedirect)
		http.Redirect(w, r, rowData[0]["redirect_url"].(string), http.StatusPermanentRedirect)
		return
	}
}