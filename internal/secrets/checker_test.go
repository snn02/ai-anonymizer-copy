package secrets

import (
	"errors"
	"testing"
)

type fakeKeyringClient struct {
	value string
	err   error
	user  string
}

func (f fakeKeyringClient) Get(service, user string) (string, error) {
	if f.user != "" && user != f.user {
		return "", errors.New("unexpected user")
	}
	return f.value, f.err
}

func TestOSSecretCheckerHasSecretReturnsTrueWhenCredentialExists(t *testing.T) {
	checker := OSSecretChecker{
		service: "ai-anonymizer/api",
		client: fakeKeyringClient{
			value: "token",
		},
	}

	ok, err := checker.HasSecret("partner-1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !ok {
		t.Fatal("expected secret to be found")
	}
}

func TestOSSecretCheckerHasSecretReturnsFalseWhenNotFound(t *testing.T) {
	checker := OSSecretChecker{
		service: "ai-anonymizer/api",
		client: fakeKeyringClient{
			err: errSecretNotFound,
		},
	}

	ok, err := checker.HasSecret("partner-1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if ok {
		t.Fatal("expected secret to be missing")
	}
}

func TestOSSecretCheckerHasSecretReturnsErrorOnStoreFailure(t *testing.T) {
	checker := OSSecretChecker{
		service: "ai-anonymizer/api",
		client: fakeKeyringClient{
			err: errors.New("credential service unavailable"),
		},
	}

	ok, err := checker.HasSecret("partner-1")
	if err == nil {
		t.Fatal("expected store error")
	}
	if ok {
		t.Fatal("expected ok=false when store returns error")
	}
}

func TestOSSecretCheckerGetSecretReturnsValue(t *testing.T) {
	checker := OSSecretChecker{
		service: "ai-anonymizer/api",
		client: fakeKeyringClient{
			value: "token",
		},
	}

	got, err := checker.GetSecret("partner-1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != "token" {
		t.Fatalf("unexpected token %q", got)
	}
}

func TestOSSecretCheckerGetSecretReturnsNotFoundSentinel(t *testing.T) {
	checker := OSSecretChecker{
		service: "ai-anonymizer/api",
		client: fakeKeyringClient{
			err: errSecretNotFound,
		},
	}

	_, err := checker.GetSecret("partner-1")
	if !errors.Is(err, ErrSecretNotFound) {
		t.Fatalf("expected ErrSecretNotFound, got %v", err)
	}
}

func TestOSSecretCheckerGetSecretTrimsPartnerIDBeforeLookup(t *testing.T) {
	checker := OSSecretChecker{
		service: "ai-anonymizer/api",
		client: fakeKeyringClient{
			value: "token",
			user:  "partner-1",
		},
	}

	got, err := checker.GetSecret("  partner-1  ")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != "token" {
		t.Fatalf("unexpected token %q", got)
	}
}

func TestOSSecretCheckerGetSecretNormalizesUUIDToLowercase(t *testing.T) {
	checker := OSSecretChecker{
		service: "ai-anonymizer/api",
		client: fakeKeyringClient{
			value: "token",
			user:  "da315a9a-b39c-11ef-a107-4cd98f59d07b",
		},
	}

	got, err := checker.GetSecret("DA315A9A-B39C-11EF-A107-4CD98F59D07B")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != "token" {
		t.Fatalf("unexpected token %q", got)
	}
}
