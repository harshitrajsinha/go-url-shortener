package routes

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
	"github.com/supabase-community/supabase-go"
)

func HandleUrlDelete(w http.ResponseWriter, r *http.Request) {
	var queryParams map[string]string = mux.Vars(r)
	var shortID string = queryParams["shortid"]

	// Get shortid from path variable
	if shortID == ""{
		w.WriteHeader(http.StatusBadRequest)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(Response{Code: http.StatusBadRequest, Message: "shortID is missing", Reference: "/api/v1/redirect/XXXXXXXX"})
		log.Println("Error shortid is missing")
		return
	}else if len(shortID) != 8{
		w.WriteHeader(http.StatusBadRequest)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(Response{Code: http.StatusBadRequest, Message: "Invalid shortID", Reference: "/api/v1/redirect/XXXXXXXX"})
		log.Println("insufficient shortID length")
	}

	// get client ip address 
	clientIpAddr := getClientIpAddr(r)
	var networkID string

	// if ipv6
	if strings.Contains(clientIpAddr, "]"){
		indexVal := strings.Index(clientIpAddr, "]:")
		if indexVal != -1{
			clientIpAddr = clientIpAddr[:indexVal+1]
		}
		networkID = convertIPv6LoopbackToIPv4(clientIpAddr)
	}else{
		networkID = strings.Split(clientIpAddr, ":")[0]
	// 	octet := strings.Split(networkID, ".")
	// 	_ = octet[len(octet)-1]
	}

	
	if clientIpAddr == ""{
		w.WriteHeader(http.StatusInternalServerError)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(Response{Code: http.StatusInternalServerError, Message: "Error fetching request"})
		log.Println("IP address not received")
		return
	}

	// Instantiate db client
	const SupabaseClientKey string = "SupabaseClient"
	client, ok := r.Context().Value(SupabaseClientKey).(*supabase.Client)
	if !ok {
		w.WriteHeader(http.StatusInternalServerError)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(Response{Code: http.StatusInternalServerError, Message: "Internal Server Error"})
		log.Fatal("Error initializing client for handler - shortid creation")
		return
	}

	// delete url
	var rowData []map[string]interface{}
	data, _, err := client.From("go_url_shortner").Delete("", "").Eq("client_ip_addr", networkID).Eq("short_id", shortID).Execute()
	if err != nil{
		w.WriteHeader(http.StatusInternalServerError)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(Response{Code: http.StatusInternalServerError, Message: "Internal Server Error"})
		log.Fatal("Error deleting redirect url from database")
		return
	}
	json.Unmarshal(data, &rowData)

	if len(rowData) == 0{
		// Redirect URL not found
		w.WriteHeader(http.StatusNotFound)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(Response{Code: http.StatusNotFound, Message: "Redirect URL not found for given shortID", Reference: "Either not found or already deleted"})
		log.Println("Redirect url not found or already deleted")
		return
	}else{
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(Response{Code: http.StatusOK, Message: "Redirect URL deleted successfully"})
		log.Println("Redirect url not found or already deleted")
		return
	}
}