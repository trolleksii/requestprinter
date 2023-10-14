package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

func getBody(r *http.Request) interface{} {
	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		log.Printf("Error reading request body: %v", err)
		return map[string]string{}
	}
	defer r.Body.Close()

	if len(bodyBytes) == 0 {
		return nil
	}

	var payload interface{}
	if err := json.Unmarshal(bodyBytes, &payload); err != nil {
		log.Printf("Error parsing JSON payload: %v", err)
		payload = string(bodyBytes)
	}
	return payload
}

func getHeaders(r *http.Request) map[string]string {
	headers := make(map[string]string)
	for k, v := range r.Header {
		for _, h := range v {
			headers[k] = h
		}
	}
	return headers
}

func main() {
	var port = flag.Int("port", 8080, "Server port, default: 8080")
	var logTime = flag.Bool("time", false, "Log the time of the request, default: false")
	var logURL = flag.Bool("url", false, "Log the URL of the request, default: false")
	var logHeaders = flag.Bool("headers", false, "Log the headers of the request, default: false")
	var logBody = flag.Bool("body", false, "Log the headers of the request, default: false")
	var logMethod = flag.Bool("method", false, "Log the method of the request, default: false")
	flag.Parse()

	if !*logTime {
		log.SetFlags(0)
	}

	if !*logURL && !*logHeaders && !*logBody && !*logMethod {
		log.SetOutput(io.Discard)
	} else {
		log.SetOutput(os.Stdout)
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		response := make(map[string]interface{})
		if *logURL {
			response["URL"] = r.URL.String()
		}
		if *logHeaders {
			response["Headers"] = getHeaders(r)
		}
		if *logBody {
			response["Body"] = getBody(r)
		}
		if *logMethod {
			response["Method"] = r.Method
		}

		jsonData, err := json.Marshal(response)
		if err != nil {
			log.Printf("Error encoding JSON: %v", err)
			return
		}
		log.Println(string(jsonData))
		fmt.Fprint(w, "Ok")
	})
	log.Printf("Starting server on :%d", *port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", *port), nil))
}
