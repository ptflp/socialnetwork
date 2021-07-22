package utils

import (
	"fmt"
	"math/rand"
	"strings"
	"time"

	"github.com/google/uuid"
)

func ProjectUUIDGen(key string) (string, error) {
	rand.Seed(time.Now().UnixNano())
	id := rand.Intn(89) + 10

	uUUID, err := uuid.NewUUID()
	if err != nil {
		return "", err
	}

	return strings.Join([]string{uUUID.String(), fmt.Sprintf("-%s%d", key, id)}, ""), nil
}
