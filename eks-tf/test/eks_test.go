package test

import (
    "testing"

    "github.com/gruntwork-io/terratest/modules/terraform"
    "github.com/stretchr/testify/assert"
)

func TestMockEKSModule(t *testing.T) {
    t.Parallel()

    tfOptions := &terraform.Options{
        TerraformDir: "../",
        VarFiles:     []string{"./mock.tfvars"},
        NoColor:      true,
    }

    // Init, apply, and defer destroy (safe because it's a mock)
    terraform.InitAndApply(t, tfOptions)
    defer terraform.Destroy(t, tfOptions)

    // Extract outputs from actual state
    clusterName := terraform.Output(t, tfOptions, "mock_cluster_name")
    nodeGroupName := terraform.Output(t, tfOptions, "mock_node_group")

    // Assertions
    assert.Equal(t, "test-mock-cluster", clusterName)
    assert.Equal(t, "test-node-group", nodeGroupName)
}
