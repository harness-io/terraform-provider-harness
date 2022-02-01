package sweep_test

import (
	"testing"

	sdk "github.com/harness-io/harness-go-sdk"
	_ "github.com/harness-io/terraform-provider-harness/internal/service/cd/application"
	_ "github.com/harness-io/terraform-provider-harness/internal/service/cd/cloudprovider"
	_ "github.com/harness-io/terraform-provider-harness/internal/service/cd/secrets"
	"github.com/harness-io/terraform-provider-harness/internal/sweep"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestMain(m *testing.M) {
	sweep.SweeperClient = sdk.NewSession(&sdk.SessionOptions{})
	resource.TestMain(m)
}
