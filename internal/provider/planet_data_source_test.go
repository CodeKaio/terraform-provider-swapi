// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccExampleDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: providerConfig + testAccExampleDataSourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.swapi_planet.dagobah", "name", "Dagobah"),
					resource.TestCheckResourceAttrSet("data.swapi_planet.dagobah", "id"),
				),
			},
		},
	})
}

const testAccExampleDataSourceConfig = `
data "swapi_planet" "dagobah" {
  name = "Dagobah"
}
`
