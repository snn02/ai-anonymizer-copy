package secrets

import (
	"errors"
	"regexp"
	"strings"

	"github.com/zalando/go-keyring"
)

const defaultService = "ai-anonymizer/api"

var errSecretNotFound = keyring.ErrNotFound
var ErrSecretNotFound = errors.New("secret not found")
var uuidLikePattern = regexp.MustCompile(`(?i)^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)

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
		if errors.Is(err, ErrSecretNotFound) {
			return false, nil
		}
		return false, err
	}
	return strings.TrimSpace(secret) != "", nil
}

func (c OSSecretChecker) GetSecret(partnerID string) (string, error) {
	normalizedPartnerID := strings.TrimSpace(partnerID)
	if normalizedPartnerID == "" {
		return "", errors.New("partner id is required")
	}
	if uuidLikePattern.MatchString(normalizedPartnerID) {
		normalizedPartnerID = strings.ToLower(normalizedPartnerID)
	}
	secret, err := c.client.Get(c.service, normalizedPartnerID)
	if err != nil {
		if errors.Is(err, errSecretNotFound) {
			return "", ErrSecretNotFound
		}
		return "", err
	}
	return secret, nil
}
