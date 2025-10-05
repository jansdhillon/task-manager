package secrets

import (
	"context"
	"fmt"
	"log"

	secretmanager "cloud.google.com/go/secretmanager/apiv1"
	"cloud.google.com/go/secretmanager/apiv1/secretmanagerpb"
)

type SecretsClient interface {
	GetSecretVersion(ctx context.Context, secretName string, versionNumber int) (any, error)
	//CreateSecretVersion(ctx context.Context, secretName string, value any) (any, error)
	GetLatestVersion(ctx context.Context, secretName string) (any, error)
}

type GcpSecretsClient struct {
	projectId string
	client    *secretmanager.Client
}

func NewGcpSecretsClient(ctx context.Context, projectId string) (*GcpSecretsClient, error) {
	gcpClient, err := secretmanager.NewClient(ctx)
	if err != nil {
		return nil, err
	}
	return &GcpSecretsClient{
		projectId: projectId,
		client:    gcpClient,
	}, nil
}

func (c *GcpSecretsClient) Close() error {
	return c.client.Close()
}

func (c *GcpSecretsClient) GetSecretVersion(ctx context.Context, secretName string, versionNumber int) (any, error) {
	defer c.Close()
	name := fmt.Sprintf("projects/%s/secrets/%s/versions/%d", c.projectId, secretName, versionNumber)
	accessRequest := &secretmanagerpb.AccessSecretVersionRequest{
		Name: name,
	}

	res, err := c.client.AccessSecretVersion(ctx, accessRequest)
	if err != nil {
		log.Fatalf("failed to access secret version: %v", err)
	}

	return string(res.Payload.Data), nil
}

func (c *GcpSecretsClient) GetLatestVersion(ctx context.Context, secretName string) (any, error) {
	defer c.Close()
	name := fmt.Sprintf("projects/%s/secrets/%s/versions/latest", c.projectId, secretName)
	accessRequest := &secretmanagerpb.AccessSecretVersionRequest{
		Name: name,
	}

	res, err := c.client.AccessSecretVersion(ctx, accessRequest)
	if err != nil {
		log.Fatalf("failed to access secret version: %v", err)
	}

	return string(res.Payload.Data), nil
}
