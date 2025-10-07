package test

import (
	"testing"

	helper "github.com/cloudposse/test-helpers/pkg/atmos/component-helper"
	"github.com/stretchr/testify/assert"
)

type ComponentSuite struct {
	helper.TestSuite
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

func TestRunSuite(t *testing.T) {
	suite := new(ComponentSuite)

	suite.AddDependency(t, "datadog-configuration", "default-test", nil)

	helper.Run(t, suite)
}
