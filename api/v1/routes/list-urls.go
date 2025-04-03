package routes

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"runtime/debug"

	"github.com/harshitrajsinha/go-url-shortener/api/v1/common"
	"github.com/supabase-community/supabase-go"
)

// utils -> Resposne struct, generateShortId func

func HandleListUrls(w http.ResponseWriter, r *http.Request) {

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

	// get all data wrt ip address
	var rowData []map[string]interface{}
	data, _, err := client.From("go_url_shortner").Select("*", "exact", false).Eq("client_ip_addr", clientIpAddr).Execute()
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(data, &rowData)
	if err != nil {
		panic(err)
	}

	if len(rowData) == 0 {
		// No data found
		w.WriteHeader(http.StatusNotFound)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(common.Response{Code: http.StatusNotFound, Message: "No data found", Reference: "Shorten a URL and then try again"})
		log.Println("No data found")
		return
	} else {
		// Retrieve all to URL
		var responseData []map[string]interface{}
		w.WriteHeader(http.StatusOK)
		for _, item := range rowData {
			responseData = append(responseData, map[string]interface{}{item["short_id"].(string): item["original_url"].(string)})
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(common.Response{Code: http.StatusOK, Message: "Data found", Data: responseData})
		log.Println("Data retrieved")
		return
	}
}
