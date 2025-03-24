package routes

import (
	"encoding/json"
	"log"
	"math/rand"
	"net/http"
	"strings"

	"github.com/asaskevich/govalidator"
	"github.com/supabase-community/supabase-go"
)

type RequestUrl struct{
	Url string `json:"url"`
}

type Response struct{
	Code int `json:"code"`
	Message string `json:"message"`
	Reference string `json:"Reference,omitempty"`
}

// Type for inserting data to supabase
type DataToInsert struct{
	RedirectURL string `json:"redirect_url"`
	ShortID string	`json:"short_id"`
	ClientIP string `json:"client_ip_addr"`
}

// Function to generate shortid of 8 characters
func generateShortId() string{
	// uuid.New().String[:8]  from "github.com/google/uuid"
	var possibleChars string= "ABCDEFGHIJKLMNOPQQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"
	const SHORT_ID_LEN = 8
	var shortId strings.Builder

	for range SHORT_ID_LEN {
		var randomIndex int = rand.Intn(len(possibleChars))
		shortId.WriteString(string(possibleChars[randomIndex]))
	}
	return shortId.String()
}

// Function to fetch client's ip address
func getClientIpAddr(r *http.Request) string{
	var clientIpAddr string
	clientIpAddr = r.Header.Get("X-Forwarded-For")
	if clientIpAddr != "" {
		// The client IP is the first one in the list
		ips := strings.Split(clientIpAddr, ",")
		if len(ips) > 0 {
			return strings.TrimSpace(ips[0]) // Return the first IP
		}
	}else if clientIpAddr = r.Header.Get("X-Real-IP"); clientIpAddr != ""{
		// Fallback to X-Real-IP
		return clientIpAddr
	}else{
		// Fallback to remoteAddr
		return r.RemoteAddr
	}
	return ""	
}

// Main handler
func HandleShortIdCreation(w http.ResponseWriter, r *http.Request) {
	var requestUrl RequestUrl
	const SupabaseClientKey string = "SupabaseClient"
	client, ok := r.Context().Value(SupabaseClientKey).(*supabase.Client)
	if !ok {
		w.WriteHeader(http.StatusInternalServerError)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(Response{Code: http.StatusInternalServerError, Message: "Internal Server Error"})
		log.Fatal("Error initializing client for handler - shortid creation")
		return
	}

	// Get url
	if err := json.NewDecoder(r.Body).Decode(&requestUrl); err != nil{
		w.WriteHeader(http.StatusInternalServerError)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(Response{Code: http.StatusBadRequest, Message: "Invalid JSON", Reference: "{url:your-url}"})
		return
	}
	if requestUrl.Url == ""{
		w.WriteHeader(http.StatusInternalServerError)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(Response{Code: http.StatusBadRequest, Message: "URL is missing", Reference: "{url:your-url}"})
		log.Fatal("URL is empty")
		return
	}

	// Check for valid URL
	if !govalidator.IsURL(requestUrl.Url){
		w.WriteHeader(http.StatusForbidden)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(Response{Code: http.StatusForbidden, Message: "URL is not valid", Reference: "{url:https://some-valid-url.domain}"})
		log.Fatal("Requested url is not valid ", requestUrl.Url)
		return
	}

	// get client ip address 
	clientIpAddr := getClientIpAddr(r)
	
	if clientIpAddr == ""{
		w.WriteHeader(http.StatusInternalServerError)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(Response{Code: http.StatusInternalServerError, Message: "Error fetching request"})
		log.Fatal("IP address not received")
		return
	}

	// check if url already exists
	var doesURLExists []map[string]interface{}
	data, _, err := client.From("go_url_shortner").Select("*", "exact", false).Eq("client_ip_addr", clientIpAddr).Eq("redirect_url", requestUrl.Url).Execute()
	if err != nil{
		w.WriteHeader(http.StatusInternalServerError)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(Response{Code: http.StatusInternalServerError, Message: "Internal Server Error"})
		log.Fatal("Error checking existing url in database")
		return
	}
	json.Unmarshal(data, &doesURLExists)
	if len(doesURLExists) != 0 {
		// If data already exists
		w.WriteHeader(http.StatusConflict)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(Response{Code: http.StatusConflict, Message: "URL already shorten", Reference: string(data)})
		log.Fatal("URL already exists in database")
		return
	}	
	
	// generate shortid
	var shortId string = generateShortId()
	if shortId == "" {
		w.WriteHeader(http.StatusInternalServerError)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(Response{Code: http.StatusInternalServerError, Message: "Error shortening URL, try after sometime"})
		log.Fatal("Shortid not generated")
		return
	}
	
	// Insert into database
	var newData DataToInsert = DataToInsert{RedirectURL: requestUrl.Url, ShortID: shortId, ClientIP: clientIpAddr}
	_, _, err = client.From("go_url_shortner").Insert(newData, false, "", "", "").Execute()
	if err != nil{
		w.WriteHeader(http.StatusInternalServerError)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(Response{Code: http.StatusInternalServerError, Message: "Error shortening URL, try after sometime"})
		log.Fatal("Error while inserting data to supabase db")
		return
	}

	// send response
	err = json.NewEncoder(w).Encode(Response{Code: http.StatusOK, Message: "Shorten url generated successfully", Reference: string(r.Host) + "/" + shortId})
	if err != nil{
		w.WriteHeader(http.StatusInternalServerError)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(Response{Code: http.StatusInternalServerError, Message: "Failed to generate short url, try after sometime"})
		log.Fatal("Error while sending response")
		return
	}	
}

/*
check if url's host is of same server then prevent
csrf prevention
*/