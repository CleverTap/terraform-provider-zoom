# Resource User testing

<https://www.terraform.io/docs/extend/testing/acceptance-tests/index.html>

#### BASIC COMMANDS TO RUN TEST FILE

1-go test <br/>
2-go test -cover (to get idea about how much percentage of your code is tested) <br />
3-go test ./... (to run all the test file in a particular folder) <br />

//In our case we are doing acceptance testing <br />

### <strong> STEPS</strong>

1-make TF_ACC = true (set environment variable,this is to run the acceptance testing) <br />

2-Hashicorp has provider some inbuilt packages which we can use to implemet our testing ie. resource("github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource") <br />

3-we set up a resource.Test and provide it with the following: <br />
PreCheck,Providers,CheckDestroy,Steps <br />

4- In each steps, we can provide couple of things <br />
Config,Check <br />

### DIFFERENT TESTING FUNCTION

<strong>1. testAccCheckItemDestroy</strong>

Runs at the end of our test after all the steps have been run and the resources have been destroyed <br />

<strong>2. testAccCheckItemBasic </strong>

Returns a string containing some Terraform code <br />

<strong>3. testAccCheckExampleItemExists </strong>

Checks that the user we asked the test to create was actually created <br />

<strong>4. TestCheckResourceAttr </strong>

It verifies that the resource attributes matches <br />

<strong>5. TestAccItem_Basic </strong>

It calls the destroy and exist function to test delete and create finction <br />

<strong>6. TestAccItem_Update</strong>

It is the testing function for update <br />
