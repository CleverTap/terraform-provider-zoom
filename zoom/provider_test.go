package zoom

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"log"
	"os"
	"testing"
)

var testAccProviders map[string]*schema.Provider
var testAccProvider *schema.Provider

func init() {
	os.Setenv("ZOOM_API_SECRET", "DEMO_VALUE")
	os.Setenv("ZOOM_API_KEY", "DEMO_VALUE")
	testAccProvider = Provider()
	testAccProviders = map[string]*schema.Provider{
		"zoom": testAccProvider,
	}
}

func TestProvider(t *testing.T) {
	if err := Provider().InternalValidate(); err != nil {
		log.Println("[ERROR]: ", err)
		t.Fatalf("err: %s", err)
	}
}

func TestProvider_impl(t *testing.T) {
	var _ *schema.Provider = Provider()
}

func testAccPreCheck(t *testing.T) {
	if v := os.Getenv("ZOOM_API_SECRET"); v == "" {
		t.Fatal("ZOOM_API_SECRET must be set for acceptance tests")
	}
	if v := os.Getenv("ZOOM_API_KEY"); v == "" {
		t.Fatal("ZOOM_API_KEY must be set for acceptance tests")
	}
}
