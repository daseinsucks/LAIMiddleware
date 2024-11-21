package registration

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

type RequestData struct {
	Wallet   string `json:"wallet"`
	APIKey   string `json:"apikey"`
	Telegram string `json:"telegram"`
}

type ResponseData struct {
	Reply string `json:"reply"`
}

func Start() {
	_ = godotenv.Load()
	port := os.Getenv("PORT")
	http.HandleFunc("/register", handler)
	fmt.Println("Starting server on " + port)
	if err := http.ListenAndServe(port, nil); err != nil {
		fmt.Printf("Failed to start server: %s\n", err)
	}
}

func handler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodOptions {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		w.WriteHeader(http.StatusOK)
		return
	}

	if r.Method != http.MethodPost {
		http.Error(w, "Only POST method is allowed", http.StatusMethodNotAllowed)
		return
	}
	//TODO: do we filter requests to some origin or are we an open book?
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")

	var reqData RequestData
	if err := json.NewDecoder(r.Body).Decode(&reqData); err != nil {
		http.Error(w, fmt.Sprintf(`{"error": "Invalid request body: %s"}`, err), http.StatusBadRequest)
		return
	}

	fmt.Printf("Received: %+v\n", reqData)

	response := ResponseData{
		Reply: "Registration complete.",
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
	if err := AddToDatabase(reqData.Wallet, reqData.APIKey, reqData.Telegram); err != nil {
		fmt.Printf("Failed to add to database: %s\n", err)
	}

}

func AddToDatabase(address, password, telegram string) error {

	connStr := os.Getenv("PG_LINK")
	if connStr == "" {
		return fmt.Errorf("PGLink isn't in .env file")
	}

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return fmt.Errorf("Error opening database: %v", err)
	}
	defer db.Close()

	//TODO: I'm not sure if this thing here is valid.
	query := `INSERT INTO users (address, password, telegram) VALUES ($1, $2, $3)`
	_, err = db.Exec(query, address, password, telegram)
	if err != nil {
		return fmt.Errorf("Error inserting data into database: %v", err)
	}

	return nil
}
