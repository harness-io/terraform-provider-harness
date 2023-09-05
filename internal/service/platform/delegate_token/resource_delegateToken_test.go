package delegatetoken_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/antihax/optional"
	"github.com/harness/harness-go-sdk/harness/nextgen"
	"github.com/harness/harness-go-sdk/harness/utils"
	"github.com/harness/terraform-provider-harness/internal/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccResourceDelegateToken(t *testing.T) {
	name := utils.RandStringBytes(5)
	account_id := os.Getenv("HARNESS_ACCOUNT_ID")

	resourceName := "harness_platform_delegatetoken.test"

	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { acctest.TestAccPreCheck(t) },
		ProviderFactories: acctest.ProviderFactories,
		CheckDestroy:      testDelegateTokenDestroy(resourceName),
		Steps: []resource.TestStep{
			{
				Config: tesAccResourceDelegateToken(name, account_id),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "token_status", "ACTIVE"),
				),
			},
		},
	})
}

func tesAccResourceDelegateToken(name string, accountId string) string {
	return fmt.Sprintf(`
		resource "harness_platform_delegatetoken" "test" {
			identifier = "%[1]s"
			name = "%[1]s"
			account_id = "%[2]s"
			token_status = "ACTIVE"
		}
		`, name, accountId)
}

func testAccGetResourceDelegateToken(resourceName string, state *terraform.State) (*nextgen.DelegateTokenDetails, error) {
	d := acctest.TestAccGetResource(resourceName, state)
	c, ctx := acctest.TestAccGetPlatformClientWithContext()

	resp, _, err := c.DelegateTokenResourceApi.GetDelegateTokens(ctx, c.AccountId, &nextgen.DelegateTokenResourceApiGetDelegateTokensOpts{
		OrgIdentifier:     buildField(d, "org_id"),
		ProjectIdentifier: buildField(d, "project_id"),
		Name:              buildField(d, "name"),
		Status:            buildField(d, "token_status"),
	})

	if err != nil {
		return nil, err
	}

	if resp.Resource == nil {
		return nil, nil
	}

	return &resp.Resource[0], nil
}

func testDelegateTokenDestroy(resourceName string) resource.TestCheckFunc {
	var token *nextgen.DelegateTokenDetails
	return func(state *terraform.State) error {
		token, _ = testAccGetResourceDelegateToken(resourceName, state)
		if token.Status != "REVOKED" {
			return fmt.Errorf("Token is not revoked : %s", token.Name)
		}

		return nil
	}
}

func buildField(r *terraform.ResourceState, field string) optional.String {
	if attr, ok := r.Primary.Attributes[field]; ok {
		return optional.NewString(attr)
	}
	return optional.EmptyString()
}
