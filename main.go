package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"time"
)

// "postgurrll/utils"
func Greet(w http.ResponseWriter, r *http.Request) {

	//major issue in postman
	// we will not be able to hardcode this json
	var data struct {
		Name  string `json:"name"`
		Class int    `json:"class"`
	}

	err := json.NewDecoder(r.Body).Decode(&data)

	if err != nil {
		http.Error(w, "InvalidBodyRequest", http.StatusBadRequest)
		fmt.Fprintf(w, "InvalidBody")
		return
	}

	fmt.Fprintf(w, data.Name)
}

func dataFetch(w http.ResponseWriter, r *http.Request) {

	var data struct {
		URL     string          `json:"url"`
		Method  string          `json:"method"`
		Payload json.RawMessage `json:"payload"`
	}
	fmt.Println("dataFetch called", data.URL)
	err := json.NewDecoder(r.Body).Decode(&data)

	// get post not working in postman not recieving the body in data struct
	// outgoing header req to confirm json?

	if err != nil {
		http.Error(w, "InvalidBodyFormat", http.StatusBadRequest)
		return
	}

	var bodyReader io.Reader

	bodyReader = bytes.NewBuffer(data.Payload)

	res, err := http.NewRequest(data.Method, data.URL, bodyReader)

	if err != nil {
		fmt.Print("error in fetching data ", data.URL, err)
		http.Error(w, "DataNotFetched", http.StatusBadRequest)
		return
	}

	start := time.Now()
	res.Header.Set("Content-Type", "application/json")
	client := http.Client{}

	respo, err := client.Do(res)

	if err != nil || respo.StatusCode >= 400 {
		http.Error(w, "failed to read response ", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	// lets set outgoing header to json for the request we are making to the url

	// r.Header.Set("Content-Type", "application/json")
	fmt.Println(respo.StatusCode)

	bodyBytes, err := io.ReadAll(respo.Body)

	if err != nil {
		http.Error(w, "failed to read response body", http.StatusInternalServerError)
		return
	}

	type response struct {
		ID        string          `json:"id"`
		LatencyMS int64           `json:"latency_ms"`
		Data      json.RawMessage `json:"data"`
	}

	responseBodyFinal := response{
		ID:        "REQ-" + strconv.Itoa(os.Getpid()),
		LatencyMS: time.Since(start).Milliseconds(),
		Data:      json.RawMessage(bodyBytes),
	}

	w.Header().Set("Content-Type", "application/json")

	err1 := json.NewEncoder(w).Encode(responseBodyFinal)

	if err1 != nil {
		http.Error(w, "Encoding failed", http.StatusInternalServerError)
		return
	}

}

func main() {
	fmt.Println("Hello World")

	http.HandleFunc("/greet", Greet)
	http.HandleFunc("/datadedo", dataFetch)

	err := http.ListenAndServe(":8080", nil)

	if err != nil {
		fmt.Printf("error int starting the server ", err)
	}
}
