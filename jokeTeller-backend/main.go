package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func fetchJoke(category string) (string, error) {
	//Static var declaration; didn't use := because it's a global variable
	var buffer bytes.Buffer
	//Fetch the joke
	url := ("https://v2.jokeapi.dev/joke/" + category)
	resp, err := http.Get(url)
	// defer initiates the close of the response body after the function returns; the close is deferred until the function returns
	defer resp.Body.Close()

	// go fundamentals idiom is to check for errors first and then use the result
	if err != nil {
		return "There was an error fetching the joke", err
	}
	// chunks are the data that is read from the response body, in this case sorted into 1024 bytes and then written to the buffer.
	// the buffer is a temporary storage for the data that is read from the response body.
	chunk := make([]byte, 1024)
	for {
		n, err := resp.Body.Read(chunk)
		// if the error is not the end of file, break the loop
		// write the chunk to the buffer
		if n > 0 {
			buffer.Write(chunk[:n])
		}
		if err == io.EOF {
			break
		}
		// if the error is not nil, return the error
		if err != nil {
			return "", fmt.Errorf("error reading joke: %w", err)
		}

	}

	// convert the buffer to a string
	jokeBytes := buffer.Bytes()
	jokeString := string(jokeBytes)
	return jokeString, nil
}

type Joke struct {
	Setup    string `json:"setup"`
	Delivery string `json:"delivery"`
	Type     string `json:"type"`
	Category string `json:"category"`
	Error    bool   `json:"error"`
	Id       int    `json:"id"`
	Safe     bool   `json:"safe"`
	Lang     string `json:"lang"`
	Flags    Flags  `json:"flags"`
	Joke     string `json:"joke"`
}
type Flags struct {
	Nsfw      bool `json:"nsfw"`
	Religious bool `json:"relgious"`
	Political bool `json:"political"`
	Racist    bool `json:"racist"`
	Sexist    bool `json:"sexist"`
	Explicit  bool `json:"explicit"`
}

func handlerForHomepage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello World, from Go Webapp!")
}
func handleJokeAPI(w http.ResponseWriter, r *http.Request) {

	// This header allows CORS for local development; in production, this should be restricted to the actual domain of the frontend application.
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:4200")
	w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK) // Handle preflight request
		return
	}

	category := r.URL.Query().Get("category")
	if category == "" {
		category = "Any"
	}
	// pulling our fetchJoke function in and defining it as the raw JSON variable
	jokeRawJSON, err := fetchJoke(category)
	if err != nil {
		fmt.Println("Error fetching joke:", err)
		return
	}

	// Translating the raw JSON being passed into programmable pieces.
	var parsedJoke Joke
	unmarshalErr := json.Unmarshal([]byte(jokeRawJSON), &parsedJoke)
	if unmarshalErr != nil {
		http.Error(w, "Error unmarshaling JSON:"+unmarshalErr.Error(), http.StatusInternalServerError)
		return
	}
	jsonResponse, marshalErr := json.Marshal(parsedJoke)
	if marshalErr != nil {
		http.Error(w, "Error marshaling JSON:"+marshalErr.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonResponse)
}

func main() {

	// adding HTTP handling for GUI implementation
	http.HandleFunc("/", handlerForHomepage)
	http.HandleFunc("/api/joke", handleJokeAPI)

	fmt.Println("server starting on port :8080...")
	serverErr := http.ListenAndServe(":8080", nil)
	if serverErr != nil {
		fmt.Println("error starting server:", serverErr)
	}
	return
}
