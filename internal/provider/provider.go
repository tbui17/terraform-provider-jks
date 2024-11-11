// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

var _ provider.Provider = &KeystoreProvider{}

type KeystoreProvider struct {
	// version is set to the provider version on release, "dev" when the
	// provider is built and ran locally, and "test" when running acceptance
	// testing.
	version string
}

// KeystoreProviderModel describes the provider data model.
type KeystoreProviderModel struct {
}

func (p *KeystoreProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "jks"
	resp.Version = p.version
}

func (p *KeystoreProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "JKS provider.",
		Attributes:  map[string]schema.Attribute{},
	}
}

func (p *KeystoreProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var data KeystoreProviderModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Configuration values are now available.
	// if data.Endpoint.IsNull() { /* ... */ }

}

func (p *KeystoreProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewKeystoreResource,
	}
}

func (p *KeystoreProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewKeystoreDataSource,
	}
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &KeystoreProvider{
			version: version,
		}
	}
}
