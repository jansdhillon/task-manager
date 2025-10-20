package db

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jansdhillon/task-manager/server/internal/config"
	"github.com/jansdhillon/task-manager/server/internal/secrets"
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

func TestConnect(t *testing.T) {
	const (
		projectID   = "test-project"
		expectedDSN = "postgresql://user:pass@localhost:5432/db"
	)

	t.Setenv(config.GCP_PROJECT_ID_ENV, projectID)

	fakeSecrets := &FakeSecretsClient{
		Secrets: map[string]map[int]string{
			config.POSTGRES_DSN_ENV: {
				1: expectedDSN,
			},
		},
	}

	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("sqlmock.New(): %v", err)
	}

	mock.ExpectPing()

	originalSecretsClient := newSecretsClient
	originalSQLOpen := sqlOpenFunc

	newSecretsClient = func(ctx context.Context, gotProjectID string) (secrets.SecretsClient, error) {
		if gotProjectID != projectID {
			return nil, fmt.Errorf("unexpected project ID: got %s", gotProjectID)
		}
		return fakeSecrets, nil
	}

	sqlOpenFunc = func(driverName, dsn string) (*sql.DB, error) {
		if driverName != "pgx" {
			return nil, fmt.Errorf("unexpected driver name: %s", driverName)
		}
		if dsn != expectedDSN {
			return nil, fmt.Errorf("unexpected DSN: %s", dsn)
		}
		return mockDB, nil
	}

	t.Cleanup(func() {
		newSecretsClient = originalSecretsClient
		sqlOpenFunc = originalSQLOpen
		_ = mockDB.Close()
	})

	db, err := Connect()
	if err != nil {
		t.Fatalf("Connect() error = %v", err)
	}

	if db != mockDB {
		t.Fatalf("Connect() returned unexpected DB instance")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet expectations: %v", err)
	}
}
