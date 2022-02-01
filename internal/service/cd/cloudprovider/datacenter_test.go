package cloudprovider_test

import (
	"errors"
	"fmt"
	"testing"

	sdk "github.com/harness-io/harness-go-sdk"
	"github.com/harness-io/harness-go-sdk/harness/cd/cac"
	"github.com/harness-io/harness-go-sdk/harness/utils"
	"github.com/harness-io/terraform-provider-harness/internal/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/stretchr/testify/require"
)

func TestAccResourceDataCenterCloudProviderConnector(t *testing.T) {

	var (
		name         = fmt.Sprintf("%s_%s", t.Name(), utils.RandStringBytes(4))
		updatedName  = fmt.Sprintf("%s_updated", name)
		resourceName = "harness_cloudprovider_datacenter.test"
	)

	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { acctest.TestAccPreCheck(t) },
		ProviderFactories: acctest.ProviderFactories,
		CheckDestroy:      testAccCloudProviderDestroy(resourceName),
		Steps: []resource.TestStep{
			{
				Config: testAccResourceDataCenterCloudProvider(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", name),
					testAccCheckDataCenterCloudProviderExists(t, resourceName, name),
				),
			},
			{
				Config: testAccResourceDataCenterCloudProvider(updatedName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDataCenterCloudProviderExists(t, resourceName, name),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccResourceDataCenterCloudProviderConnector_DeleteUnderlyingResource(t *testing.T) {

	var (
		name         = fmt.Sprintf("%s_%s", t.Name(), utils.RandStringBytes(4))
		resourceName = "harness_cloudprovider_datacenter.test"
	)

	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { acctest.TestAccPreCheck(t) },
		ProviderFactories: acctest.ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceDataCenterCloudProvider(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", name),
					testAccCheckDataCenterCloudProviderExists(t, resourceName, name),
				),
			},
			{
				PreConfig: func() {
					acctest.TestAccConfigureProvider()
					c := acctest.TestAccProvider.Meta().(*sdk.Session)
					cp, err := c.CDClient.CloudProviderClient.GetPhysicalDatacenterCloudProviderByName(name)
					require.NoError(t, err)
					require.NotNil(t, cp)

					err = c.CDClient.CloudProviderClient.DeleteCloudProvider(cp.Id)
					require.NoError(t, err)
				},
				Config:             testAccResourceDataCenterCloudProvider(name),
				PlanOnly:           true,
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func testAccResourceDataCenterCloudProvider(name string) string {
	return fmt.Sprintf(`
		resource "harness_cloudprovider_datacenter" "test" {
			name = "%[1]s"

			usage_scope {
				environment_filter_type = "NON_PRODUCTION_ENVIRONMENTS"
			}
			
			usage_scope {
				environment_filter_type = "PRODUCTION_ENVIRONMENTS"
			}
		}	
`, name)
}

func testAccCheckDataCenterCloudProviderExists(t *testing.T, resourceName, cloudProviderName string) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		cp := &cac.PhysicalDatacenterCloudProvider{}
		err := testAccGetCloudProvider(resourceName, state, cp)
		if err != nil {
			return err
		}
		return nil
	}
}

func testAccGetCloudProvider(resourceName string, state *terraform.State, respObj interface{}) error {
	r := acctest.TestAccGetResource(resourceName, state)
	if r == nil {
		return errors.New("Resource not found")
	}

	c := acctest.TestAccGetApiClientFromProvider()
	name := r.Primary.Attributes["name"]

	err := c.CDClient.ConfigAsCodeClient.GetCloudProviderByName(name, respObj)
	if err != nil {
		return err
	}

	return nil
}

func testAccCloudProviderDestroy(resourceName string) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		cp := &cac.PhysicalDatacenterCloudProvider{}
		err := testAccGetCloudProvider(resourceName, state, &cp)
		if err != nil {
			return err
		}

		if !cp.IsEmpty() {
			return fmt.Errorf("cloud Provider still exists")
		}

		return nil
	}
}
