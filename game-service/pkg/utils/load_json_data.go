package utils

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
)

// LoadJSONData reads JSON data from a file and unmarshals it into the target structure
func LoadJSONData(filePath string, target interface{}) error {
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("cannot open file: %w", err)
	}
	defer file.Close()

	byteValue, err := io.ReadAll(file)
	if err != nil {
		return fmt.Errorf("cannot read file: %w", err)
	}

	err = json.Unmarshal(byteValue, target)
	if err != nil {
		return fmt.Errorf("cannot parse json: %w", err)
	}

	return nil
}
