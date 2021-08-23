package zoom

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccUserDataSource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccUserDataSourceConfig(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(
						"data.zoom_user.user1", "email", "user@gmail.com"),
				),
			},
		},
	})
}

func testAccUserDataSourceConfig() string {
	return fmt.Sprintf(`	  
		resource "zoom_user" "user1" {
			email        = "user@gmail.com"
			first_name   = "FirstName"
			last_name    = "LastName"
			license_type =  1
		}
		data "zoom_user" "user1" {
			email = "user@gmail.com"
		}
	`)
}
