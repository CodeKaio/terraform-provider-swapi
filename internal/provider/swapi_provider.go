// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"github.com/codekaio/terraform-provider-swapi/internal/swapi"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"os"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure SWAPIProvider satisfies various provider interfaces.
var _ provider.Provider = &SWAPIProvider{}

// SWAPIProvider defines the provider implementation.
type SWAPIProvider struct {
	// version is set to the provider version on release, "dev" when the
	// provider is built and ran locally, and "test" when running acceptance
	// testing.
	version string
}

// SWAPIProviderModel describes the provider data model.
type SWAPIProviderModel struct {
	Endpoint types.String `tfsdk:"endpoint"`
	ApiKey   types.String `tfsdk:"api_key"`
}

func (p *SWAPIProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "swapi"
	resp.Version = p.version
}

func (p *SWAPIProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"endpoint": schema.StringAttribute{
				MarkdownDescription: "Endpoint to connect SWAPI. Can be set by environement variable SWAPI_ENDPOINT.",
				Optional:            true,
			},
			"api_key": schema.StringAttribute{
				MarkdownDescription: "API Key. Can be set by environement variable SWAPI_APIKEY.",
				Optional:            true,
				Sensitive:           true,
			},
		},
	}
}

func (p *SWAPIProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var config SWAPIProviderModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// default values with env vars
	endpoint := os.Getenv("SWAPI_ENDPOINT")
	api_key := os.Getenv("SWAPI_APIKEY")

	// override env vars with provider attributes
	if !config.Endpoint.IsNull() {
		endpoint = config.Endpoint.ValueString()
	}
	if !config.ApiKey.IsNull() {
		api_key = config.ApiKey.ValueString()
	}

	if endpoint == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("endpoint"),
			"Unknown SWAPI Endpoint",
			"The provider cannot create the SW API client as there is an unknown configuration value for the SWAPI host.",
		)
	}

	if api_key == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("api_key"),
			"Unknown SWAPI API key",
			"The provider cannot create the SW API client as there is an unknown configuration value for the SWAPI key.",
		)
	}

	if resp.Diagnostics.HasError() {
		// no need to continue if we have an error
		return
	}

	// Example client configuration for config sources and resources
	client := swapi.NewSWAPIClient(endpoint, api_key)
	resp.DataSourceData = client
	resp.ResourceData = client
}

func (p *SWAPIProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewPlanetResource,
	}
}

func (p *SWAPIProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewPlanetDataSource,
	}
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &SWAPIProvider{
			version: version,
		}
	}
}
