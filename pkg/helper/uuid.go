package helper

import "github.com/google/uuid"

func GenerateUUIDString() string {
	return uuid.New().String()
}

func SetUUIDIfEmpty(id *string) {
	if *id == "" {
		*id = GenerateUUIDString()
	}
}
