package helpers

import (
	// "encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

func GetAvailableProducts() {
	url := "https://endoflife.date/api/all.json"

	req, err := http.NewRequest("GET", url, nil)

	if err != nil {
		log.Fatal(err)
	}

	req.Header.Add("Accept", "application/json")

	client := &http.Client{}

	resp, err := client.Do(req)

	if err != nil {
		log.Fatal(err)
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(string(body))
}

func GetProduct(product string, version string) []byte {
	url := fmt.Sprintf("https://endoflife.date/api/%s.json", product)

	if version != "" {
		url = fmt.Sprintf("https://endoflife.date/api/%s/%s.json", product, version)
	}

	req, err := http.NewRequest("GET", url, nil)

	if err != nil {
		log.Fatalf("Failed to send request to the API: %v", err)
	}

	req.Header.Add("Accept", "application/json")

	res, err := http.DefaultClient.Do(req)

	if err != nil {
		log.Fatalf("Failed to fetch data from the API: %v", err)
	}

	defer res.Body.Close()
	body, _ := io.ReadAll(res.Body)

	// if top && version != "" {
	// 	log.Fatal("Cant top for specific version")
	// } else {
	// 	if len(body) > 3 {
	// 		body = body[:3] // Slice to get the top 3 items
	// 	}
	// }

	return body
}

func ExportToFile(outputData []byte, outputFolder string) {

	err := os.MkdirAll(outputFolder, os.ModePerm)

	if err != nil {
		log.Fatalf("Failed to create a folder: %v", err)
	}

	filePath := outputFolder + "/output.json"

	body := outputData

	file, err := os.Create(filePath)

	if err != nil {
		log.Fatalf("Failed to create file: %v", err)
	}

	defer file.Close()

	_, err = fmt.Fprintln(file, string(body))

	if err != nil {
		log.Fatalf("Error writing to the file: %v", err)
	}

	fmt.Println("Content written to file successfully.")
}
