package utils

import "fmt"

func ComputeKey(address, namespace string) string {
	return fmt.Sprintf("%s%s", namespace, address)
}
