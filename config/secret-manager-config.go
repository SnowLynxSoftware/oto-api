package config

import (
	"context"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
)

type SecretManagerConfig struct {
	cloudEnv string
	client   *secretsmanager.Client
}

func NewSecretManagerConfig(cloudEnv string) *SecretManagerConfig {
	config, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Fatal(err)
	}

	client := secretsmanager.NewFromConfig(config)

	return &SecretManagerConfig{
		cloudEnv: cloudEnv,
		client:   client,
	}
}

func (s *SecretManagerConfig) GetDebugMode() (bool, error) {
	input := &secretsmanager.GetSecretValueInput{
		SecretId: aws.String(s.cloudEnv + "_DEBUG_MODE"),
	}

	result, err := s.client.GetSecretValue(context.TODO(), input)
	if err != nil {
		return false, err
	}

	secretString := *result.SecretString
	if secretString == "true" {
		return true, nil
	} else {
		return false, nil
	}
}

func (s *SecretManagerConfig) GetDBConnectionString() (string, error) {
	input := &secretsmanager.GetSecretValueInput{
		SecretId: aws.String(s.cloudEnv + "_DB_CONNECTION_STRING"),
	}

	result, err := s.client.GetSecretValue(context.TODO(), input)
	if err != nil {
		return "", err
	}

	secretString := *result.SecretString
	return secretString, nil
}
