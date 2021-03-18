package aws

import (
	"errors"
	"fmt"
	"log"
	"os"
	"reflect"

	"github.com/adaptavist/terraform-wrapper/v1/pkg/terraform"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sts"
)

type ConfigAWS struct {
	RoleARN      string   `mapstructure:"role_name" yaml:"role_name"`           // RoleName specifies the role we'll hand to terraform to access its target account
	TryRoleNames []string `mapstructure:"try_role_names" yaml:"try_role_names"` // TryRoleNames is a list of roles we try to assume if RoleName is not provided
	AccountIDVar string   `mapstructure:"account_id_var" yaml:"account_id_var"` // Which env var is used to when using roles in RoleName and TryRoleNames
	RoleTFVar    string   `mapstructure:"role_tf_var" yaml:"role_tf_var"`       // What var are we setting if we find a role
}

func (c ConfigAWS) IsEmpty() bool {
	return reflect.DeepEqual(c, ConfigAWS{})
}

// Detect for an available role and set it as a variable ready for a Terraform provider.
func (c ConfigAWS) Configure(opts *terraform.Opts) (err error) {
	if c.RoleTFVar == "" {
		return errors.New("we don't know what terraform variable we need to set for the role ARN as aws.tf_role_arn is empty")
	}

	// Roles been pinned, that's good enough for us
	if c.RoleARN != "" {
		opts.Variables[c.RoleTFVar] = c.RoleARN
		return
	}

	// Lets try our role names
	if len(c.TryRoleNames) > 0 {
		if c.AccountIDVar == "" {
			return errors.New("aws.try_role_names requires aws.account_id_var to be aws")
		}

		accountID, ok := os.LookupEnv(c.AccountIDVar)

		if !ok || accountID == "" {
			return fmt.Errorf("%s not found or empty", c.AccountIDVar)
		}

		roleArn, err := FindRole(accountID, c.TryRoleNames)

		if err != nil {
			return err
		}

		opts.Variables[c.RoleTFVar] = roleArn
	}
	return
}

// MakeRoleArn generates a role ARN from an account ID and role name
func MakeRoleArn(accountID, roleName string) string {
	return fmt.Sprintf("arn:aws:iam::%s:role/%s", accountID, roleName)
}

// FindRole detects if a role has been set or iterates over TryRoleNames to test for a valid one.
func FindRole(accountID string, roleNames []string) (roleArn string, err error) {
	// We only create a session at this point
	sessionName := "testing"
	sess := session.Must(session.NewSession())
	svc := sts.New(sess)
	for _, roleName := range roleNames {
		testArn := MakeRoleArn(accountID, roleName)
		log.Printf("trying %s", testArn)
		_, err = svc.AssumeRole(&sts.AssumeRoleInput{
			RoleArn:         &testArn,
			RoleSessionName: &sessionName,
		})

		if err != nil {
			fmt.Println(err.Error())
			continue
		} else {
			roleArn = testArn
			break
		}
	}

	if roleArn == "" {
		err = errors.New("failed to find suitable role")
	}

	return
}
