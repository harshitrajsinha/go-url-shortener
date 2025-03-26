package routes

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/asaskevich/govalidator"
	"github.com/gorilla/mux"
	"github.com/supabase-community/supabase-go"
)

type UpdateRequest struct{
	Url string `json:"url"`
}

func HandleUrlUpdate(w http.ResponseWriter, r *http.Request) {
	var updateRequest UpdateRequest

	var queryParams map[string]string = mux.Vars(r)
	var shortID string = queryParams["shortid"]

	// Get shortid from path variable
	if shortID == ""{
		w.WriteHeader(http.StatusBadRequest)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(Response{Code: http.StatusBadRequest, Message: "shortID is missing", Reference: "/api/v1/udpate/XXXXXXXX"})
		log.Println("Error shortid is missing")
		return
	}else if len(shortID) != 8{
		w.WriteHeader(http.StatusBadRequest)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(Response{Code: http.StatusBadRequest, Message: "Invalid shortID", Reference: "/api/v1/udpate/XXXXXXXX"})
		log.Println("insufficient shortID length")
		return
	}

	// Get request body
	if err := json.NewDecoder(r.Body).Decode(&updateRequest); err != nil{
		w.WriteHeader(http.StatusBadRequest)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(Response{Code: http.StatusBadRequest, Message: "Invalid JSON", Reference: "{short-id: shortid, url:url-to-update}"})
		return
	}
	if updateRequest.Url == ""{
		w.WriteHeader(http.StatusBadRequest)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(Response{Code: http.StatusBadRequest, Message: "URL is missing", Reference: "{short-id: shortid, url:url-to-update}"})
		log.Println("URL is empty")
		return
	}

	// Check for valid URL
	if !govalidator.IsURL(updateRequest.Url){
		w.WriteHeader(http.StatusForbidden)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(Response{Code: http.StatusForbidden, Message: "URL is not valid", Reference: "{url:https://some-valid-url.domain}"})
		log.Println("Requested url is not valid ", updateRequest.Url)
		return
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
		log.Fatal("Error initializing client for handler")
		return
	}

	// get data to update
	var rowData []map[string]interface{}
	data, _, err := client.From("go_url_shortner").Select("*", "exact", false).Eq("client_ip_addr", networkID).Eq("short_id", shortID).Execute()
	if err != nil{
		w.WriteHeader(http.StatusInternalServerError)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(Response{Code: http.StatusInternalServerError, Message: "Internal Server Error"})
		log.Fatal("Error fetching redirect url from database")
		return
	}
	json.Unmarshal(data, &rowData)
	log.Println(rowData)
	if len(rowData) == 0{
		// Data not found
		w.WriteHeader(http.StatusNotFound)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(Response{Code: http.StatusNotFound, Message: "Data not found", Reference: "Check shortid and try again"})
		log.Println("PUT/Data not found")
		return
	}else{
		// update data
		newData := map[string]interface{}{"redirect_url": updateRequest.Url}
		result, _,  err:= client.From("go_url_shortner").Update(newData, "", "exact").Eq("client_ip_addr", networkID).Eq("short_id", shortID).Execute()
		if err != nil{
			w.WriteHeader(http.StatusInternalServerError)
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(Response{Code: http.StatusInternalServerError, Message: "Internal Server Error", Reference: "Error while updating data"})
			log.Println("PUT/Error while updating data")
			return
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(Response{Code: http.StatusOK, Message: "Data updated successfully", Reference: string(result)})
		log.Println("Data updated successfully")
		return
	}
}