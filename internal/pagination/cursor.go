package pagination

import (
	"encoding/base64"
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
)

type NearbyCursor struct {
	DistanceMetres float64   `json:"d"`
	UserID         uuid.UUID `json:"id"`
}

func EncodeCursor(c NearbyCursor) (string, error) {
	raw, err := json.Marshal(c)
	if err != nil {
		return "", fmt.Errorf("encode cursor: %w", err)
	}
	return base64.URLEncoding.EncodeToString(raw), nil
}

func DecodeCursor(s string) (NearbyCursor, error) {
	var c NearbyCursor
	raw, err := base64.URLEncoding.DecodeString(s)
	if err != nil {
		return c, fmt.Errorf("decode cursor: invalid encoding: %w", err)
	}
	if err := json.Unmarshal(raw, &c); err != nil {
		return c, fmt.Errorf("decode cursor: invalid payload: %w", err)
	}
	return c, nil
}