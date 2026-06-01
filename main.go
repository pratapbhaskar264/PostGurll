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

	err := json.NewDecoder(r.Body).Decode(&data) // r.body is the actual data from client

	if err != nil {
		http.Error(w, "InvalidBodyRequest", http.StatusBadRequest)
		fmt.Fprintf(w, "InvalidBody")
		return
	}

	fmt.Fprintf(w, data.Name)
}

func dataFetch(w http.ResponseWriter, r *http.Request) {

	var data struct {
		URL     string              `json:"url"`
		Method  string              `json:"method"`
		Headers map[string][]string `json:"headers`
		Payload json.RawMessage     `json:"payload"`
	}

	err := json.NewDecoder(r.Body).Decode(&data)

	// get post not working in postman not recieving the body in data struct
	// outgoing header req to confirm json?

	if err != nil {
		http.Error(w, "InvalidBodyFormat", http.StatusBadRequest)
		return
	}

	valid := validateJSON(data.Payload)

	if !valid {
		http.Error(w, "Invalid Payload Format", http.StatusBadRequest)
		return
	}

	if data.Method == "" {
		http.Error(w, "HTTP method is required", http.StatusBadRequest)
		return
	}
	var bodyReader io.Reader

	bodyReader = bytes.NewBuffer(data.Payload)

	fmt.Print(data.Method)

	res, err := http.NewRequest(data.Method, data.URL, bodyReader) // envelop created request to the url with method
	//and body has nothing to do with content type

	if err != nil {
		fmt.Print("error in fetching data ", data.URL, err)
		http.Error(w, "DataNotFetched", http.StatusBadRequest)
		return
	}

	start := time.Now()
	res.Header.Set("Content-Type", "application/json") // headers are set for the request we are making to the url

	// every single header to other machine
	for key, values := range data.Headers {
		for _, val1 := range values {
			res.Header.Add(key, val1)
		}
	}

	// var redirectHops []string
	// Do this right at the start of your function
	redirectHops := make([]string, 0)

	client := http.Client{
		//custom hook that will justify latency due to redirections and also give us the urls of the redirections
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			redirectHops = append(redirectHops, req.URL.String())
			return nil
		},
	}

	respo, err := client.Do(res) // opens a low-level network socket connection to the target server over the internet,
	// streams your headers and payload bytes across the wire, and waits to hand you back the response (respo).

	if err != nil || respo.StatusCode >= 400 {
		http.Error(w, "failed to read response ", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	// lets set outgoing header to json for the request we are making to the url
	//OR
	//"Hey, Postgurrll is about to send you a JSON object containing
	//  the telemetry ID, the latency, and the target's data. Prepare your UI to format it as JSON."

	// r.Header.Set("Content-Type", "application/json")
	fmt.Println(respo.StatusCode)

	bodyBytes, err := io.ReadAll(respo.Body)

	targetHeaders := respo.Header

	if err != nil {
		http.Error(w, "failed to read response body", http.StatusInternalServerError)
		return
	}

	type response struct {
		ID        string              `json:"id"`
		LatencyMS int64               `json:"latency_ms"`
		Hops      []string            `json:"hops"`
		Headers   map[string][]string `json:"headers"`
		Data      json.RawMessage     `json:"data"`
	}

	responseBodyFinal := response{
		ID:        "REQ-" + strconv.Itoa(os.Getpid()),
		LatencyMS: time.Since(start).Milliseconds(),
		Hops:      redirectHops,
		Headers:   targetHeaders,
		Data:      json.RawMessage(bodyBytes),
	}

	w.Header().Set("Content-Type", "application/json")

	err1 := json.NewEncoder(w).Encode(responseBodyFinal)

	if err1 != nil {
		http.Error(w, "Encoding failed", http.StatusInternalServerError)
		return
	}

}

func validateJSON(jsonData []byte) bool { // json validator
	return json.Valid(jsonData)
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
