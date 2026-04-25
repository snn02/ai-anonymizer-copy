package secrets

import (
	"errors"
	"testing"
)

type fakeKeyringClient struct {
	value string
	err   error
}

func (f fakeKeyringClient) Get(service, user string) (string, error) {
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
