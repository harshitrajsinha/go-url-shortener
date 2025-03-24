package routes

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/supabase-community/supabase-go"
)

// utils -> Resposne struct, generateShortId func


func HandleListUrls(w http.ResponseWriter, r *http.Request) {

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

	// get all data wrt ip address
	var rowData []map[string]interface{}
	data, _, err := client.From("go_url_shortner").Select("*", "exact", false).Eq("client_ip_addr", networkID).Execute()
	if err != nil{
		w.WriteHeader(http.StatusInternalServerError)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(Response{Code: http.StatusInternalServerError, Message: "Internal Server Error"})
		log.Fatal("Error fetching redirect url from database")
		return
	}
	json.Unmarshal(data, &rowData)
	
	if len(rowData) == 0{
		// No data found
		w.WriteHeader(http.StatusNotFound)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(Response{Code: http.StatusNotFound, Message: "No data found", Reference: "Shorten a URL and then try again"})
		log.Println("No data found")
		return
	}else{
		// Retrieve all to URL
		var responseData []map[string]string
		w.WriteHeader(http.StatusOK)
		for _, item := range rowData{
			responseData = append(responseData, map[string]string{item["short_id"].(string): item["redirect_url"].(string)})
		}
		dataToSend, _ := json.Marshal(responseData)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(Response{Code: http.StatusOK, Message: "Data found", Reference: string(dataToSend)})
		log.Println("Data retrieved")
		return
	}
}