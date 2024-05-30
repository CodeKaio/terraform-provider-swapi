// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"fmt"
	"github.com/codekaio/terraform-provider-swapi/internal/swapi"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &PlanetResource{}
var _ resource.ResourceWithImportState = &PlanetResource{}

func NewPlanetResource() resource.Resource {
	return &PlanetResource{}
}

// PlanetResource defines the resource implementation.
type PlanetResource struct {
	client *swapi.SWAPIClient
}

// PlanetResourceModel describes the resource data model.
type PlanetResourceModel struct {
	Id         types.String `tfsdk:"id"`
	Name       types.String `tfsdk:"name"`
	Diameter   types.Int64  `tfsdk:"diameter"`
	Population types.Int64  `tfsdk:"population"`
}

func (r *PlanetResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_planet"
}

func (r *PlanetResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Planet resource",

		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				MarkdownDescription: "Planet Name",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Planet id",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"diameter": schema.Int64Attribute{
				MarkdownDescription: "Planet diameter",
				Optional:            true,
			},
			"population": schema.Int64Attribute{
				MarkdownDescription: "Planet population",
				Optional:            true,
			},
		},
	}
}

func (r *PlanetResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*swapi.SWAPIClient)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *http.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	r.client = client
}

func (r *PlanetResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var state PlanetResourceModel

	// Read Terraform plan state into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	var planet swapi.Planet
	planet.Name = state.Name.ValueString()
	planet.Diameter = state.Diameter.ValueInt64()
	planet.Population = state.Population.ValueInt64()

	_, err := r.client.CreateOrUpdatePlanet(&planet)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating planet",
			fmt.Sprintf("Error creating planet %s, got error %s", state.Name.ValueString(), err))
		return
	}

	state.Id = types.StringValue(planet.Id)
	state.Name = types.StringValue(planet.Name)
	state.Diameter = types.Int64Value(planet.Diameter)
	state.Population = types.Int64Value(planet.Population)

	// Write logs using the tflog package
	// Documentation: https://terraform.io/plugin/log
	tflog.Trace(ctx, "created planet")

	// Save state into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *PlanetResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state PlanetResourceModel

	// Read Terraform prior state into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	planet, err := r.client.ReadPlanetById(state.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading planet",
			fmt.Sprintf("Error reading planet %s, got error %s", state.Id.ValueString(), err))
		return
	}

	state.Id = types.StringValue(planet.Id)
	state.Name = types.StringValue(planet.Name)
	state.Population = types.Int64Value(planet.Population)
	state.Diameter = types.Int64Value(planet.Diameter)

	// Save updated state into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *PlanetResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var state PlanetResourceModel

	// Read Terraform plan state into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	var planet swapi.Planet
	planet.Name = state.Name.ValueString()
	planet.Id = state.Id.ValueString()

	_, err := r.client.CreateOrUpdatePlanet(&planet)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating planet",
			fmt.Sprintf("Error creating planet %s, got error %s", state.Name.ValueString(), err))
		return
	}

	state.Id = types.StringValue(planet.Id)
	state.Name = types.StringValue(planet.Name)

	// Save updated state into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *PlanetResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state PlanetResourceModel

	// Read Terraform prior state state into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.DeletePlanet(state.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error deleting planet",
			fmt.Sprintf("Error reading planet %s, got error %s", state.Id.ValueString(), err))
		return
	}
}

func (r *PlanetResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
