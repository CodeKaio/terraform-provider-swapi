// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccExampleResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: providerConfig + testAccPlanetResourceConfig("Hoth"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("swapi_planet.test", "name", "Hoth"),
					resource.TestCheckResourceAttrSet("swapi_planet.test", "id"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "swapi_planet.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Update and Read testing
			{
				Config: providerConfig + testAccPlanetResourceConfig("Mustafar"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("swapi_planet.test", "name", "Mustafar"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccPlanetResourceConfig(planetName string) string {
	return fmt.Sprintf(`
resource "swapi_planet" "test" {
  name = %[1]q
}
`, planetName)
}
