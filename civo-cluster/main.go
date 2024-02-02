package main

import (
	"context"
	"strings"
	"time"
)

type CivoCluster struct{}

// example usage: "dagger call cluster-list --api-key <api-key>"
func (m *CivoCluster) ClusterList(ctx context.Context,
	// apiKey API key used to against the Civo API. Found at https://dashboard.civo.com/account/api
	apiToken *Secret,
	// region The region to list clusters in
	region string,
) (string, error) {
	c, err := civoContainer(ctx, apiToken)
	if err != nil {
		panic(err)
	}

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
	c, err := civoContainer(ctx, apiToken)
	if err != nil {
		return "", err
	}

	return c.
		WithEnvVariable("CACHE_BUSTER", time.Now().String()).
		WithExec([]string{"k3s", "get", name, "--region", region}).
		Stdout(ctx)
}

func civoContainer(ctx context.Context, apiToken *Secret) (*Container, error) {
	platform, err := dag.DefaultPlatform(ctx)
	if err != nil {
		return nil, err
	}
	platfromSplit := strings.SplitN(string(platform), "/", 2)

	return dag.Container().
		From("alpine:latest").
		WithExec([]string{"apk", "add", "curl"}).
		WithExec([]string{"curl", "-L", "-o", "/tmp/civo.tar.gz", "https://github.com/civo/cli/releases/download/v1.0.73/civo-1.0.73-" + platfromSplit[0] + "-" + platfromSplit[1] + ".tar.gz"}).
		WithExec([]string{"tar", "-xvf", "/tmp/civo.tar.gz", "-C", "/tmp"}).
		WithExec([]string{"mv", "/tmp/civo", "/usr/local/bin/civo"}).
		WithExec([]string{"chmod", "+x", "/usr/local/bin/civo"}).
		WithSecretVariable("CIVO_TOKEN", apiToken).
		WithEntrypoint([]string{"civo"}), nil

}
