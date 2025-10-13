package db

import (
	"context"
	"fmt"
	"log"
	"testing"

	"github.com/jansdhillon/task-manager/server/internal/config"
)

type FakeSecretsClient struct {
	// name: { version*: val}
	Secrets map[string]map[int]string
}

func (fsc *FakeSecretsClient) GetSecretVersion(ctx context.Context, secretName string, secretVersion int) (string, error) {
	if _, ok := fsc.Secrets[secretName]; !ok {
		return "", fmt.Errorf("secret not found!")
	}

	scv, ok := fsc.Secrets[secretName][secretVersion]

	if !ok {
		return "", fmt.Errorf("secret version not found!")
	}

	return scv, nil
}

func (fsc *FakeSecretsClient) GetLatestVersion(ctx context.Context, secretName string) (string, error) {
	if _, ok := fsc.Secrets[secretName]; !ok {
		return "", fmt.Errorf("secret not found!")
	}

	max := 1
	for _, vMap := range fsc.Secrets {
		for ver := range vMap {
			if ver > max {
				max = ver
			}
		}

		return vMap[max], nil
	}

	return "", nil

}

func TestGetDsn(t *testing.T) {
	testCases := []struct {
		desc string
		dsn  string
	}{
		{
			desc: "getting the dsn works",
			dsn:  "postgresql://hello@world:5432",
		},
	}
	for _, tC := range testCases {
		sc := &FakeSecretsClient{
			Secrets: map[string]map[int]string{
				config.POSTGRES_DSN_ENV: {
					1: tC.dsn,
				},
			},
		}
		t.Run(tC.desc, func(t *testing.T) {
			dsn, err := GetDsn(context.Background(), sc)

			if err != nil {
				t.Fatalf("error getting dsn: %v", err)
			}

			if dsn != tC.dsn {
				t.Fatalf("want: %s, got: %s", tC.dsn, dsn)
			}

			log.Printf("expected dsn: %s, got dsn: %s", tC.dsn, dsn)
		})
	}
}
