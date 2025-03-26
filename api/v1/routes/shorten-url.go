package routes

import (
	"encoding/json"
	"log"
	"math/rand"
	"net/http"
	"strings"

	"github.com/asaskevich/govalidator"
	"github.com/harshitrajsinha/go-url-shortener/api/v1/common"
	"github.com/supabase-community/supabase-go"
)

type RequestUrl struct{
	Url string `json:"url"`
}

// Type for inserting data to supabase
type DataToInsert struct{
	RedirectURL string `json:"redirect_url"`
	ShortID string	`json:"short_id"`
	ClientIP string `json:"client_ip_addr"`
}

type contextKey string

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

// Main handler
func HandleShortIdCreation(w http.ResponseWriter, r *http.Request) {
	var requestUrl RequestUrl

	// Get request body
	if err := json.NewDecoder(r.Body).Decode(&requestUrl); err != nil{
		w.WriteHeader(http.StatusBadRequest)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(common.Response{Code: http.StatusBadRequest, Message: "Invalid JSON", Reference: "{url:your-url}"})
		return
	}
	if requestUrl.Url == ""{
		w.WriteHeader(http.StatusBadRequest)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(common.Response{Code: http.StatusBadRequest, Message: "URL is missing", Reference: "{url:your-url}"})
		log.Println("URL is empty")
		return
	}

	// Check for valid URL
	if !govalidator.IsURL(requestUrl.Url){
		w.WriteHeader(http.StatusForbidden)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(common.Response{Code: http.StatusForbidden, Message: "URL is not valid", Reference: "{url:https://some-valid-url.domain}"})
		log.Println("Requested url is not valid ", requestUrl.Url)
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

	// Initialize DB client
	const SupabaseClientKey string = "SupabaseClient"
	client, ok := r.Context().Value(SupabaseClientKey).(*supabase.Client)
	if !ok {
		w.WriteHeader(http.StatusInternalServerError)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(common.Response{Code: http.StatusInternalServerError, Message: "Internal Server Error"})
		log.Fatal("Error initializing client for handler")
		return
	}

	// check if url already exists
	var doesURLExists []map[string]interface{}
	data, _, err := client.From("go_url_shortner").Select("*", "exact", false).Eq("client_ip_addr", clientIpAddr).Eq("redirect_url", requestUrl.Url).Execute()
	if err != nil{
		w.WriteHeader(http.StatusInternalServerError)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(common.Response{Code: http.StatusInternalServerError, Message: "Internal Server Error"})
		log.Fatal("Error checking existing url in database")
		return
	}
	json.Unmarshal(data, &doesURLExists)
	if len(doesURLExists) != 0 {
		// If data already exists
		w.WriteHeader(http.StatusConflict)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(common.Response{Code: http.StatusConflict, Message: "URL already shorten", Data: doesURLExists})
		log.Println("URL already exists in database")
		return
	}	
	
	// generate shortid
	var shortId string = generateShortId()
	if shortId == "" {
		w.WriteHeader(http.StatusInternalServerError)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(common.Response{Code: http.StatusInternalServerError, Message: "Error shortening URL, try after sometime"})
		log.Println("Shortid not generated")
		return
	}
	
	// Insert into database
	var newData DataToInsert = DataToInsert{RedirectURL: requestUrl.Url, ShortID: shortId, ClientIP: clientIpAddr}
	_, _, err = client.From("go_url_shortner").Insert(newData, false, "", "", "").Execute()
	if err != nil{
		w.WriteHeader(http.StatusInternalServerError)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(common.Response{Code: http.StatusInternalServerError, Message: "Error shortening URL, try after sometime"})
		log.Fatal("Error while inserting data to supabase db")
		return
	}

	// send response
	err = json.NewEncoder(w).Encode(common.Response{Code: http.StatusOK, Message: "Shorten url generated successfully", Data: map[string]string{"shortened-url": string(r.Host) + "/api/v1/redirect/" + shortId}})
	if err != nil{
		w.WriteHeader(http.StatusInternalServerError)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(common.Response{Code: http.StatusInternalServerError, Message: "Failed to generate short url, try after sometime"})
		log.Fatal("Error while sending response")
		return
	}	
}

/*
https not present
check if url's host is of same server then prevent
csrf prevention
*/