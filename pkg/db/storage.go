package db

import "fmt"

type StorageImplementation struct{}

func (s *StorageImplementation) Message() {
	fmt.Printf("storage implementation")
}
