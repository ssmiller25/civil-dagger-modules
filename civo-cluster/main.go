package main

import (
	"context"
	"strings"
	"time"
)

const (
	civoVersion = "1.0.75"
)

type CivoCluster struct{}

// example usage: "dagger call cluster-list --api-key <api-key> --region <region>"
func (m *CivoCluster) ClusterList(ctx context.Context,
	// apiKey API key used to against the Civo API. Found at https://dashboard.civo.com/account/api
	apiToken *Secret,
	// region The region to list clusters in
	region string,
) (string, error) {
	c := civoContainer(apiToken)
	return c.
		// with cache buster of time.now
		WithEnvVariable("CACHE_BUSTER", time.Now().String()).
		WithExec([]string{"k3s", "list", "--region", region}).
		Stdout(ctx)
}

func (m *CivoCluster) ClusterShow(ctx context.Context,
	apiToken *Secret,
	region string,
	name string,
) (string, error) {
	c := civoContainer(apiToken)
	return c.
		WithEnvVariable("CACHE_BUSTER", time.Now().String()).
		WithExec([]string{"k3s", "get", name, "--region", region}).
		Stdout(ctx)
}

func (m *CivoCluster) Version(ctx context.Context) (string, error) {
	c := civoContainer(nil)
	return c.
		WithExec([]string{"version"}).
		Stdout(ctx)
}

func civoContainer(apiToken *Secret) *Container {
	ctx := context.Background()
	platform, err := dag.DefaultPlatform(ctx)
	if err != nil {
		panic(err)
	}
	platformSplit := strings.SplitN(string(platform), "/", 2)

	container := dag.Container().
		From("alpine:latest").
		WithExec([]string{"apk", "add", "curl"}).
		WithExec([]string{"curl", "-L", "-o", "/tmp/civo.tar.gz", "https://github.com/civo/cli/releases/download/v" + civoVersion + "/civo-" + civoVersion + "-" + platformSplit[0] + "-" + platformSplit[1] + ".tar.gz"}).
		WithExec([]string{"tar", "-xvf", "/tmp/civo.tar.gz", "-C", "/tmp"}).
		WithExec([]string{"mv", "/tmp/civo", "/usr/local/bin/civo"}).
		WithExec([]string{"chmod", "+x", "/usr/local/bin/civo"}).
		WithEntrypoint([]string{"civo"})

	if apiToken != nil {
		container = container.WithSecretVariable("CIVO_TOKEN", apiToken)
	}
	return container
}
