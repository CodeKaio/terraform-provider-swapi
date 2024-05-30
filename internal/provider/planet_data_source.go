// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"fmt"
	"github.com/codekaio/terraform-provider-swapi/internal/swapi"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ datasource.DataSource = &PlanetDataSource{}

func NewPlanetDataSource() datasource.DataSource {
	return &PlanetDataSource{}
}

// PlanetDataSource defines the data source implementation.
type PlanetDataSource struct {
	client *swapi.SWAPIClient
}

// PlanetDataSourceModel describes the data source data model.
type PlanetDataSourceModel struct {
	Id         types.String `tfsdk:"id"`
	Name       types.String `tfsdk:"name"`
	Diameter   types.Int64  `tfsdk:"diameter"`
	Population types.Int64  `tfsdk:"population"`
}

func (d *PlanetDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_planet"
}

func (d *PlanetDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Planet data source",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "Planet Id",
				Computed:            true,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "Planet Name",
				Required:            true,
			},
			"diameter": schema.Int64Attribute{
				MarkdownDescription: "Planet diameter",
				Computed:            true,
			},
			"population": schema.Int64Attribute{
				MarkdownDescription: "Planet population",
				Computed:            true,
			},
		},
	}
}

func (d *PlanetDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*swapi.SWAPIClient)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *SWAPIClient, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	d.client = client
}

func (d *PlanetDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state PlanetDataSourceModel

	// Read Terraform configuration state into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	planet, err := d.client.ReadPlanetByName(state.Name.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error reading planet",
			fmt.Sprintf("Error reading planet %s, got error %s", state.Name.ValueString(), err))
		return
	}

	state.Id = types.StringValue(planet.Id)
	state.Name = types.StringValue(planet.Name)
	state.Diameter = types.Int64Value(planet.Diameter)
	state.Population = types.Int64Value(planet.Population)

	tflog.Trace(ctx, "read a state source")

	// Save state into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
