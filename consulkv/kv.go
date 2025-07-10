package consulkv

import (
	"consuldiff/kvtypes"
	"consuldiff/storage"
	"encoding/base64"
	"log"

	"github.com/hashicorp/consul/api"
)

func FetchKV(client *api.Client, prefix string) ([]kvtypes.KVExportEntry, error) {
	kv := client.KV()
	pairs, _, err := kv.List(prefix, nil)
	if err != nil {
		return nil, err
	}

	var resp string
	result := []kvtypes.KVExportEntry{}
	for _, pair := range pairs {
		if pair.Value != nil {
			resp = string(pair.Value)
		} else {
			resp = ""
		}
		result = append(result, kvtypes.KVExportEntry{
			Key:   pair.Key,
			Flags: pair.Flags,
			Value: resp,
		})
	}
	return result, nil
}

func FetchKVBase64(client *api.Client, prefix string) ([]kvtypes.KVExportEntry, error) {
	kv := client.KV()
	pairs, _, err := kv.List(prefix, nil)
	if err != nil {
		return nil, err
	}

	var result []kvtypes.KVExportEntry
	for _, pair := range pairs {
		encoded := ""
		if pair.Value != nil {
			encoded = base64.StdEncoding.EncodeToString(pair.Value)
		}
		result = append(result, kvtypes.KVExportEntry{
			Key:   pair.Key,
			Flags: pair.Flags,
			Value: encoded,
		})
	}
	return result, nil
}

func LogKVDiff(prev []kvtypes.KVExportEntry, curr []kvtypes.KVExportEntry, filepath string) {
	// Convert slices to maps for easier comparison
	prevMap := make(map[string]string)
	for _, entry := range prev {
		prevMap[entry.Key] = entry.Value
	}

	currMap := make(map[string]string)
	for _, entry := range curr {
		currMap[entry.Key] = entry.Value
	}

	// Detect added and modified keys
	for k, v := range currMap {
		if oldVal, ok := prevMap[k]; !ok {
			log.Printf("[+] Added: %s = %s", k, v)
		} else if oldVal != v {
			log.Printf("[~] Modified: %s    Old: %s    New: %s", k, oldVal, v)
		}
	}

	// Detect deleted keys
	for k, v := range prevMap {
		if _, ok := currMap[k]; !ok {
			log.Printf("[-] Deleted: %s = %s", k, v)
		}
	}

	// Save current state
	storage.WriteMapToFile(filepath, curr)
}
