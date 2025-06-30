package storage

import (
	"encoding/json"
	"errors"
	"os"
)

// WriteMapIfNotExists writes data to filename only if the file doesn't already exist.
func WriteMapToFile(filename string, data map[string]string) error {
	// O_CREATE: create if not exist
	// O_EXCL: fail if file already exists
	// O_WRONLY: write-only
	file, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		if errors.Is(err, os.ErrExist) {
			return errors.New("file already exists")
		}
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	return encoder.Encode(data)
}

// ReadMapFromFile reads a JSON file into a map[string]interface{}.
func ReadMapFromFile(filepath string) (map[string]string, error) {
	file, err := os.Open(filepath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var result map[string]string
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&result)
	if err != nil {
		return nil, err
	}

	return result, nil
}
