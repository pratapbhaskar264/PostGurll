package main

import (
	"encoding/json"
	"fmt"
	"net/http"
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
		URL string `json:"url"`
	}

	err := json.NewDecoder(r.Body).Decode(&data)

	if err != nil {
		http.Error(w, "InvalidBodyFormat", http.StatusBadRequest)
		return
	}
	// var responseBodyFinal struct {
	// 	StartTime   time.Time `json:"startTime"`
	// 	UserId      int       `json:"userId"`
	// 	Id          int       `json:"id"`
	// 	Title       string    `json:"title"`
	// 	IsCompleted bool      `json:"completed"`
	// 	EndTime     time.Time `json:"endTime"`
	// }
	// responseBodyFinal.StartTime = time.Now()
	res, err := http.Get(data.URL)
	// responseBodyFinal.EndTime = time.Now()

	if err != nil {
		http.Error(w, "DataNotFetched", http.StatusBadRequest)
		return
	}

	defer res.Body.Close()

	if res.StatusCode != 200 {
		http.Error(w, "NotOk", http.StatusBadRequest)
		return
	}

	var responseBody struct {
		UserId      int    `json:"userId"`
		Id          int    `json:"id"`
		Title       string `json:"title"`
		IsCompleted bool   `json:"completed"`
	}
	er := json.NewDecoder(res.Body).Decode(&responseBody)
	// responseBodyFinal.UserId = responseBody.UserId
	// responseBodyFinal.Id = responseBody.Id
	// responseBodyFinal.Title = responseBody.Title
	// responseBodyFinal.IsCompleted = responseBody.IsCompleted

	if er != nil {
		http.Error(w, "DataFormatMismatched", http.StatusBadRequest)
		return
	}

	// fmt.Fprintf(w, responseBody.Title)

	w.Header().Set("Content-Type", "application/json")

	e := json.NewEncoder(w).Encode(responseBody)

	if e != nil {
		http.Error(w, "DataFormatMismatched", http.StatusBadRequest)
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
