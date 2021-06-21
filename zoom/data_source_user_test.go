package zoom

import (
	"fmt"
	"testing"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccUserDataSource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccUserDataSourceConfig(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(
						"data.zoom_user.user5", "id", "ui17ec38@iiitsurat.ac.in"),
				),
			},
		},
	})
}

func testAccUserDataSourceConfig() string {
	return fmt.Sprintf(`	  
	resource "zoom_user" "user5" {
		email        = "ui17co15@iitsurat.ac.in"
		first_name   = "ekansh"
		last_name    = "singh"
		license_type =  1
	  }
	data "zoom_user" "user5" {
		id = "ui17ec38@iiitsurat.ac.in"
	}
	`)
}
