// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"reflect"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &KeystoreResource{}

func NewKeystoreResource() resource.Resource {
	return &KeystoreResource{}
}

// KeystoreResource defines the resource implementation.
type KeystoreResource struct {
}

// KeystoreResourceModel describes the resource data model.
type KeystoreResourceModel struct {
	Id                 types.String `tfsdk:"id"`
	Password           types.String `tfsdk:"password"`
	File               types.String `tfsdk:"file"`
	CommonName         types.String `tfsdk:"common_name"`
	Organization       types.String `tfsdk:"organization"`
	OrganizationalUnit types.String `tfsdk:"organizational_unit"`
	Locality           types.String `tfsdk:"locality"`
	State              types.String `tfsdk:"state"`
	Country            types.String `tfsdk:"country"`
}

func (r KeystoreResourceModel) ToKeystoreModel() KeystoreModel {
	return KeystoreModel{
		Password: r.Password.ValueString(),
		File:     r.File.ValueString(),
		DistinguishedName: DistinguishedName{
			CommonName:         r.CommonName.ValueString(),
			Organization:       r.Organization.ValueString(),
			OrganizationalUnit: r.OrganizationalUnit.ValueString(),
			Locality:           r.Locality.ValueString(),
			State:              r.State.ValueString(),
			Country:            r.Country.ValueString(),
		},
	}
}

func (r *KeystoreResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_keystore"
}

func (r *KeystoreResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {

	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		Description: `
        Keystore resource which creates a base64 encoded PKCS12 keystore file valid for 25 years using the keytool utility.
        The machine running Terraform needs to have the keytool utility installed. https://docs.oracle.com/javase/8/docs/technotes/tools/unix/keytool.html.
        The file is persisted solely within the Terraform state as base64 text.
        `,

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "Generated UUID for the keystore",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"password": schema.StringAttribute{
				Description: "Password for the keystore and the single key in the keystore",
				Required:    true,
				Sensitive:   true,
			},
			"file": schema.StringAttribute{

				Description: "Base64 encoded keystore file",
				Computed:    true,
				Sensitive:   true,
			},
			"common_name": schema.StringAttribute{
				Description: "Common Name (CN)",
				Optional:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplaceIfConfigured(),
				},
			},
			"organization": schema.StringAttribute{
				Description: "Organization (O)",
				Optional:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplaceIfConfigured(),
				},
			},
			"organizational_unit": schema.StringAttribute{
				Description: "Organizational Unit (OU)",
				Optional:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplaceIfConfigured(),
				},
			},
			"locality": schema.StringAttribute{
				Description: "Locality (L)",
				Optional:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplaceIfConfigured(),
				},
			},
			"state": schema.StringAttribute{
				Description: "State (S)",
				Optional:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplaceIfConfigured(),
				},
			},
			"country": schema.StringAttribute{
				Description: "Country (C)",
				Optional:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplaceIfConfigured(),
				},
			},
		},
	}
}

func (r *KeystoreResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}
}

func (r *KeystoreResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data KeystoreResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	model := data.ToKeystoreModel()

	b64File, err := model.CreateKeystoreBase64()
	model.File = b64File

	if err != nil {
		resp.Diagnostics.AddError("Error during create operation", err.Error())
		return
	}

	data.Id = types.StringValue(uuid.New().String())
	data.File = types.StringValue(model.File)

	// Write logs using the tflog package
	// Documentation: https://terraform.io/plugin/log
	tflog.Trace(ctx, "created a resource")

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *KeystoreResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data KeystoreResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *KeystoreResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data KeystoreResourceModel
	var oldData KeystoreResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &oldData)...)

	if resp.Diagnostics.HasError() {
		return
	}

	if reflect.DeepEqual(data, oldData) {
		resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
		return
	}

	newModel := data.ToKeystoreModel()
	oldModel := oldData.ToKeystoreModel()
	b64File, err := oldModel.UpdateKeystoreBase64(newModel.Password)
	if err != nil {
		resp.Diagnostics.AddError("Error during update operation", err.Error())
		return
	}

	newModel.File = b64File
	data.File = types.StringValue(newModel.File)

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *KeystoreResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data KeystoreResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}
}
