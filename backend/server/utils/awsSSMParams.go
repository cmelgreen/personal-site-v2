package utils

import (
	"context"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ssm"
)

// AwsSSM extends AWS ssm.SSM
type AwsSSM struct {
	*ssm.SSM
}

// NewSSM creates a new AWS connection returns a Simple Service Manager session
func NewSSM(region string) *AwsSSM {
	sess := session.New()

	return &AwsSSM{ssm.New(sess,
		&aws.Config{
			Region: aws.String(region),
		})}

}

// GetParams returns map of key:value SSM Parameters as listed in paramsToGet along with any error fectching them
func (svc *AwsSSM) GetParams(ctx context.Context, encrpyted bool, root string, paramsToGet []string) (map[string]string, error) {
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
