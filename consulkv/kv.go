package consulkv

import (
	"consuldiff/storage"
	"log"

	"github.com/hashicorp/consul/api"
)

func FetchKV(client *api.Client, prefix string) (map[string]string, error) {
	kv := client.KV()
	pairs, _, err := kv.List(prefix, nil)
	if err != nil {
		return nil, err
	}

	result := make(map[string]string)
	for _, pair := range pairs {
		result[pair.Key] = string(pair.Value)
	}
	return result, nil
}

func DiffKV(prev, curr map[string]string, filepath string) {
	// Detect added or changed keys
	for k, v := range curr {
		if oldVal, ok := prev[k]; !ok {
			log.Printf("[+] Added: %s = %s", k, v)
		} else if oldVal != v {
			log.Printf("[~] Modified: %s    Old: %s    New: %s", k, oldVal, v)
		}
	}

	// Detect deleted keys
	for k := range prev {
		if _, ok := curr[k]; !ok {
			log.Printf("[-] Deleted: %s = %s", k, prev[k])
		}
	}

	storage.WriteMapToFile(filepath, curr)
}
