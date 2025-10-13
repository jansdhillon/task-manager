package secrets

import (
	"context"
	"fmt"
	"testing"
)

type FakeGoogleSecretsClient struct {
	Secrets map[string]map[int]string
	Err     error
}

func (f *FakeGoogleSecretsClient) GetSecretVersion(ctx context.Context, name string, version int) (string, error) {
	if f.Err != nil {
		return "", f.Err
	}
	if versions, ok := f.Secrets[name]; ok {
		if val, ok := versions[version]; ok {
			return val, nil
		}
		return "", fmt.Errorf("version %d not found for secret %q", version, name)
	}
	return "", fmt.Errorf("secret %q not found", name)
}

func (f *FakeGoogleSecretsClient) GetLatestVersion(ctx context.Context, name string) (string, error) {
	if f.Err != nil {
		return "", f.Err
	}
	versions, ok := f.Secrets[name]
	if !ok || len(versions) == 0 {
		return "", fmt.Errorf("no versions found for secret %q", name)
	}
	var latest int
	for v := range versions {
		if v > latest {
			latest = v
		}
	}
	return versions[latest], nil
}

func GetSecretValue(ctx context.Context, client SecretsClient, name string, versionNumber int) (any, error) {
	return client.GetSecretVersion(ctx, name, versionNumber)
}

func TestGetSecretVersion(t *testing.T) {
	testCases := []struct {
		desc          string
		fakeClient    *FakeGoogleSecretsClient
		secretName    string
		version       int
		expectedValue string
		expectError   bool
	}{
		{
			desc: "should return specific secret version",
			fakeClient: &FakeGoogleSecretsClient{
				Secrets: map[string]map[int]string{"my-secret": {1: "foo", 2: "bar"}},
			},
			secretName:    "my-secret",
			version:       2,
			expectedValue: "bar",
		},
		{
			desc: "should fail on unknown version",
			fakeClient: &FakeGoogleSecretsClient{
				Secrets: map[string]map[int]string{"my-secret": {1: "foo"}},
			},
			secretName:  "my-secret",
			version:     2,
			expectError: true,
		},
		{
			desc: "should fail on missing secret",
			fakeClient: &FakeGoogleSecretsClient{
				Secrets: map[string]map[int]string{},
			},
			secretName:  "unknown",
			version:     1,
			expectError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			got, err := GetSecretValue(context.Background(), tc.fakeClient, tc.secretName, tc.version)
			if (err != nil) != tc.expectError {
				t.Fatalf("unexpected error state: %v", err)
			}
			if got != tc.expectedValue {
				t.Errorf("got %v, want %v", got, tc.expectedValue)
			}
		})
	}
}

func TestGetLatestVersion(t *testing.T) {
	testCases := []struct {
		desc          string
		fakeClient    *FakeGoogleSecretsClient
		secretName    string
		expectedValue string
		expectError   bool
	}{
		{
			desc: "should return the latest secret version",
			fakeClient: &FakeGoogleSecretsClient{
				Secrets: map[string]map[int]string{"my-secret": {1: "foo", 2: "bar"}},
			},
			secretName:    "my-secret",
			expectedValue: "bar",
		},
		{
			desc: "should fail on missing secret",
			fakeClient: &FakeGoogleSecretsClient{
				Secrets: map[string]map[int]string{},
			},
			secretName:  "unknown",
			expectError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			got, err := tc.fakeClient.GetLatestVersion(context.Background(), tc.secretName)
			if (err != nil) != tc.expectError {
				t.Fatalf("unexpected error state: %v", err)
			}
			if got != tc.expectedValue {
				t.Errorf("got %v, want %v", got, tc.expectedValue)
			}
		})
	}
}
