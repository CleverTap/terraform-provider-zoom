package zoom

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccItem_Basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckItemBasic(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("zoom_user.user1", "email", "user@gmail.com"),
					resource.TestCheckResourceAttr("zoom_user.user1", "first_name", "FirstName"),
					resource.TestCheckResourceAttr("zoom_user.user1", "last_name", "LastName"),
				),
			},
		},
	})
}

func testAccCheckItemBasic() string {
	return fmt.Sprintf(`
resource "zoom_user" "user1" {
      email        = "user@gmail.com"
      first_name   = "FirstName"
      last_name    = "LastName"
      license_type =  1
    }
`)
}

func TestAccItem_Update(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckItemUpdatePre(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"zoom_user.user1", "email", "user@gmail.com"),
					resource.TestCheckResourceAttr(
						"zoom_user.user1", "first_name", "FirstName"),
					resource.TestCheckResourceAttr(
						"zoom_user.user1", "last_name", "LastName"),
				),
			},
			{
				Config: testAccCheckItemUpdatePost(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"zoom_user.user1", "email", "user@gmail.com"),
					resource.TestCheckResourceAttr(
						"zoom_user.user1", "first_name", "FirstName"),
					resource.TestCheckResourceAttr(
						"zoom_user.user1", "last_name", "LastName"),
				),
			},
		},
	})
}

func testAccCheckItemUpdatePre() string {
	return fmt.Sprintf(`
resource "zoom_user" "user1" {
	email        = "user@gmail.com"
	first_name   = "FirstName"
	last_name    = "LastName"
	status       = "activate"
	license_type =  1
}
`)
}

func testAccCheckItemUpdatePost() string {
	return fmt.Sprintf(`
resource "zoom_user" "user1" {
	email        = "user@gmail.com"
	first_name   = "FirstName"
	last_name    = "LastName"
	status       = "activate"
	license_type =  1
}
`)
}
