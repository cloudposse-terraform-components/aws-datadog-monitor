package test

import (
	"fmt"
	"os"
	"strings"
	"testing"

	helper "github.com/cloudposse/test-helpers/pkg/atmos/component-helper"
	awsTerratest "github.com/gruntwork-io/terratest/modules/aws"
	"github.com/gruntwork-io/terratest/modules/random"
	"github.com/stretchr/testify/assert"
)

type ComponentSuite struct {
	helper.TestSuite

	DDAPIKey  string
	DDHost    string
	RandomID  string
	AWSRegion string
}

func (s *ComponentSuite) TestBasicDatadogMonitor() {
	const component = "datadog-monitor/basic"
	const stack = "default-test"
	const awsRegion = "us-east-2"

	defer s.DestroyAtmosComponent(s.T(), component, stack, nil)
	options, _ := s.DeployAtmosComponent(s.T(), component, stack, nil)
	assert.NotNil(s.T(), options)

	s.DriftTest(component, stack, nil)
}

func (s *ComponentSuite) TestEnabledFlag() {
	const component = "datadog-monitor/disabled"
	const stack = "default-test"
	s.VerifyEnabledFlag(component, stack, nil)
}

func (s *ComponentSuite) SetupSuite() {
	s.InitConfig()
	s.Config.ComponentDestDir = "components/terraform/datadog-monitor"

	// Store the Datadog API key in SSM for the duration of the test.
	// The API key is obtained from GitHub secrets as an environment variable
	// in the test pipelines.
	s.DDAPIKey = os.Getenv("DD_API_KEY")
	if s.DDAPIKey == "" {
		s.T().Fatal("DD_API_KEY environment variable must be set")
	}

	s.RandomID = strings.ToLower(random.UniqueId())
	s.AWSRegion = "us-east-2"
	s.DDHost = "us5.datadoghq.com"

	if !s.Config.SkipDeployDependencies {
		secretPath := fmt.Sprintf("/datadog/%s/api_key", s.RandomID)
		awsTerratest.PutParameter(s.T(), s.AWSRegion, secretPath, "Datadog API Key", s.DDAPIKey)

		inputs := map[string]any{
			"datadog_site_url": s.DDHost,
			"datadog_secrets_source_store_account_region": s.AWSRegion,
			"datadog_secrets_source_store_account_stage":  "default",
			"datadog_secrets_source_store_account_tenant": "test",
			"datadog_api_secret_key":                      s.RandomID,
		}
		s.AddDependency(s.T(), "datadog-configuration", "default-test", &inputs)
	}

	s.TestSuite.SetupSuite()
}

func (s *ComponentSuite) TearDownSuite() {
	s.TestSuite.TearDownSuite()
	if !s.Config.SkipDestroyDependencies {
		secretPath := fmt.Sprintf("/datadog/%s/api_key", s.RandomID)
		awsTerratest.DeleteParameter(s.T(), s.AWSRegion, secretPath)
	}
}

func TestRunSuite(t *testing.T) {
	suite := new(ComponentSuite)

	helper.Run(t, suite)
}
