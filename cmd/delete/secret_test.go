package delete

import (
	"testing"

	"github.com/mas2020-golang/cryptex/packages/utils"
)

// TestRemoveSecret_SecretExists tests removing a secret that exists in the box
func TestRemoveSecret_SecretExists(t *testing.T) {
	// Setup: Create a box with multiple secrets
	box := &utils.Box{
		Name:    "test-box",
		Version: "1.0.0",
		Secrets: []*utils.Secret{
			{Name: "secret1", Pwd: "password1"},
			{Name: "secret2", Pwd: "password2"},
			{Name: "secret3", Pwd: "password3"},
		},
	}

	initialCount := len(box.Secrets)

	// Execute: Remove the middle secret
	deleted, err := removeSecret("secret2", box)

	// Assert
	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}

	if !deleted {
		t.Error("Expected deleted to be true, got false")
	}

	if len(box.Secrets) != initialCount-1 {
		t.Errorf("Expected %d secrets, got %d", initialCount-1, len(box.Secrets))
	}

	// Verify the correct secret was removed
	for _, s := range box.Secrets {
		if s.Name == "secret2" {
			t.Error("Secret 'secret2' should have been removed but still exists")
		}
	}

	// Verify remaining secrets are intact
	if box.Secrets[0].Name != "secret1" || box.Secrets[1].Name != "secret3" {
		t.Error("Remaining secrets are not in the expected order")
	}
}

// TestRemoveSecret_SecretDoesNotExist tests removing a secret that doesn't exist
func TestRemoveSecret_SecretDoesNotExist(t *testing.T) {
	// Setup: Create a box with secrets
	box := &utils.Box{
		Name:    "test-box",
		Version: "1.0.0",
		Secrets: []*utils.Secret{
			{Name: "secret1", Pwd: "password1"},
			{Name: "secret2", Pwd: "password2"},
		},
	}

	initialCount := len(box.Secrets)

	// Execute: Try to remove a non-existent secret
	deleted, err := removeSecret("nonexistent", box)

	// Assert
	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}

	if deleted {
		t.Error("Expected deleted to be false, got true")
	}

	if len(box.Secrets) != initialCount {
		t.Errorf("Expected %d secrets (unchanged), got %d", initialCount, len(box.Secrets))
	}
}

// TestRemoveSecret_EmptyBox tests removing from a box with an empty secrets slice
func TestRemoveSecret_EmptyBox(t *testing.T) {
	// Setup: Create a box with empty secrets slice
	box := &utils.Box{
		Name:    "test-box",
		Version: "1.0.0",
		Secrets: []*utils.Secret{},
	}

	// Execute: Try to remove a secret
	deleted, err := removeSecret("any-secret", box)

	// Assert
	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}

	if deleted {
		t.Error("Expected deleted to be false, got true")
	}

	if len(box.Secrets) != 0 {
		t.Errorf("Expected 0 secrets, got %d", len(box.Secrets))
	}
}

// TestRemoveSecret_NilSecrets tests removing from a box with nil secrets
func TestRemoveSecret_NilSecrets(t *testing.T) {
	// Setup: Create a box with nil secrets
	box := &utils.Box{
		Name:    "test-box",
		Version: "1.0.0",
		Secrets: nil,
	}

	// Execute: Try to remove a secret
	deleted, err := removeSecret("any-secret", box)

	// Assert
	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}

	if deleted {
		t.Error("Expected deleted to be false, got true")
	}

	if box.Secrets != nil {
		t.Error("Expected secrets to remain nil")
	}
}

// TestRemoveSecret_FirstSecret tests removing the first secret in the list
func TestRemoveSecret_FirstSecret(t *testing.T) {
	// Setup: Create a box with multiple secrets
	box := &utils.Box{
		Name:    "test-box",
		Version: "1.0.0",
		Secrets: []*utils.Secret{
			{Name: "secret1", Pwd: "password1"},
			{Name: "secret2", Pwd: "password2"},
			{Name: "secret3", Pwd: "password3"},
		},
	}

	// Execute: Remove the first secret
	deleted, err := removeSecret("secret1", box)

	// Assert
	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}

	if !deleted {
		t.Error("Expected deleted to be true, got false")
	}

	if len(box.Secrets) != 2 {
		t.Errorf("Expected 2 secrets, got %d", len(box.Secrets))
	}

	// Verify the first secret was removed
	if box.Secrets[0].Name != "secret2" {
		t.Errorf("Expected first secret to be 'secret2', got '%s'", box.Secrets[0].Name)
	}
}

// TestRemoveSecret_LastSecret tests removing the last secret in the list
func TestRemoveSecret_LastSecret(t *testing.T) {
	// Setup: Create a box with multiple secrets
	box := &utils.Box{
		Name:    "test-box",
		Version: "1.0.0",
		Secrets: []*utils.Secret{
			{Name: "secret1", Pwd: "password1"},
			{Name: "secret2", Pwd: "password2"},
			{Name: "secret3", Pwd: "password3"},
		},
	}

	// Execute: Remove the last secret
	deleted, err := removeSecret("secret3", box)

	// Assert
	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}

	if !deleted {
		t.Error("Expected deleted to be true, got false")
	}

	if len(box.Secrets) != 2 {
		t.Errorf("Expected 2 secrets, got %d", len(box.Secrets))
	}

	// Verify the last secret was removed
	if box.Secrets[len(box.Secrets)-1].Name != "secret2" {
		t.Errorf("Expected last secret to be 'secret2', got '%s'", box.Secrets[len(box.Secrets)-1].Name)
	}
}

// TestRemoveSecret_OnlySecret tests removing the only secret in the box
func TestRemoveSecret_OnlySecret(t *testing.T) {
	// Setup: Create a box with only one secret
	box := &utils.Box{
		Name:    "test-box",
		Version: "1.0.0",
		Secrets: []*utils.Secret{
			{Name: "only-secret", Pwd: "password1"},
		},
	}

	// Execute: Remove the only secret
	deleted, err := removeSecret("only-secret", box)

	// Assert
	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}

	if !deleted {
		t.Error("Expected deleted to be true, got false")
	}

	if len(box.Secrets) != 0 {
		t.Errorf("Expected 0 secrets, got %d", len(box.Secrets))
	}
}

// TestRemoveSecret_ComplexSecret tests removing a secret with all fields populated
func TestRemoveSecret_ComplexSecret(t *testing.T) {
	// Setup: Create a box with a complex secret
	box := &utils.Box{
		Name:    "test-box",
		Version: "1.0.0",
		Secrets: []*utils.Secret{
			{
				Name:        "simple-secret",
				Pwd:         "password1",
				Login:       "user1",
				Url:         "https://example1.com",
				Notes:       "Some notes",
				Version:     "1.0.0",
				LastUpdated: "2023-01-01",
			},
			{
				Name:        "complex-secret",
				Pwd:         "password2",
				Login:       "user2",
				Url:         "https://example2.com",
				Notes:       "More notes",
				Version:     "2.0.0",
				LastUpdated: "2023-01-02",
				Others: map[string]string{
					"key1": "value1",
					"key2": "value2",
				},
			},
			{
				Name:        "another-secret",
				Pwd:         "password3",
				Login:       "user3",
				Url:         "https://example3.com",
			},
		},
	}

	// Execute: Remove the complex secret
	deleted, err := removeSecret("complex-secret", box)

	// Assert
	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}

	if !deleted {
		t.Error("Expected deleted to be true, got false")
	}

	if len(box.Secrets) != 2 {
		t.Errorf("Expected 2 secrets, got %d", len(box.Secrets))
	}

	// Verify the complex secret was removed
	for _, s := range box.Secrets {
		if s.Name == "complex-secret" {
			t.Error("Secret 'complex-secret' should have been removed but still exists")
		}
	}

	// Verify remaining secrets
	if box.Secrets[0].Name != "simple-secret" || box.Secrets[1].Name != "another-secret" {
		t.Error("Remaining secrets are not as expected")
	}
}

// TestRemoveSecret_CaseSensitive tests that secret names are case-sensitive
func TestRemoveSecret_CaseSensitive(t *testing.T) {
	// Setup: Create a box with a secret
	box := &utils.Box{
		Name:    "test-box",
		Version: "1.0.0",
		Secrets: []*utils.Secret{
			{Name: "MySecret", Pwd: "password1"},
		},
	}

	// Execute: Try to remove with different case
	deleted, err := removeSecret("mysecret", box)

	// Assert
	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}

	if deleted {
		t.Error("Expected deleted to be false (case-sensitive), got true")
	}

	if len(box.Secrets) != 1 {
		t.Errorf("Expected 1 secret (unchanged), got %d", len(box.Secrets))
	}
}
