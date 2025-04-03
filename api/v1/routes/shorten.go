package routes

import (
	"encoding/json"
	"errors"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"runtime/debug"
	"strings"

	"github.com/asaskevich/govalidator"
	"github.com/harshitrajsinha/go-url-shortener/api/v1/common"
	"github.com/supabase-community/supabase-go"
)

type RequestUrl struct {
	Url string `json:"url"`
}

// Type for inserting data to supabase
type DataToInsert struct {
	RedirectURL string `json:"original_url"`
	ShortID     string `json:"short_id"`
	ClientIP    string `json:"client_ip_addr"`
}

// Function to generate shortid of 8 characters
func generateShortId() string {
	// uuid.New().String[:8]  from "github.com/google/uuid"
	var possibleChars string = "ABCDEFGHIJKLMNOPQQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"
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

	defer func() {
		if r := recover(); r != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(common.Response{Code: http.StatusInternalServerError, Message: "Internal Server Error"})
			log.Println("Error occured: ", r)
			debug.PrintStack()
		}
	}()

	var requestUrl RequestUrl

	// Get request body
	if err := json.NewDecoder(r.Body).Decode(&requestUrl); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(common.Response{Code: http.StatusBadRequest, Message: "Invalid JSON", Reference: "{url:your-url}"})
		return
	}
	if requestUrl.Url == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(common.Response{Code: http.StatusBadRequest, Message: "URL is missing", Reference: "{url:your-url}"})
		log.Println("URL is empty")
		return
	}

	// Check for valid URL
	if !govalidator.IsURL(requestUrl.Url) {
		w.WriteHeader(http.StatusForbidden)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(common.Response{Code: http.StatusForbidden, Message: "URL is not valid", Reference: "{url:https://some-valid-url.domain}"})
		log.Println("Requested url is not valid ", requestUrl.Url)
		return
	}

	// Check if URL contains https
	parsedURL, err := url.Parse(requestUrl.Url)
	if err != nil {
		panic(err)
	}
	if parsedURL.Scheme != "https" {
		w.WriteHeader(http.StatusForbidden)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(common.Response{Code: http.StatusForbidden, Message: "URL is not valid", Reference: "{url:https://some-valid-url.domain}"})
		log.Println("Requested url is not valid ", requestUrl.Url)
		return
	}

	// get client ip address
	clientIpAddr := common.GetClientIpAddr(r)

	if clientIpAddr == "" {

		err := errors.New("ip address not received")
		panic(err)
	}

	// Initialize DB client
	const SupabaseClientKey string = "SupabaseClient"
	client, ok := r.Context().Value(SupabaseClientKey).(*supabase.Client)
	if !ok {
		err := errors.New("error initializing db client for handler")
		panic(err)
	}

	// check if url already exists
	var doesURLExists []map[string]interface{}
	data, _, err := client.From("go_url_shortner").Select("*", "exact", false).Eq("client_ip_addr", clientIpAddr).Eq("original_url", requestUrl.Url).Execute()
	if err != nil {
		panic(err)
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
		err := errors.New("shortid not generated")
		panic(err)
	}

	// Insert into database
	var newData DataToInsert = DataToInsert{RedirectURL: requestUrl.Url, ShortID: shortId, ClientIP: clientIpAddr}
	_, _, err = client.From("go_url_shortner").Insert(newData, false, "", "", "").Execute()
	if err != nil {
		panic(err)
	}

	// send response
	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(common.Response{Code: http.StatusCreated, Message: "Shorten url generated successfully", Data: []map[string]interface{}{{"shortened-url": string(r.Host) + "/" + shortId}}})
	if err != nil {
		panic(err)
	}
}
