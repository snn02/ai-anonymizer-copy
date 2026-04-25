package secrets

import (
	"errors"
	"strings"

	"github.com/zalando/go-keyring"
)

const defaultService = "ai-anonymizer/api"

var errSecretNotFound = keyring.ErrNotFound

type keyringClient interface {
	Get(service, user string) (string, error)
}

type keyringAdapter struct{}

func (keyringAdapter) Get(service, user string) (string, error) {
	return keyring.Get(service, user)
}

type OSSecretChecker struct {
	service string
	client  keyringClient
}

func NewOSSecretChecker() OSSecretChecker {
	return OSSecretChecker{
		service: defaultService,
		client:  keyringAdapter{},
	}
}

func (c OSSecretChecker) HasSecret(partnerID string) (bool, error) {
	secret, err := c.GetSecret(partnerID)
	if err != nil {
		if errors.Is(err, errSecretNotFound) {
			return false, nil
		}
		return false, err
	}
	return strings.TrimSpace(secret) != "", nil
}

func (c OSSecretChecker) GetSecret(partnerID string) (string, error) {
	if strings.TrimSpace(partnerID) == "" {
		return "", errors.New("partner id is required")
	}
	secret, err := c.client.Get(c.service, partnerID)
	if err != nil {
		return "", err
	}
	return secret, nil
}
