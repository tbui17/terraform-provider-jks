// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/compare"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
)

type KeystoreTestContext struct {
	Namespace string
	Type      string
	Name      string
	FullName  string
}

const (
	TestNamespace            = "jks"
	TestType                 = "keystore"
	TestResourceInstanceName = "test"
	TestResourceFullName     = TestNamespace + "_" + TestType + "." + TestResourceInstanceName
)

func TestAccKeystoreResource(t *testing.T) {

	model := KeystoreModel{
		Password: "MyPassword12345",
	}

	model2 := model
	model2.Password = "MyPassword"

	valuesDifferCtx := statecheck.CompareValue(compare.ValuesDiffer())
	valuesSameCtx := statecheck.CompareValue(compare.ValuesSame())

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: ToTfResourceString(model),
				Check: resource.ComposeAggregateTestCheckFunc(
					testCheckResourceAttrLengthGreater(TestResourceFullName, "id", 0),
					resource.TestCheckResourceAttr(TestResourceFullName, "password", model.Password),
					testCheckResourceAttrLengthGreater(TestResourceFullName, "base64_text", 1000),
				),
				ConfigStateChecks: []statecheck.StateCheck{
					valuesDifferCtx.AddStateValue(
						TestResourceFullName,
						tfjsonpath.New("base64_text"),
					),
					valuesSameCtx.AddStateValue(
						TestResourceFullName,
						tfjsonpath.New("id"),
					),
				},
			},
			// Update and Read testing
			{
				Config: ToTfResourceString(model2),
				ConfigStateChecks: []statecheck.StateCheck{
					valuesDifferCtx.AddStateValue(
						TestResourceFullName,
						tfjsonpath.New("base64_text"),
					),
					valuesSameCtx.AddStateValue(
						TestResourceFullName,
						tfjsonpath.New("id"),
					),
				},
				Check: resource.ComposeAggregateTestCheckFunc(
					testCheckKeepEncodedFile(),
					resource.TestCheckResourceAttr(TestResourceFullName, "password", model2.Password),
				),
			},

			// Delete testing automatically occurs in TestCase
		},
	})
}

func testCheckKeepEncodedFile() resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rootModule := s.RootModule()
		rs, ok := rootModule.Resources[TestResourceFullName]
		if !ok {
			return fmt.Errorf("not found: %s", TestResourceFullName)
		}

		val, ok := rs.Primary.Attributes["base64_text"]
		if !ok {
			return fmt.Errorf("attribute not found: %s", "base64_text")
		}
		current := rootModule.Resources[TestResourceFullName].Primary.Attributes["base64_text"]
		if len(current) < 1000 {
			return fmt.Errorf("attribute %s is %s, expected length > 1000", "base64_text", current)
		}
		if current != val {
			return fmt.Errorf("attribute %s is %s, expected %s", "base64_text", val, current)
		}

		return nil
	}
}

func testCheckResourceAttrLengthGreater(resourceName, attributeName string, minLength int) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("not found: %s", resourceName)
		}

		attr, ok := rs.Primary.Attributes[attributeName]
		if !ok {
			return fmt.Errorf("attribute not found: %s", attributeName)
		}

		if len(attr) < minLength {
			return fmt.Errorf("attribute %s length is %d, less than %d", attributeName, len(attr), minLength)
		}

		return nil
	}
}

func ToTfResourceString(m KeystoreModel) string {
	return fmt.Sprintf(`
	resource "jks_keystore" "test" {
		password = %q
		
		}`, m.Password,
	)
}
