package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

type RequestBody struct {
	Data string `json:"data"`
}

type Resource struct {
	ResourceType string `json:"resourceType"`
}

func webhookHandler(w http.ResponseWriter, r *http.Request) {
	var (
		// reqBody         RequestBody
		res             Resource
		responseMessage string
		payload         map[string]interface{}
	)

	// Load .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	w.Header().Set("Content-Type", "application/json")
	statusCode := http.StatusOK

	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		printResponse(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	jsonData, _ := json.MarshalIndent(payload, "", "  ")
	fmt.Println("body payload:", string(jsonData))

	if resourceType, ok := payload["resourceType"].(string); ok {
		res.ResourceType = resourceType
	}

	// if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
	// 	printResponse(w, http.StatusBadRequest, "Invalid request body")
	// 	return
	// }
	// // print body data
	// fmt.Println(reqBody.Data)

	// bodyDecoded, _ := base64.StdEncoding.DecodeString(reqBody.Data)

	// decoded, _ := base64.StdEncoding.DecodeString(os.Getenv("PRIVATE_KEY_ORG"))
	// privOrg, _ := ecies.ParseECPrivateKeyPEM(decoded)

	// decoded, _ = base64.StdEncoding.DecodeString(os.Getenv("PUBLIC_KEY_SS"))
	// pubSS, _ := ecies.ParseECPublicKeyPEM(decoded)

	// decrypted, err := ecies.Decrypt(bodyDecoded, privOrg, pubSS)
	// if err != nil {
	// 	printResponse(w, http.StatusBadRequest, err.Error())
	// 	return
	// }
	// err = json.Unmarshal(decrypted, &res)
	// if err != nil {
	// 	printResponse(w, http.StatusBadRequest, err.Error())
	// 	return
	// }

	switch res.ResourceType {
	case "CoverageEligibilityResponse", "ChargeItemResponse", "BillingStatus", "ClaimResponse", "PaymentReconciliation", "PaymentNotice":
		responseMessage = "Faskes Notified Successfully"
	case "CoverageEligibilityRequest", "ChargeItem", "Claim":
		responseMessage = "Insurance Notified Successfully"
	case "Bundle":
		responseMessage = "SISRUTE Notified Sucessfully"
	default:
		statusCode = http.StatusBadRequest
		responseMessage = "Invalid Resource Type"
	}
	fmt.Println(responseMessage)
	printResponse(w, statusCode, responseMessage)
}

func printResponse(w http.ResponseWriter, code int, msg string) {
	w.WriteHeader(code)
	fmt.Fprintf(w, `{"message": "%s"}`, msg)
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("PORT environment variable is required")
	}

	http.HandleFunc("/webhook", webhookHandler)
	log.Printf("Server is running on port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

func printPayload(r *http.Request) {
	body, _ := io.ReadAll(r.Body)
	fmt.Println("body payload:", string(body))
}
