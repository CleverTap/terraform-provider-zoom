package zoom

import (
	"terraform-provider-zoom/client"
	"strings"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"fmt"
	"regexp"
	"log"
	"time"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
)

func validateName(v interface{}, k string) (ws []string, es []error) {
	var errs []error
	var warns []string
	value, ok := v.(string)
	if !ok {
		errs = append(errs, fmt.Errorf("Expected name to be string"))
		return warns, errs
	}
	nameRegex := regexp.MustCompile("^[A-Za-z]\\w{5,29}$")

	if !(nameRegex.MatchString(k)) {
		errs = append(errs, fmt.Errorf("Expected name is not valid .Got %s", value))
		return warns, errs
	}
	return
}

func validateEmail(v interface{}, k string) (ws []string, es []error) {
	var errs []error
	var warns []string
	value := v.(string)

	var emailRegex = regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`)

	if !(emailRegex.MatchString(value)) {
		errs = append(errs, fmt.Errorf("Expected EmailId is not valid  %s", k))
		return warns, errs
	}
	return
}

func resourceUser() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceUserCreate,
		ReadContext:   resourceUserRead,
		UpdateContext: resourceUserUpdate,
		DeleteContext: resourceUserDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceUserImporter,
		},
		Schema: map[string]*schema.Schema{

			"email": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				ValidateFunc: validateEmail,
			},
			"first_name": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				ValidateFunc: validateName,
			},
			"last_name": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				ValidateFunc: validateName,
			},
			"status": &schema.Schema{
				Type:        schema.TypeString,
				Optional :   true,
				Computed:    true,
			},
			"id": &schema.Schema{
				Type:        schema.TypeString,
				Computed :   true,
			},
			"license_type": &schema.Schema{
				Type:        schema.TypeInt,
				Required:    true,
			},
			"pmi": &schema.Schema{
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
			},
			"role_name": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
			},
			"department": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
			},
			"job_title": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
			},
			"location": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
			},
		},
	}

}

func resourceUserCreate(ctx context.Context,d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	apiClient := m.(*client.Client)
	user := client.User{
		Email:   d.Get("email").(string),
		FirstName: d.Get("first_name").(string),
		LastName:  d.Get("last_name").(string),
		Type:      d.Get("license_type").(int),
	}

	var err error
	retryErr := resource.Retry(2*time.Minute, func() *resource.RetryError {
		if err = apiClient.NewItem(&user); err != nil {
			if apiClient.IsRetry(err) {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	if retryErr != nil {
		time.Sleep(2 * time.Second)
		return diag.FromErr(retryErr)
	}

	if err != nil {
		log.Println("[ERROR]: ",err)
		return diag.FromErr(err)
	}
	d.SetId(user.Email)
	resourceUserRead(ctx,d,m)
	return diags
}

func resourceUserRead(ctx context.Context,d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	apiClient := m.(*client.Client)
	userId := d.Id()

	retryErr := resource.Retry(2*time.Minute, func() *resource.RetryError {
		user, err := apiClient.GetItem(userId)
		if err != nil {
			if apiClient.IsRetry(err) {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		if len(user.Email) > 0{
			d.SetId(user.Email)
			d.Set("email", user.Email)
			d.Set("first_name", user.FirstName)
			d.Set("last_name", user.LastName)
			d.Set("license_type", user.Type)
			d.Set("pmi",user.Pmi)
			d.Set("status",user.Status)
			d.Set("role_name", user.RoleName)
			d.Set("department",user.Department)
			d.Set("job_title", user.JobTitle)
			d.Set("location", user.Location)
		}
		return nil
	})
	if retryErr!=nil {
		if strings.Contains(retryErr.Error(), "User does not exist")==true {
			d.SetId("")
			return diags
		}
		return diag.FromErr(retryErr)
	}
	return diags
}

func resourceUserUpdate(ctx context.Context,d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var _ diag.Diagnostics
	apiClient := m.(*client.Client)
	var diags diag.Diagnostics
	if d.HasChange("email") {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "User not allowed to change email",
			Detail:   "User not allowed to change email",
		})
		return diags
	}
	user := client.User{
		Email:   d.Get("email").(string),
		FirstName: d.Get("first_name").(string),
		LastName:  d.Get("last_name").(string),
		Type:  d.Get("license_type").(int),
		Department:  d.Get("department").(string),
		JobTitle:  d.Get("job_title").(string),
		Location:  d.Get("location").(string),
	}

	var err_deactivate error
	var user_status string
	status := d.Get("status").(string)
	if(status=="active"){
		user_status = "activate"
	}else if(status=="inactive"){
		user_status = "deactivate"
	}
	retryErrDeac := resource.Retry(2*time.Minute, func() *resource.RetryError {
		if err_deactivate = apiClient.DeactivateUser(user.Email, user_status); err_deactivate != nil {
			if apiClient.IsRetry(err_deactivate) {
				return resource.RetryableError(err_deactivate)
			}
			return resource.NonRetryableError(err_deactivate)
		}
		return nil
	})
	if retryErrDeac != nil {
		time.Sleep(2 * time.Second)
		return diag.FromErr(retryErrDeac)
	}
	if err_deactivate != nil {
		return diag.FromErr(err_deactivate)
	}

	var err error
	retryErr := resource.Retry(2*time.Minute, func() *resource.RetryError {
		if err = apiClient.UpdateItem(&user); err != nil {
			if apiClient.IsRetry(err) {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	if retryErr != nil {
		time.Sleep(2 * time.Second)
		return diag.FromErr(retryErr)
	}
	if err != nil {
		return diag.FromErr(err)
	}
	return diags
}

func resourceUserDelete(ctx context.Context,d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	apiClient := m.(*client.Client)
	userId := d.Id()

	var err error
	retryErr := resource.Retry(2*time.Minute, func() *resource.RetryError {
		if err = apiClient.DeleteItem(userId); err != nil {
			if apiClient.IsRetry(err) {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	if retryErr != nil {
		time.Sleep(2 * time.Second)
		return diag.FromErr(retryErr)
	}
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId("")
	return diags
}

func resourceUserImporter(ctx context.Context, d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	apiClient := m.(*client.Client)
	userId := d.Id()
	user, err := apiClient.GetItem(userId)
	if err!=nil{
		return nil, err
	}
	d.Set("email", user.Email)
		d.Set("first_name", user.FirstName)
		d.Set("last_name", user.LastName)
		d.Set("license_type", user.Type)
		d.Set("pmi",user.Pmi)
		d.Set("status",user.Status)
		d.Set("role_name", user.RoleName)
		d.Set("department",user.Department)
		d.Set("job_title", user.JobTitle)
		d.Set("location", user.Location)
	return []*schema.ResourceData{d}, nil
}
