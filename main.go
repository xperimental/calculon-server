package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

type Request struct {
	Operation  string `json:"operation"`
	OperandOne int    `json:"operandOne"`
	OperandTwo int    `json:"operandTwo"`
}

type Response struct {
	Error        bool   `json:"isError"`
	ErrorMessage string `json:"errorMessage,omitempty"`
	Result       int    `json:"result"`
}

const (
	OPERATION_MULTIPLY = "MULTIPLY"
	OPERATION_DIVIDE   = "DIVIDE"
)

func calculateHandler(rw http.ResponseWriter, req *http.Request) {
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		log.Fatal(err)
	}
	var request Request
	err = json.Unmarshal(body, &request)
	if err != nil {
		log.Fatal(err)
	}
	var response Response
	switch request.Operation {
	case OPERATION_MULTIPLY:
		response.Result = request.OperandOne * request.OperandTwo
		break
	case OPERATION_DIVIDE:
		if request.OperandTwo == 0 {
			response.Error = true
			response.ErrorMessage = "Division by zero!"
		} else {
			response.Result = request.OperandOne / request.OperandTwo
		}
		break
	default:
		response.Error = true
		response.ErrorMessage = "Invalid operation: " + request.Operation
		break
	}
	log.Printf("Request: %+v Response: %+v", request, response)
	responseData, err := json.Marshal(response)
	if err != nil {
		log.Fatal(err)
	}
	rw.Header().Set("Content-Type", "application/json")
	if response.Error {
		rw.WriteHeader(400)
	}
	_, err = rw.Write(responseData)
	if err != nil {
		log.Fatal(err)
	}
}

func getPort() string {
	port := os.Getenv("PORT")
	if len(port) == 0 {
		port = "8080"
	}
	return port
}

func main() {
	http.Handle("/", http.FileServer(http.Dir("web/")))
	http.HandleFunc("/calc", calculateHandler)

	port := getPort()
	log.Printf("Listening on port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
