package database

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/supabase-community/supabase-go"
)

// Declare client from supabase db
var SupabaseClient *supabase.Client
var err error

func init() {

	//Load supabase variables from env file
	err = godotenv.Load()
	if err != nil{
		log.Fatal("Error setting up env vars - ", err)
	}
	supabase_url := os.Getenv("SUPABASE_URL")
	supabase_key := os.Getenv("SUPABASE_KEY")
	
	log.Println("Initializing Supabase Client...")

	// Initialize supabase client
	SupabaseClient, err = supabase.NewClient(supabase_url, supabase_key, &supabase.ClientOptions{})
	if err != nil {
		log.Fatal("Cannot initalize supabase client - ", err)
	}
	log.Println("Supabase Client initialized")
}
