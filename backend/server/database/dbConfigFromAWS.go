package database

import (
	"context"
	"fmt"

	"github.com/spf13/viper"
)

var ssmParams = []string{
	"database",
	"host",
	"port",
	"user",
	"password",
}

// DBConfigFromAWS implements DBConfig
type DBConfigFromAWS struct {
	BaseAWSRegion  string
	BaseAWSRoot    string
	BaseConfigName string
	BaseConfigPath string
	WithEncrpytion bool
}

// ConfigString returns database connection string based on AWS_ROOT and remote SSM parameters
func (dbConfig DBConfigFromAWS) ConfigString(ctx context.Context) (string, error) {
	err := dbConfig.loadBaseConfigFromDotEnv()
	if err != nil {
		return "", err
	}

	awsRegion := viper.GetString(dbConfig.BaseAWSRegion)
	ssmRoot := viper.GetString(dbConfig.BaseAWSRoot)

	svc := NewSSM(awsRegion)

	params, err := svc.GetParams(ctx, dbConfig.WithEncrpytion, ssmRoot, ssmParams)
	if err != nil {
		return "", err
	}

	configString := fmt.Sprintf(
		"database=%s host=%s port=%s user=%s password=%s",
		params["database"],
		params["host"],
		params["port"],
		params["user"],
		params["password"],
	)

	return configString, nil
}

// Pull AWS_ROOT and AWS_REGION from .env file
func (dbConfig DBConfigFromAWS) loadBaseConfigFromDotEnv() error {
	viper.SetConfigName(dbConfig.BaseConfigName)
	viper.AddConfigPath(dbConfig.BaseConfigPath)

	err := viper.ReadInConfig()
	if err != nil {
		return err
	}

	return nil
}
