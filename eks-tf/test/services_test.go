package test

import (
    "fmt"
    "testing"
    "time"
    "net/http"

    "github.com/gruntwork-io/terratest/modules/terraform"
    "github.com/gruntwork-io/terratest/modules/k8s"
    "github.com/stretchr/testify/assert"
)

func TestServicesAreDeployed(t *testing.T) {
    t.Parallel()

    terraformOptions := &terraform.Options{
        TerraformDir: "../infra",
    }

    terraform.InitAndApply(t, terraformOptions)
    defer terraform.Destroy(t, terraformOptions)

    options := k8s.NewKubectlOptions("", "", "default")

    services := []string{
        "api-gateway",
        "backend-app",
        "redis-store",
        "stream-service",
    }

    for _, svc := range services {
        k8s.WaitUntilServiceAvailable(t, options, svc, 10, 5*time.Second)
    }
}

func TestApiGatewayResponse(t *testing.T) {
    options := k8s.NewKubectlOptions("", "", "default")

    // Create port-forward object (8080 on local to 80 on service)
    portForward := k8s.KubectlPortForward{
        KubectlOptions: options,
        ResourceType:   "svc",
        ResourceName:   "api-gateway",
        LocalPort:      8080,
        RemotePort:     80,
    }

    // Start port-forwarding in background
    stopCh := make(chan struct{})
    defer close(stopCh)
    err := k8s.StartKubectlPortForward(t, &portForward, stopCh)
    assert.NoError(t, err)

    url := fmt.Sprintf("http://localhost:%d", portForward.LocalPort)

    // Retry HTTP request until it succeeds or times out
    success := false
    for i := 0; i < 10; i++ {
        resp, err := http.Get(url)
        if err == nil && resp.StatusCode == 200 {
            success = true
            resp.Body.Close()
            break
        }
        time.Sleep(3 * time.Second)
    }

    assert.True(t, success, "Expected successful response from API gateway")
}