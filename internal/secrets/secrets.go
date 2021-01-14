package secrets

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/secretsmanager"
	"github.com/aws/aws-sdk-go/service/secretsmanager/secretsmanageriface"
)

// Cache update secrets from secret manager
type Cache struct {
	ssmsvc secretsmanageriface.SecretsManagerAPI
}

func NewCache(awscfg *aws.Config) *Cache {
	sess := session.Must(session.NewSession(awscfg))

	return &Cache{
		ssmsvc: secretsmanager.New(sess),
	}
}

func (sc *Cache) GetValue(key string) (string, error) {
	val, err := sc.ssmsvc.GetSecretValue(&secretsmanager.GetSecretValueInput{SecretId: aws.String(key)})
	if err != nil {
		return "", err
	}

	return aws.StringValue(val.SecretString), nil
}
