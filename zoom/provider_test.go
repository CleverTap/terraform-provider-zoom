package zoom

import(
	"os"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"testing"
	"log"
	"io/ioutil"

)

var testAccProviders map[string]*schema.Provider
var testAccProvider *schema.Provider

func init() {
	file, err := os.Open("../acctoken.txt")
    if err != nil {
        log.Fatal(err)
    }
	token, err := ioutil.ReadAll(file)
	os.Setenv("ZOOM_TOKEN", string(token))
	testAccProvider = Provider()
	testAccProviders = map[string]*schema.Provider{
		"zoom": testAccProvider,
	}
}

func TestProvider(t *testing.T) {
	if err := Provider().InternalValidate(); err != nil {
		log.Println("[ERROR]: ",err)
		t.Fatalf("err: %s", err)
	}
}

func TestProvider_impl(t *testing.T)  {
	var _ *schema.Provider = Provider()
}

func testAccPreCheck(t *testing.T) {
	if v := os.Getenv("ZOOM_TOKEN"); v == "" {
		t.Fatal("ZOOM_TOKEN must be set for acceptance tests")
	}
}
