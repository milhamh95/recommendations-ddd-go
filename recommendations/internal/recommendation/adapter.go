package recommendation

import (
	"errors"
	"net/http"
)

type PartnershipAdaptor struct {
	client *http.Client
	url    string
}

func NewPartnershipAdaptor(client *http.Client, url string) (*PartnershipAdaptor, error) {
	if client == nil {
		return nil, errors.New("client can't be nil")
	}

	if url == "" {
		return nil, errors.New("url can't be empty")
	}

	return &PartnershipAdaptor{client: client, url: url}, nil
}
