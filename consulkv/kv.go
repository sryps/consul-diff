package consulkv

import (
	"consuldiff/storage"
	"encoding/base64"
	"log"

	"github.com/hashicorp/consul/api"
)

func FetchKV(client *api.Client, prefix string) ([]map[string]string, error) {
	kv := client.KV()
	pairs, _, err := kv.List(prefix, nil)
	if err != nil {
		return nil, err
	}

	result := make(map[string]string)
	for _, pair := range pairs {
		result[pair.Key] = string(pair.Value)
	}
	return []map[string]string{result}, nil
}

func FetchKVBase64(client *api.Client, prefix string) ([]map[string]string, error) {
	kv := client.KV()
	pairs, _, err := kv.List(prefix, nil)
	if err != nil {
		return nil, err
	}

	result := make(map[string]string)
	for _, pair := range pairs {
		if pair.Value != nil {
			result[pair.Key] = base64.StdEncoding.EncodeToString(pair.Value)
		} else {
			result[pair.Key] = ""
		}
	}

	return []map[string]string{result}, nil
}

func LogKVDiff(prev, curr []map[string]string, filepath string) {
	// Detect added or changed keys
	c := curr[0]
	p := prev[0]
	for k, v := range c {
		if oldVal, ok := p[k]; !ok {
			log.Printf("[+] Added: %s = %s", k, v)
		} else if oldVal != v {
			log.Printf("[~] Modified: %s    Old: %s    New: %s", k, oldVal, v)
		}
	}

	// Detect deleted keys
	for k := range p {
		if _, ok := c[k]; !ok {
			log.Printf("[-] Deleted: %s = %s", k, p[k])
		}
	}

	storage.WriteMapToFile(filepath, curr)
}
