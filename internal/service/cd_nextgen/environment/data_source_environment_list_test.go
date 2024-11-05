package environment_test

import (
	"fmt"
	"testing"

	"github.com/harness/harness-go-sdk/harness/utils"
	"github.com/harness/terraform-provider-harness/internal/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceEnvironmentList(t *testing.T) {

	id := fmt.Sprintf("%s_%s", t.Name(), utils.RandStringBytes(6))
	name := id
	resourceName := "data.harness_platform_environment_list.test"

	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { acctest.TestAccPreCheck(t) },
		ProviderFactories: acctest.ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceEnvironmentList(id, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "environments.0.identifier", id),
					resource.TestCheckResourceAttr(resourceName, "environments.0.name", name),
				),
			},
		},
	})
}

func TestAccDataSourceEnvironmentListAtOrgLevel(t *testing.T) {

	id := fmt.Sprintf("%s_%s", t.Name(), utils.RandStringBytes(6))
	name := id
	resourceName := "data.harness_platform_environment_list.test"

	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { acctest.TestAccPreCheck(t) },
		ProviderFactories: acctest.ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceEnvironmentListOrgLevel(id, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "environments.0.identifier", id),
					resource.TestCheckResourceAttr(resourceName, "environments.0.name", name),
				),
			},
		},
	})
}

func testAccDataSourceEnvironmentList(id string, name string) string {
	return fmt.Sprintf(`
		resource "harness_platform_organization" "test" {
			identifier = "%[1]s"
			name = "%[2]s"
		}

		resource "harness_platform_project" "test" {
			identifier = "%[1]s"
			name = "%[2]s"
			org_id = harness_platform_organization.test.id
			color = "#472848"
		}

		resource "harness_platform_environment" "test" {
			identifier = "%[1]s"
			name = "%[2]s"
			org_id = harness_platform_project.test.org_id
			project_id = harness_platform_project.test.id
			color = "#0063F7"
			type = "PreProduction"
		}

		data "harness_platform_environment_list" "test" {
			org_id = harness_platform_organization.test.id
			project_id = harness_platform_environment.test.project_id
		}

`, id, name)
}

func testAccDataSourceEnvironmentListOrgLevel(id string, name string) string {
	return fmt.Sprintf(`
		resource "harness_platform_organization" "test" {
			identifier = "%[1]s"
			name = "%[2]s"
		}

		resource "harness_platform_environment" "test" {
			identifier = "%[1]s"
			org_id = harness_platform_organization.test.id
			name = "%[2]s"
			color = "#0063F7"
			type = "PreProduction"
		}

		data "harness_platform_environment_list" "test" {
			org_id = harness_platform_environment.test.org_id
		}

`, id, name)
}
