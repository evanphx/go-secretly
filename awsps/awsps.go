package awsps

import (
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ssm"
	secretly "github.com/evanphx/go-secretly"
)

type ps struct{}

func (g *ps) Get(path string) (string, error) {
	sess, err := session.NewSession()
	if err != nil {
		return "", err
	}

	svc := ssm.New(sess)

	input := &ssm.GetParameterInput{}
	input.Name = &path
	input.SetWithDecryption(true)

	out, err := svc.GetParameter(input)
	if err != nil {
		return "", err
	}

	if out.Parameter == nil {
		return "", secretly.ErrNoSuchSecret
	}

	return *out.Parameter.Value, nil
}

func (g *ps) Put(path, val string) error {
	sess, err := session.NewSession()
	if err != nil {
		return err
	}

	svc := ssm.New(sess)

	input := &ssm.PutParameterInput{}
	input.SetType("SecureString")
	input.Name = &path
	input.Value = &val

	_, err = svc.PutParameter(input)
	return err
}

func init() {
	secretly.AddBackend("awsps", &ps{})
}
