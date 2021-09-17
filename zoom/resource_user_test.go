package zoom

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccUser_Basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckUserBasic(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("zoom_user.user1", "email", "user@gmail.com"),
					resource.TestCheckResourceAttr("zoom_user.user1", "first_name", "FirstName"),
					resource.TestCheckResourceAttr("zoom_user.user1", "last_name", "LastName"),
				),
			},
		},
	})
}

func testAccCheckUserBasic() string {
	return fmt.Sprintf(`
		resource "zoom_user" "user1" {
			email        = "user@gmail.com"
			first_name   = "FirstName"
			last_name    = "LastName"
			license_type =  1
		}
	`)
}

func TestAccUser_Update(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckUserUpdatePre(),
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
				Config: testAccCheckUserUpdatePost(),
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

func testAccCheckUserUpdatePre() string {
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

func testAccCheckUserUpdatePost() string {
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
