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
	Id 	   types.String `tfsdk:"id"`
	Password   types.String `tfsdk:"password"`
	Base64Text types.String `tfsdk:"base64_text"`
	
}

func (r KeystoreResourceModel) ToKeystoreModel() KeystoreModel {
	return KeystoreModel{
		Password: r.Password.ValueString(),
		Base64Text: r.Base64Text.ValueString(),
	}
}

func (r *KeystoreResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_keystore"
}

func (r *KeystoreResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		Description: `
		Keystore resource which creates a base64 encoded keystore file using the keytool utility.
		The machine running Terraform needs to have the keytool utility installed. https://docs.oracle.com/javase/8/docs/technotes/tools/unix/keytool.html.
		File is persisted solely within the Terraform state as base64 text.
		`,

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"password": schema.StringAttribute{
				Required:  true,
				Sensitive: true,
			},
			"base64_text": schema.StringAttribute{
				Computed:  true,
				Sensitive: true,
				

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

	model:= data.ToKeystoreModel()
	
	b64File, err := model.CreateKeystoreBase64()
	model.Base64Text = b64File
	

	if err != nil {
		resp.Diagnostics.AddError("Error during create operation", err.Error())
		return
	}

	data.Id = types.StringValue(uuid.New().String())
	data.Base64Text = types.StringValue(model.Base64Text)


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

	newModel.Base64Text = b64File
	data.Base64Text = types.StringValue(newModel.Base64Text)


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
