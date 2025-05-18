package utils

import (
	"encoding/json"
	"log"
	"os"
)

func WriteJSON(filename string, data any) {
	content, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		log.Fatalf("Failed to marshal JSON for %s: %v\n", filename, err)
	}
	if err := os.WriteFile(filename, content, 0644); err != nil {
		log.Fatalf("Failed to write file %s: %v\n", filename, err)
	}
}
