package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func getBody(r *http.Request) interface{} {
	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		log.Error().Err(err).Msg("Error reading request body")
		return map[string]string{}
	}
	defer r.Body.Close()

	if len(bodyBytes) == 0 {
		return nil
	}

	var payload interface{}
	if err := json.Unmarshal(bodyBytes, &payload); err != nil {
		log.Error().Err(err).Msg("Error parsing JSON payload")
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
	var logURL = flag.Bool("url", false, "Log the URL of the request, default: false")
	var logHeaders = flag.Bool("headers", false, "Log the headers of the request, default: false")
	var logBody = flag.Bool("body", false, "Log the headers of the request, default: false")
	var logMethod = flag.Bool("method", false, "Log the method of the request, default: false")
	flag.Parse()

	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	if !*logURL && !*logHeaders && !*logBody && !*logMethod {
		zerolog.SetGlobalLevel(zerolog.FatalLevel)
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
			log.Error().Err(err).Msg("Error encoding JSON")
			return
		}
		log.Info().RawJSON("Request", jsonData).Msg("")
		fmt.Fprint(w, "Ok")
	})
	log.Info().Msgf("Starting server on :%d", *port)
	log.Fatal().Err(http.ListenAndServe(fmt.Sprintf(":%d", *port), nil)).Msg("")
}
