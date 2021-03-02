package aws

import (
	"context"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ssm"
)

// SSM extends AWS ssm.SSM
type SSM struct {
	*ssm.SSM
}

// NewSSM creates a new AWS connection returns a Simple Service Manager session
func NewSSM(region string) *SSM {
	sess := session.New()

	return &SSM{ssm.New(sess,
		&aws.Config{
			Region: aws.String(region),
		})}
}

// GetParams returns map of key:value SSM Parameters as listed in paramsToGet along with any error fectching them
func (svc *SSM) GetParams(ctx context.Context, encrpyted bool, root string, paramsToGet []string) (map[string]string, error) {
	params := make(map[string]string, len(paramsToGet))
	var paramsToGetPaths []*string

	// Concat parameter names to SSM path e.g. value to /path/value
	for _, paramToGet := range paramsToGet {
		paramPath := root + paramToGet
		paramsToGetPaths = append(paramsToGetPaths, &paramPath)
	}

	// Get all parameters with single call
	output, err := svc.GetParametersWithContext(ctx,
		&ssm.GetParametersInput{
			Names:          paramsToGetPaths,
			WithDecryption: aws.Bool(encrpyted),
		})

	// Trim parameter paths back to names for map keys e.g. /path/value to value
	for _, param := range output.Parameters {
		key := strings.TrimPrefix(*param.Name, root)
		val := *param.Value
		params[key] = val
	}

	return params, err
}

// PutParam creates a new paramater in the aws parameter store
func (svc *SSM) PutParam(ctx context.Context, encrypted bool, root, key, value string) error {
	var paramType string 
	if encrypted {
		paramType = ssm.ParameterTypeSecureString
	} else {
		paramType = ssm.ParameterTypeString
	}

	name := root + key

	input := &ssm.PutParameterInput{
		Type: &paramType,
		Name: &name,
		Value: &value,
	}

	_, err := svc.PutParameterWithContext(ctx, input)

	return err
}

// DeleteParam removes a paramter from the paramter store
func (svc *SSM) DeleteParam(ctx context.Context, root, key string) error {
	name := root + key

	_, err := svc.DeleteParameterWithContext(ctx, &ssm.DeleteParameterInput{Name: &name})

	return err
}
