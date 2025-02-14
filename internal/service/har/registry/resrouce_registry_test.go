package registry_test

import (
	"fmt"
	"github.com/harness/terraform-provider-harness/internal/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"os"
	"testing"
)

// Virtual Docker Registry
func TestAccResourceVirtualDockerRegistry(t *testing.T) {
	id := fmt.Sprintf("tf_auto_virtual_docker_registry")
	resourceName := "harness_platform_har_registry.test"
	accountId := os.Getenv("HARNESS_ACCOUNT_ID")

	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { acctest.TestAccPreCheck(t) },
		ProviderFactories: acctest.ProviderFactories,
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					acctest.TestAccConfigureProvider()
					_, _ = acctest.TestAccGetHarClientWithContext()
				},
				Config: testAccResourceVirtualDockerRegistry(id, accountId),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "identifier", id),
				),
			},
		},
	})
}
func TestOrgResourceVirtualDockerRegistry(t *testing.T) {
	id := fmt.Sprintf("tf_auto_virtual_docker_registry")
	resourceName := "harness_platform_har_registry.test"
	accountId := os.Getenv("HARNESS_ACCOUNT_ID")

	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { acctest.TestAccPreCheck(t) },
		ProviderFactories: acctest.ProviderFactories,
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					acctest.TestAccConfigureProvider()
					_, _ = acctest.TestAccGetHarClientWithContext()
				},
				Config: testOrgResourceVirtualDockerRegistry(id, accountId),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "identifier", id),
				),
			},
		},
	})
}
func TestProjectResourceVirtualDockerRegistry(t *testing.T) {
	id := fmt.Sprintf("tf_auto_virtual_docker_registry")
	resourceName := "harness_platform_har_registry.test"
	accountId := os.Getenv("HARNESS_ACCOUNT_ID")

	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { acctest.TestAccPreCheck(t) },
		ProviderFactories: acctest.ProviderFactories,
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					acctest.TestAccConfigureProvider()
					_, _ = acctest.TestAccGetHarClientWithContext()
				},
				Config: testProjResourceVirtualDockerRegistry(id, accountId),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "identifier", id),
				),
			},
		},
	})
}

// Upstream Docker Registry --- UserPassword
func TestAccResourceUpstreamDockerRegistry(t *testing.T) {
	id := fmt.Sprintf("tf_auto_upstream_docker_registry")
	resourceName := "harness_platform_har_registry.test"
	accountId := os.Getenv("HARNESS_ACCOUNT_ID")

	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { acctest.TestAccPreCheck(t) },
		ProviderFactories: acctest.ProviderFactories,
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					acctest.TestAccConfigureProvider()
					_, _ = acctest.TestAccGetHarClientWithContext()
				},
				Config: testAccResourceUpstreamDockerRegistry(id, accountId),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "identifier", id),
				),
			},
		},
	})
}
func TestOrgResourceUpstreamDockerRegistry(t *testing.T) {
	id := fmt.Sprintf("tf_auto_upstream_docker_registry")
	resourceName := "harness_platform_har_registry.test"
	accountId := os.Getenv("HARNESS_ACCOUNT_ID")

	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { acctest.TestAccPreCheck(t) },
		ProviderFactories: acctest.ProviderFactories,
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					acctest.TestAccConfigureProvider()
					_, _ = acctest.TestAccGetHarClientWithContext()
				},
				Config: testOrgResourceUpstreamDockerRegistry(id, accountId),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "identifier", id),
				),
			},
		},
	})
}
func TestProjectResourceUpstreamDockerRegistry(t *testing.T) {
	id := fmt.Sprintf("tf_auto_upstream_docker_registry")
	resourceName := "harness_platform_har_registry.test"
	accountId := os.Getenv("HARNESS_ACCOUNT_ID")

	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { acctest.TestAccPreCheck(t) },
		ProviderFactories: acctest.ProviderFactories,
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					acctest.TestAccConfigureProvider()
					_, _ = acctest.TestAccGetHarClientWithContext()
				},
				Config: testProjResourceUpstreamDockerRegistry(id, accountId),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "identifier", id),
				),
			},
		},
	})
}

func testAccResourceVirtualDockerRegistry(id string, accId string) string {
	return fmt.Sprintf(`

 resource "harness_platform_har_registry" "test" {
   identifier   = "%[1]s"
   space_ref    = "%[2]s"
   package_type = "DOCKER"

   config {
    type = "VIRTUAL"
   }
   parent_ref = "%[2]s"
 }
`, id, accId)
}

func testOrgResourceVirtualDockerRegistry(id string, accId string) string {
	return fmt.Sprintf(`
 resource "harness_platform_organization" "test" {
  identifier = "%[1]s_org"
  name = "%[1]s"
 }

 resource "harness_platform_har_registry" "test" {
   identifier   = "%[1]s"
   space_ref    = "%[2]s/${harness_platform_organization.test.identifier}"
   package_type = "DOCKER"

   config {
    type = "VIRTUAL"
   }
   parent_ref = "%[2]s/${harness_platform_organization.test.identifier}"
 }
`, id, accId)
}

func testProjResourceVirtualDockerRegistry(id string, accId string) string {
	return fmt.Sprintf(`
 resource "harness_platform_organization" "test" {
  identifier = "%[1]s_org"
  name = "%[1]s"
 }

 resource "harness_platform_project" "test" {
  identifier = "%[1]s_project"
  name = "%[1]s"
  org_id = harness_platform_organization.test.id
  color = "#472848"
 }
 resource "harness_platform_har_registry" "test" {
   identifier   = "%[1]s"
   space_ref    = "%[2]s/${harness_platform_organization.test.identifier}/${harness_platform_project.test.identifier}"
   package_type = "DOCKER"

   config {
    type = "VIRTUAL"
   }
   parent_ref = "%[2]s/${harness_platform_organization.test.identifier}/${harness_platform_project.test.identifier}"
 }
`, id, accId)
}

func testAccResourceUpstreamDockerRegistry(id string, accId string) string {
	return fmt.Sprintf(`
 resource "harness_platform_har_registry" "test" {
   identifier   = "%[1]s"
   space_ref    = "%[2]s"
   package_type = "DOCKER"

   config {
		type = "UPSTREAM"
		auth_type = "UserPassword"
		source = "Dockerhub"
		auth {
			auth_type = "UserPassword"
			user_name = "username"
			secret_identifier = "Secret_Token"
			secret_space_path = "%[2]s"
		}
   }
   parent_ref = "%[2]s"
 }
`, id, accId)
}

func testOrgResourceUpstreamDockerRegistry(id string, accId string) string {
	return fmt.Sprintf(`
 resource "harness_platform_organization" "test" {
  identifier = "%[1]s_org"
  name = "%[1]s"
 }

 resource "harness_platform_har_registry" "test" {
   identifier   = "%[1]s"
   space_ref    = "%[2]s/${harness_platform_organization.test.identifier}"
   package_type = "DOCKER"

   config {
		type = "UPSTREAM"
		auth_type = "UserPassword"
		source = "Dockerhub"
		auth {
			auth_type = "UserPassword"
			user_name = "username"
			secret_identifier = "Secret_Token"
			secret_space_path = "%[2]s"
		}
   }
   parent_ref = "%[2]s/${harness_platform_organization.test.identifier}"
 }
`, id, accId)
}

func testProjResourceUpstreamDockerRegistry(id string, accId string) string {
	return fmt.Sprintf(`
 resource "harness_platform_organization" "test" {
  identifier = "%[1]s_org"
  name = "%[1]s"
 }

 resource "harness_platform_project" "test" {
  identifier = "%[1]s_project"
  name = "%[1]s"
  org_id = harness_platform_organization.test.id
  color = "#472848"
 }
 resource "harness_platform_har_registry" "test" {
   identifier   = "%[1]s"
   space_ref    = "%[2]s/${harness_platform_organization.test.identifier}/${harness_platform_project.test.identifier}"
   package_type = "DOCKER"

   config {
		type = "UPSTREAM"
		auth_type = "UserPassword"
		source = "Dockerhub"
		auth {
			auth_type = "UserPassword"
			user_name = "username"
			secret_identifier = "Secret_Token"
			secret_space_path = "%[2]s"
		}
   }
   parent_ref = "%[2]s/${harness_platform_organization.test.identifier}/${harness_platform_project.test.identifier}"
 }
`, id, accId)
}

// Virtual Helm Registry
func TestAccResourceVirtualHelmRegistry(t *testing.T) {
	id := fmt.Sprintf("tf_auto_virtual_helm_registry")
	resourceName := "harness_platform_har_registry.test"
	accountId := os.Getenv("HARNESS_ACCOUNT_ID")

	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { acctest.TestAccPreCheck(t) },
		ProviderFactories: acctest.ProviderFactories,
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					acctest.TestAccConfigureProvider()
					_, _ = acctest.TestAccGetHarClientWithContext()
				},
				Config: testAccResourceVirtualHelmRegistry(id, accountId),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "identifier", id),
				),
			},
		},
	})
}
func TestOrgResourceVirtualHelmRegistry(t *testing.T) {
	id := fmt.Sprintf("tf_auto_virtual_helm_registry")
	resourceName := "harness_platform_har_registry.test"
	accountId := os.Getenv("HARNESS_ACCOUNT_ID")

	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { acctest.TestAccPreCheck(t) },
		ProviderFactories: acctest.ProviderFactories,
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					acctest.TestAccConfigureProvider()
					_, _ = acctest.TestAccGetHarClientWithContext()
				},
				Config: testOrgResourceVirtualHelmRegistry(id, accountId),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "identifier", id),
				),
			},
		},
	})
}
func TestProjectResourceVirtualHelmRegistry(t *testing.T) {
	id := fmt.Sprintf("tf_auto_virtual_helm_registry")
	resourceName := "harness_platform_har_registry.test"
	accountId := os.Getenv("HARNESS_ACCOUNT_ID")

	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { acctest.TestAccPreCheck(t) },
		ProviderFactories: acctest.ProviderFactories,
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					acctest.TestAccConfigureProvider()
					_, _ = acctest.TestAccGetHarClientWithContext()
				},
				Config: testProjResourceVirtualHelmRegistry(id, accountId),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "identifier", id),
				),
			},
		},
	})
}

// Upstream Helm Registry --- UserPass
func TestAccResourceUpstreamHelmRegistry(t *testing.T) {
	id := fmt.Sprintf("tf_auto_upstream_helm_registry")
	resourceName := "harness_platform_har_registry.test"
	accountId := os.Getenv("HARNESS_ACCOUNT_ID")

	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { acctest.TestAccPreCheck(t) },
		ProviderFactories: acctest.ProviderFactories,
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					acctest.TestAccConfigureProvider()
					_, _ = acctest.TestAccGetHarClientWithContext()
				},
				Config: testAccResourceUpstreamHelmRegistry(id, accountId),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "identifier", id),
				),
			},
		},
	})
}
func TestOrgResourceUpstreamHelmRegistry(t *testing.T) {
	id := fmt.Sprintf("tf_auto_upstream_helm_registry")
	resourceName := "harness_platform_har_registry.test"
	accountId := os.Getenv("HARNESS_ACCOUNT_ID")

	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { acctest.TestAccPreCheck(t) },
		ProviderFactories: acctest.ProviderFactories,
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					acctest.TestAccConfigureProvider()
					_, _ = acctest.TestAccGetHarClientWithContext()
				},
				Config: testOrgResourceUpstreamHelmRegistry(id, accountId),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "identifier", id),
				),
			},
		},
	})
}
func TestProjectResourceUpstreamHelmRegistry(t *testing.T) {
	id := fmt.Sprintf("tf_auto_upstream_helm_registry")
	resourceName := "harness_platform_har_registry.test"
	accountId := os.Getenv("HARNESS_ACCOUNT_ID")

	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { acctest.TestAccPreCheck(t) },
		ProviderFactories: acctest.ProviderFactories,
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					acctest.TestAccConfigureProvider()
					_, _ = acctest.TestAccGetHarClientWithContext()
				},
				Config: testProjResourceUpstreamHelmRegistry(id, accountId),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "identifier", id),
				),
			},
		},
	})
}

func testAccResourceVirtualHelmRegistry(id string, accId string) string {
	return fmt.Sprintf(`

 resource "harness_platform_har_registry" "test" {
   identifier   = "%[1]s"
   space_ref    = "%[2]s"
   package_type = "HELM"

   config {
    type = "VIRTUAL"
   }
   parent_ref = "%[2]s"
 }
`, id, accId)
}

func testOrgResourceVirtualHelmRegistry(id string, accId string) string {
	return fmt.Sprintf(`
 resource "harness_platform_organization" "test" {
  identifier = "%[1]s_org"
  name = "%[1]s"
 }

 resource "harness_platform_har_registry" "test" {
   identifier   = "%[1]s"
   space_ref    = "%[2]s/${harness_platform_organization.test.identifier}"
   package_type = "HELM"

   config {
    type = "VIRTUAL"
   }
   parent_ref = "%[2]s/${harness_platform_organization.test.identifier}"
 }
`, id, accId)
}

func testProjResourceVirtualHelmRegistry(id string, accId string) string {
	return fmt.Sprintf(`
 resource "harness_platform_organization" "test" {
  identifier = "%[1]s_org"
  name = "%[1]s"
 }

 resource "harness_platform_project" "test" {
  identifier = "%[1]s_project"
  name = "%[1]s"
  org_id = harness_platform_organization.test.id
  color = "#472848"
 }
 resource "harness_platform_har_registry" "test" {
   identifier   = "%[1]s"
   space_ref    = "%[2]s/${harness_platform_organization.test.identifier}/${harness_platform_project.test.identifier}"
   package_type = "HELM"

   config {
    type = "VIRTUAL"
   }
   parent_ref = "%[2]s/${harness_platform_organization.test.identifier}/${harness_platform_project.test.identifier}"
 }
`, id, accId)
}

func testAccResourceUpstreamHelmRegistry(id string, accId string) string {
	return fmt.Sprintf(`
 resource "harness_platform_har_registry" "test" {
   identifier   = "%[1]s"
   space_ref    = "%[2]s"
   package_type = "HELM"

   config {
		type = "UPSTREAM"
		auth_type = "UserPassword"
		source = "Custom"
		url = "https://har-registry.default.svc.cluster.local"
		auth {
			auth_type = "UserPassword"
			user_name = "username"
			secret_identifier = "Secret_Token"
			secret_space_path = "%[2]s"
		}
   }
   parent_ref = "%[2]s"
 }
`, id, accId)
}

func testOrgResourceUpstreamHelmRegistry(id string, accId string) string {
	return fmt.Sprintf(`
 resource "harness_platform_organization" "test" {
  identifier = "%[1]s_org"
  name = "%[1]s"
 }

 resource "harness_platform_har_registry" "test" {
   identifier   = "%[1]s"
   space_ref    = "%[2]s/${harness_platform_organization.test.identifier}"
   package_type = "HELM"

   config {
		type = "UPSTREAM"
		auth_type = "UserPassword"
		source = "Custom"
		url = "https://har-registry.default.svc.cluster.local"
		auth {
			auth_type = "UserPassword"
			user_name = "username"
			secret_identifier = "Secret_Token"
			secret_space_path = "%[2]s"
		}
   }
   parent_ref = "%[2]s/${harness_platform_organization.test.identifier}"
 }
`, id, accId)
}

func testProjResourceUpstreamHelmRegistry(id string, accId string) string {
	return fmt.Sprintf(`
 resource "harness_platform_organization" "test" {
  identifier = "%[1]s_org"
  name = "%[1]s"
 }

 resource "harness_platform_project" "test" {
  identifier = "%[1]s_project"
  name = "%[1]s"
  org_id = harness_platform_organization.test.id
  color = "#472848"
 }
 resource "harness_platform_har_registry" "test" {
   identifier   = "%[1]s"
   space_ref    = "%[2]s/${harness_platform_organization.test.identifier}/${harness_platform_project.test.identifier}"
   package_type = "HELM"

   config {
		type = "UPSTREAM"
		auth_type = "UserPassword"
		source = "Custom"
		url = "https://har-registry.default.svc.cluster.local"
		auth {
			auth_type = "UserPassword"
			user_name = "username"
			secret_identifier = "Secret_Token"
			secret_space_path = "%[2]s"
		}
   }
   parent_ref = "%[2]s/${harness_platform_organization.test.identifier}/${harness_platform_project.test.identifier}"
 }
`, id, accId)
}

// Upstream Docker Registry --- Anonymous
func TestAccResourceUpstreamDockerAnonymousRegistry(t *testing.T) {
	id := fmt.Sprintf("tf_auto_upstream_docker_registry")
	resourceName := "harness_platform_har_registry.test"
	accountId := os.Getenv("HARNESS_ACCOUNT_ID")

	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { acctest.TestAccPreCheck(t) },
		ProviderFactories: acctest.ProviderFactories,
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					acctest.TestAccConfigureProvider()
					_, _ = acctest.TestAccGetHarClientWithContext()
				},
				Config:             testAccResourceUpstreamDockerAnonymousRegistry(id, accountId),
				ExpectNonEmptyPlan: true,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "identifier", id),
				),
			},
		},
	})
}
func TestOrgResourceUpstreamDockerAnonymousRegistry(t *testing.T) {
	id := fmt.Sprintf("tf_auto_upstream_docker_registry")
	resourceName := "harness_platform_har_registry.test"
	accountId := os.Getenv("HARNESS_ACCOUNT_ID")

	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { acctest.TestAccPreCheck(t) },
		ProviderFactories: acctest.ProviderFactories,
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					acctest.TestAccConfigureProvider()
					_, _ = acctest.TestAccGetHarClientWithContext()
				},
				Config:             testOrgResourceUpstreamDockerAnonymousRegistry(id, accountId),
				ExpectNonEmptyPlan: true,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "identifier", id),
				),
			},
		},
	})
}
func TestProjectResourceUpstreamDockerAnonymousRegistry(t *testing.T) {
	id := fmt.Sprintf("tf_auto_upstream_docker_registry")
	resourceName := "harness_platform_har_registry.test"
	accountId := os.Getenv("HARNESS_ACCOUNT_ID")

	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { acctest.TestAccPreCheck(t) },
		ProviderFactories: acctest.ProviderFactories,
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					acctest.TestAccConfigureProvider()
					_, _ = acctest.TestAccGetHarClientWithContext()
				},
				Config:             testProjResourceUpstreamDockerAnonymousRegistry(id, accountId),
				ExpectNonEmptyPlan: true,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "identifier", id),
				),
			},
		},
	})
}

func testAccResourceUpstreamDockerAnonymousRegistry(id string, accId string) string {
	return fmt.Sprintf(`
 resource "harness_platform_har_registry" "test" {
   identifier   = "%[1]s"
   space_ref    = "%[2]s"
   package_type = "DOCKER"

   config {
		type = "UPSTREAM"
		auth_type = "Anonymous"
		source = "Dockerhub"
		auth {
			auth_type = "Anonymous"
		}
   }
   parent_ref = "%[2]s"
 }
`, id, accId)
}

func testOrgResourceUpstreamDockerAnonymousRegistry(id string, accId string) string {
	return fmt.Sprintf(`
 resource "harness_platform_organization" "test" {
  identifier = "%[1]s_org"
  name = "%[1]s"
 }

 resource "harness_platform_har_registry" "test" {
   identifier   = "%[1]s"
   space_ref    = "%[2]s/${harness_platform_organization.test.identifier}"
   package_type = "DOCKER"

   config {
		type = "UPSTREAM"
		auth_type = "UserPassword"
		source = "Dockerhub"
		auth {
			auth_type = "Anonymous"
		}
   }
   parent_ref = "%[2]s/${harness_platform_organization.test.identifier}"
 }
`, id, accId)
}

func testProjResourceUpstreamDockerAnonymousRegistry(id string, accId string) string {
	return fmt.Sprintf(`
 resource "harness_platform_organization" "test" {
  identifier = "%[1]s_org"
  name = "%[1]s"
 }

 resource "harness_platform_project" "test" {
  identifier = "%[1]s_project"
  name = "%[1]s"
  org_id = harness_platform_organization.test.id
  color = "#472848"
 }
 resource "harness_platform_har_registry" "test" {
   identifier   = "%[1]s"
   space_ref    = "%[2]s/${harness_platform_organization.test.identifier}/${harness_platform_project.test.identifier}"
   package_type = "DOCKER"

   config {
		type = "UPSTREAM"
		auth_type = "UserPassword"
		source = "Dockerhub"
		auth {
			auth_type = "Anonymous"
		}
   }
   parent_ref = "%[2]s/${harness_platform_organization.test.identifier}/${harness_platform_project.test.identifier}"
 }
`, id, accId)
}

// Upstream Helm Registry --- Anonymous
func TestAccResourceUpstreamHelmAnonymousRegistry(t *testing.T) {
	id := fmt.Sprintf("tf_auto_upstream_helm_registry")
	resourceName := "harness_platform_har_registry.test"
	accountId := os.Getenv("HARNESS_ACCOUNT_ID")

	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { acctest.TestAccPreCheck(t) },
		ProviderFactories: acctest.ProviderFactories,
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					acctest.TestAccConfigureProvider()
					_, _ = acctest.TestAccGetHarClientWithContext()
				},
				Config:             testAccResourceUpstreamHelmAnonymousRegistry(id, accountId),
				ExpectNonEmptyPlan: true,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "identifier", id),
				),
			},
		},
	})
}
func TestOrgResourceUpstreamHelmAnonymousRegistry(t *testing.T) {
	id := fmt.Sprintf("tf_auto_upstream_helm_registry")
	resourceName := "harness_platform_har_registry.test"
	accountId := os.Getenv("HARNESS_ACCOUNT_ID")

	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { acctest.TestAccPreCheck(t) },
		ProviderFactories: acctest.ProviderFactories,
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					acctest.TestAccConfigureProvider()
					_, _ = acctest.TestAccGetHarClientWithContext()
				},
				Config:             testOrgResourceUpstreamHelmAnonymousRegistry(id, accountId),
				ExpectNonEmptyPlan: true,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "identifier", id),
				),
			},
		},
	})
}
func TestProjectResourceUpstreamHelmAnonymousRegistry(t *testing.T) {
	id := fmt.Sprintf("tf_auto_upstream_helm_registry")
	resourceName := "harness_platform_har_registry.test"
	accountId := os.Getenv("HARNESS_ACCOUNT_ID")

	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { acctest.TestAccPreCheck(t) },
		ProviderFactories: acctest.ProviderFactories,
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					acctest.TestAccConfigureProvider()
					_, _ = acctest.TestAccGetHarClientWithContext()
				},
				Config:             testProjResourceUpstreamHelmAnonymousRegistry(id, accountId),
				ExpectNonEmptyPlan: true,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "identifier", id),
				),
			},
		},
	})
}

func testAccResourceUpstreamHelmAnonymousRegistry(id string, accId string) string {
	return fmt.Sprintf(`
 resource "harness_platform_har_registry" "test" {
   identifier   = "%[1]s"
   space_ref    = "%[2]s"
   package_type = "HELM"

   config {
		type = "UPSTREAM"
		auth_type = "UserPassword"
		source = "Custom"
		url = "https://har-registry.default.svc.cluster.local"
		auth {
			auth_type = "Anonymous"
		}
   }
   parent_ref = "%[2]s"
 }
`, id, accId)
}

func testOrgResourceUpstreamHelmAnonymousRegistry(id string, accId string) string {
	return fmt.Sprintf(`
 resource "harness_platform_organization" "test" {
  identifier = "%[1]s_org"
  name = "%[1]s"
 }

 resource "harness_platform_har_registry" "test" {
   identifier   = "%[1]s"
   space_ref    = "%[2]s/${harness_platform_organization.test.identifier}"
   package_type = "HELM"

   config {
		type = "UPSTREAM"
		auth_type = "UserPassword"
		source = "Custom"
		url = "https://har-registry.default.svc.cluster.local"
		auth {
			auth_type = "Anonymous"
		}
   }
   parent_ref = "%[2]s/${harness_platform_organization.test.identifier}"
 }
`, id, accId)
}

func testProjResourceUpstreamHelmAnonymousRegistry(id string, accId string) string {
	return fmt.Sprintf(`
 resource "harness_platform_organization" "test" {
  identifier = "%[1]s_org"
  name = "%[1]s"
 }

 resource "harness_platform_project" "test" {
  identifier = "%[1]s_project"
  name = "%[1]s"
  org_id = harness_platform_organization.test.id
  color = "#472848"
 }
 resource "harness_platform_har_registry" "test" {
   identifier   = "%[1]s"
   space_ref    = "%[2]s/${harness_platform_organization.test.identifier}/${harness_platform_project.test.identifier}"
   package_type = "HELM"

   config {
		type = "UPSTREAM"
		auth_type = "UserPassword"
		source = "Custom"
		url = "https://har-registry.default.svc.cluster.local"
		auth {
			auth_type = "Anonymous"
		}
   }
   parent_ref = "%[2]s/${harness_platform_organization.test.identifier}/${harness_platform_project.test.identifier}"
 }
`, id, accId)
}
