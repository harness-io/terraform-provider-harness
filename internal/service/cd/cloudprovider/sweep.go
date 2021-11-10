package cloudprovider

import (
	"strings"

	"github.com/harness-io/terraform-provider-harness/internal/sweep"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func init() {
	resource.AddTestSweepers("harness_cloudprovider", &resource.Sweeper{
		Name: "harness_cloudprovider",
		F:    testSweepCloudProviders,
	})
}

func testSweepCloudProviders(r string) error {
	c := sweep.SweeperClient

	limit := 100
	offset := 0
	hasMore := true

	for hasMore {

		cloudProviders, _, err := c.CDClient.CloudProviderClient.ListCloudProviders(limit, offset)
		if err != nil {
			return err
		}

		for _, cp := range cloudProviders {
			if strings.HasPrefix(cp.Name, "Test") {
				if err = c.CDClient.CloudProviderClient.DeleteCloudProvider(cp.Id); err != nil {
					return err
				}
			}
		}

		hasMore = len(cloudProviders) == limit
	}

	return nil
}
