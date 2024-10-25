package environment_group_test

import (
	"fmt"
	"testing"

	"github.com/harness/harness-go-sdk/harness/utils"
	"github.com/harness/terraform-provider-harness/internal/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceEnvironmentGroup(t *testing.T) {

	id := fmt.Sprintf("%s_%s", t.Name(), utils.RandStringBytes(6))
	name := id
	resourceName := "data.harness_platform_environment_group.test"

	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { acctest.TestAccPreCheck(t) },
		ProviderFactories: acctest.ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceEnvironmentGroup(id, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "identifier", id),
					resource.TestCheckResourceAttr(resourceName, "org_id", id),
					resource.TestCheckResourceAttr(resourceName, "project_id", id),
					resource.TestCheckResourceAttr(resourceName, "color", "#0063F7"),
				),
			},
		},
	})
}

func TestAccDataSourceEnvironmentGroupOrgLevel(t *testing.T) {

	id := fmt.Sprintf("%s_%s", t.Name(), utils.RandStringBytes(6))
	name := id
	resourceName := "data.harness_platform_environment_group.test"

	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { acctest.TestAccPreCheck(t) },
		ProviderFactories: acctest.ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceEnvironmentGroupOrgLevel(id, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "identifier", id),
					resource.TestCheckResourceAttr(resourceName, "org_id", id),
					resource.TestCheckResourceAttr(resourceName, "color", "#0063F7"),
				),
			},
		},
	})
}

func TestAccDataSourceEnvironmentGroupAccountLevel(t *testing.T) {

	id := fmt.Sprintf("%s_%s", t.Name(), utils.RandStringBytes(6))
	name := id
	resourceName := "data.harness_platform_environment_group.test"

	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { acctest.TestAccPreCheck(t) },
		ProviderFactories: acctest.ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceEnvironmentGroupOrgLevel(id, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "identifier", id),
					resource.TestCheckResourceAttr(resourceName, "color", "#0063F7"),
				),
			},
		},
	})
}

func testAccDataSourceEnvironmentGroup(id string, name string) string {
	return fmt.Sprintf(`
		resource "harness_platform_organization" "test" {
			identifier = "%[1]s"
			name = "%[2]s"
		}

		resource "harness_platform_project" "test" {
			identifier = "%[1]s"
			name = "%[2]s"
			org_id = harness_platform_organization.test.id
			color = "#0063F7"
		}

		resource "harness_platform_environment_group" "test" {
			identifier = "%[1]s"
			org_id = harness_platform_organization.test.id
			project_id = harness_platform_project.test.id
			color = "#0063F7"
			yaml = <<-EOT
			     environmentGroup:
			                 name: "%[1]s"
			                 identifier: "%[1]s"
			                 description: "temp"
			                 orgIdentifier: ${harness_platform_organization.test.id}
			                 projectIdentifier: ${harness_platform_project.test.id}
			                 envIdentifiers: []
		  EOT
		}

		data "harness_platform_environment_group" "test" {
			identifier = "%[1]s"
			org_id = harness_platform_environment_group.test.id
			project_id = harness_platform_environment_group.test.id
		}
`, id, name)
}

func testAccDataSourceEnvironmentGroupOrgLevel(id string, name string) string {
	return fmt.Sprintf(`
		resource "harness_platform_organization" "test" {
			identifier = "%[1]s"
			name = "%[2]s"
		}

		resource "harness_platform_environment_group" "test" {
			identifier = "%[1]s"
			org_id = harness_platform_organization.test.id
			color = "#0063F7"
			yaml = <<-EOT
			     environmentGroup:
			                 name: "%[1]s"
			                 identifier: "%[1]s"
			                 description: "temp"
			                 orgIdentifier: ${harness_platform_organization.test.id}
			                 envIdentifiers: []
		  EOT
		}

		data "harness_platform_environment_group" "test" {
			identifier = "%[1]s"
			org_id = harness_platform_environment_group.test.id
		}
`, id, name)
}

func testAccDataSourceEnvironmentGroupAccountLevel(id string, name string) string {
	return fmt.Sprintf(`
		resource "harness_platform_environment_group" "test" {
			identifier = "%[1]s"
			color = "#0063F7"
			yaml = <<-EOT
			     environmentGroup:
			                 name: "%[1]s"
			                 identifier: "%[1]s"
			                 description: "temp"
			                 envIdentifiers: []
		  EOT
		}

		data "harness_platform_environment_group" "test" {
			identifier = "%[1]s"
		}
`, id, name)
}
