package main

import (
	"encoding/json"
	"log"
	"net/http"
	"regexp"
	"strings"

	"github.com/thewolfnl/ibanlib"
)

type serverConfigStruct struct {
	host     string
	port     string
	sanitize bool
}

type responseStruct struct {
	StatusCode       int    `json:"statusCode"`
	IBANCode         string `json:"ibancode,omitempty"`
	IBANPretty       string `json:"ibanpretty,omitempty"`
	ValidationResult bool   `json:"valid"`
	Message          string `json:"message"`
}

type requestStruct struct {
	IBAN    string `json:"iban,omitempty"`
	BBAN    string `json:"bban,omitempty"`
	Country string `json:"country,omitempty"`
}

var serverConfig = serverConfigStruct{
	host:     "localhost",
	port:     "3000",
	sanitize: false,
}

func startServer() {
	log.Println("Start service")
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/validate", validationHandler)
	http.HandleFunc("/bban2iban", bban2ibanHandler)

	log.Printf("%s\n", "Serving on http://"+serverConfig.host+":"+serverConfig.port+"/")
	err := http.ListenAndServe(serverConfig.host+":"+serverConfig.port, nil)
	if err != nil {
		log.Fatal(err)
	}
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	respond(w, http.StatusOK, "To use this service, send a POST request to '/validate' with the IBAN to validate as a 'iban' field value")
}

func validationHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		respond(w, http.StatusMethodNotAllowed, "Method Not Allowed; Request should be a POST request")
		return
	}
	data = nil

	iban, response := extractValue(r, "iban")
	if response != nil {
		returnResponse(w, response)
		return
	}

	message := "Invalid Iban"
	IBAN := ibanlib.ConvertStringToIban(iban)
	if IBAN == nil {
		returnResponse(w, &responseStruct{
			StatusCode:       http.StatusOK,
			ValidationResult: false,
			Message:          message,
		})
		return
	}

	valid := ibanlib.ValidateIban(IBAN)
	if valid {
		message = "Valid Iban"
	}

	returnResponse(w, &responseStruct{
		StatusCode:       http.StatusOK,
		IBANCode:         IBAN.Code(),
		IBANPretty:       IBAN.PrettyCode(),
		ValidationResult: valid,
		Message:          message,
	})
}

func bban2ibanHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		respond(w, http.StatusMethodNotAllowed, "Method Not Allowed; Request should be a POST request")
		return
	}
	data = nil

	bban, response := extractValue(r, "bban")
	if response != nil {
		returnResponse(w, response)
		return
	}
	country, response := extractValue(r, "country")
	if response != nil {
		returnResponse(w, response)
		return
	}

	message := "Invalid Iban"
	IBAN := ibanlib.ConvertBBANStringToIban(bban, country)
	if IBAN == nil {
		returnResponse(w, &responseStruct{
			StatusCode:       http.StatusOK,
			ValidationResult: false,
			Message:          message,
		})
		return
	}

	valid := ibanlib.ValidateIban(IBAN)
	if valid {
		message = "Valid Iban"
	}

	returnResponse(w, &responseStruct{
		StatusCode:       http.StatusOK,
		IBANCode:         IBAN.Code(),
		IBANPretty:       IBAN.PrettyCode(),
		ValidationResult: valid,
		Message:          message,
	})
}

func respond(w http.ResponseWriter, status int, message string) {
	response := responseStruct{
		StatusCode: status,
		Message:    message,
	}
	returnResponse(w, &response)
}

func returnResponse(w http.ResponseWriter, response *responseStruct) {
	JSON, err := json.Marshal(response)
	outputError(err)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(response.StatusCode)
	w.Write(JSON)
}

var data *requestStruct

func extractValue(r *http.Request, key string) (string, *responseStruct) {
	value := ""
	switch r.Header.Get("Content-Type") {
	case "application/json":
		if data == nil {
			decoder := json.NewDecoder(r.Body)
			err := decoder.Decode(&data)
			dieOnError(err)
		}
		switch strings.ToUpper(key) {
		default:
		case "IBAN":
			value = data.IBAN
		case "BBAN":
			value = data.BBAN
		case "COUNTRY":
			value = data.Country
		}
	case "application/x-www-form-urlencoded":
		err := r.ParseForm()
		dieOnError(err)
		value = r.PostForm.Get(key)
	default:
		return "", &responseStruct{
			StatusCode: http.StatusBadRequest,
			Message:    "Invalid Content-Type",
		}
	}
	if value == "" {
		return "", &responseStruct{
			StatusCode: http.StatusBadRequest,
			Message:    "Missing '" + key + "' value",
		}
	}
	if serverConfig.sanitize {
		// Strip unwanted chars
		value = regexp.MustCompile("[^A-Z0-9]+").ReplaceAllString(value, "")
	}
	return value, nil
}
