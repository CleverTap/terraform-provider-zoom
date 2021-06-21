package zoom

import(
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccItem_Basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckItemBasic(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("zoom_user.user1", "email", "tapendrakmr639@gmail.com"),
					resource.TestCheckResourceAttr("zoom_user.user1", "first_name", "Ekansh"),
					resource.TestCheckResourceAttr("zoom_user.user1", "last_name", "Singh"),
				),
			},
		},
	})
}

func testAccCheckItemBasic() string {
	return fmt.Sprintf(`
resource "zoom_user" "user1" {
      email        = "tapendrakmr639@gmail.com"
      first_name   = "Ekansh"
      last_name    = "Singh"
      license_type =  1
    }
`)
}

func TestAccItem_Update(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckItemUpdatePre(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"zoom_user.user1", "email", "ekansh6336@gmail.com"),
					resource.TestCheckResourceAttr(
						"zoom_user.user1", "first_name", "Ekansh"),
					resource.TestCheckResourceAttr(
						"zoom_user.user1", "last_name", "Singh"),	
				),
			},
			{
				Config: testAccCheckItemUpdatePost(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"zoom_user.user1", "email", "ekansh6336@gmail.com"),
					resource.TestCheckResourceAttr(
						"zoom_user.user1", "first_name", "Ekansh"),
					resource.TestCheckResourceAttr(
						"zoom_user.user1", "last_name", "kumar"),
				),
			},
		},
	})
}

func testAccCheckItemUpdatePre() string {
	return fmt.Sprintf(`
resource "zoom_user" "user1" {
	email        = "ekansh6336@gmail.com"
	first_name   = "Ekansh"
	last_name    = "Singh"
	status       = "activate"
	license_type =  1
}
`)
}

func testAccCheckItemUpdatePost() string {
	return fmt.Sprintf(`
resource "zoom_user" "user1" {
	email        = "ekansh6336@gmail.com"
	first_name   = "Ekansh"
	last_name    = "kumar"
	status       = "activate"
	license_type =  1
}
`)
}
