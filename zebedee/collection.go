package zebedee

import (
	"encoding/json"
	"os"
	"time"
)

type Collection struct {
	PendingDeletes []PendingDelete `json:"pendingDeletes"`
	PublishEndDate time.Time       `json:"publishEndDate"`
}

type PendingDelete struct {
	Root PendingDeleteRoot `json:"root"`
}

type PendingDeleteRoot struct {
	URI string `json:"uri"`
}

func GetCollectionFromFile(filename string) (*Collection, error) {
	body, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var col Collection

	err = json.Unmarshal(body, &col)
	if err != nil {
		return nil, err
	}
	return &col, nil
}
