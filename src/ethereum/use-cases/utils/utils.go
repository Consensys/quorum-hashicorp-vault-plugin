package utils

import "fmt"

func ComputeKey(address, namespace string) string {
	path := fmt.Sprintf("ethereum/accounts/%s", address)
	if namespace != "" {
		path = fmt.Sprintf("%s/%s", namespace, path)
	}

	return path
}
