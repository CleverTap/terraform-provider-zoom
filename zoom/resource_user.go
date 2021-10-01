package zoom

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"regexp"
	"strings"
	"terraform-provider-zoom/client"
	"time"
)

func validateEmail(v interface{}, k string) (warns []string, errs []error) {
	value := v.(string)
	var emailRegex = regexp.MustCompile(`^[A-Za-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`)
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
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validateEmail,
			},
			"first_name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"last_name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"status": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"license_type": &schema.Schema{
				Type:     schema.TypeInt,
				Required: true,
			},
			"pmi": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"use_pmi": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"timezone": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"language": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"vanity_name": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"host_key": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"cms_user_id": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"company": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"group_id": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"manager": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"pronouns": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"pronouns_option": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"phone_numbers": &schema.Schema{
				Type:     schema.TypeSet,
				Optional: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"country": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"code": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"number": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"label": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
					},
				},
			},
			"role_name": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"department": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"job_title": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"location": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
		},
	}

}

func resourceUserCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	apiClient := m.(*client.Client)
	user := client.User{}
	if v, ok := d.GetOk("email"); ok {
		user.Email = v.(string)
	}
	if v, ok := d.GetOk("first_name"); ok {
		user.FirstName = v.(string)
	}
	if v, ok := d.GetOk("last_name"); ok {
		user.LastName = v.(string)
	}
	if v, ok := d.GetOk("license_type"); ok {
		user.Type = v.(int)
	}
	retryErr := resource.Retry(time.Duration(apiClient.TimeoutMinutes)*time.Minute, func() *resource.RetryError {
		if err := apiClient.NewUser(&user); err != nil {
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
	d.SetId(user.Email)
	resourceUserRead(ctx, d, m)
	return diags
}

func resourceUserRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	apiClient := m.(*client.Client)
	email := d.Id()
	retryErr := resource.Retry(time.Duration(apiClient.TimeoutMinutes)*time.Minute, func() *resource.RetryError {
		user, err := apiClient.GetUser(email)
		if err != nil {
			if apiClient.IsRetry(err) {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		d.Set("email", user.Email)
		d.Set("first_name", user.FirstName)
		d.Set("last_name", user.LastName)
		d.Set("license_type", user.Type)
		d.Set("pmi", user.Pmi)
		d.Set("use_pmi", user.UsePmi)
		d.Set("timezone", user.Timezone)
		d.Set("vanity_name", user.VanityName)
		d.Set("host_key", user.HostKey)
		d.Set("cms_user_id", user.CmsUserId)
		d.Set("language", user.Language)
		d.Set("company", user.Company)
		d.Set("group_id", user.GroupId)
		d.Set("manager", user.Manager)
		d.Set("pronouns", user.Pronouns)
		d.Set("pronouns_option", user.PronounsOption)
		d.Set("status", user.Status)
		var phoneNumbers []interface{}
		for _, v := range user.PhoneNumbers {
			phoneNumber := make(map[string]string)
			phoneNumber["country"] = v.Country
			phoneNumber["code"] = v.Code
			phoneNumber["number"] = v.Number
			phoneNumber["label"] = v.Label
			phoneNumbers = append(phoneNumbers, phoneNumber)
		}
		d.Set("phone_numbers", phoneNumbers)
		d.Set("role_name", user.RoleName)
		d.Set("department", user.Department)
		d.Set("job_title", user.JobTitle)
		d.Set("location", user.Location)
		return nil
	})
	if retryErr != nil {
		time.Sleep(2 * time.Second)
		if strings.Contains(retryErr.Error(), "404") == true {
			d.SetId("")
			return diags
		}
		return diag.FromErr(retryErr)
	}
	return diags
}

func resourceUserUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var _ diag.Diagnostics
	apiClient := m.(*client.Client)
	var diags diag.Diagnostics
	if d.HasChange("email") {
		if v, ok := d.GetOk("email"); ok {
			newEmail := v.(string)
			retryErr := resource.Retry(time.Duration(apiClient.TimeoutMinutes)*time.Minute, func() *resource.RetryError {
				if err := apiClient.ChangeEmail(d.Id(), newEmail); err != nil {
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
			d.SetId(newEmail)
		}
	}
	if d.HasChange("status") {
		if v, ok := d.GetOk("status"); ok {
			status := v.(string)
			var action string
			if status == "active" {
				action = "activate"
			} else if status == "inactive" {
				action = "deactivate"
			}
			retryErr := resource.Retry(time.Duration(apiClient.TimeoutMinutes)*time.Minute, func() *resource.RetryError {
				if err := apiClient.ChangeUserStatus(d.Id(), action); err != nil {
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
		}
	} else {
		user := client.User{}
		if d.HasChange("first_name") {
			if v, ok := d.GetOk("first_name"); ok {
				user.FirstName = v.(string)
			}
		}
		if d.HasChange("last_name") {
			if v, ok := d.GetOk("last_name"); ok {
				user.LastName = v.(string)
			}
		}
		if d.HasChange("license_type") {
			if v, ok := d.GetOk("license_type"); ok {
				user.Type = v.(int)
			}
		}
		if d.HasChange("pmi") {
			if v, ok := d.GetOk("pmi"); ok {
				user.Pmi = v.(int)
			}
		}
		if d.HasChange("use_pmi") {
			if v, ok := d.GetOk("use_pmi"); ok {
				user.UsePmi = v.(*bool)
			} else {
				user.UsePmi = v.(*bool)
			}
		}
		if d.HasChange("timezone") {
			if v, ok := d.GetOk("timezone"); ok {
				user.Timezone = v.(string)
			}
		}
		if d.HasChange("language") {
			if v, ok := d.GetOk("language"); ok {
				user.Language = v.(string)
			}
		}
		if d.HasChange("vanity_name") {
			if v, ok := d.GetOk("vanity_name"); ok {
				user.VanityName = v.(string)
			}
		}
		if d.HasChange("host_key") {
			if v, ok := d.GetOk("host_key"); ok {
				user.HostKey = v.(string)
			}
		}
		if d.HasChange("cms_user_id") {
			if v, ok := d.GetOk("cms_user_id"); ok {
				user.CmsUserId = v.(string)
			}
		}
		if d.HasChange("company") {
			if v, ok := d.GetOk("company"); ok {
				user.Company = v.(string)
			}
		}
		if d.HasChange("group_id") {
			if v, ok := d.GetOk("group_id"); ok {
				user.GroupId = v.(string)
			}
		}
		if d.HasChange("manager") {
			if v, ok := d.GetOk("manager"); ok {
				user.Manager = v.(string)
			}
		}
		if d.HasChange("pronouns") {
			if v, ok := d.GetOk("pronouns"); ok {
				user.Pronouns = v.(string)
			}
		}
		if d.HasChange("pronouns_option") {
			if v, ok := d.GetOk("pronouns_option"); ok {
				user.PronounsOption = v.(int)
			}
		}
		if d.HasChange("department") {
			if v, ok := d.GetOk("department"); ok {
				user.Department = v.(string)
			}
		}
		if d.HasChange("job_title") {
			if v, ok := d.GetOk("job_title"); ok {
				user.JobTitle = v.(string)
			}
		}
		if d.HasChange("location") {
			if v, ok := d.GetOk("location"); ok {
				user.Location = v.(string)
			}
		}
		if d.HasChange("phone_numbers") {
			if v, ok := d.GetOk("phone_numbers"); ok {
				phoneNumbers := v.(*schema.Set).List()
				var phoneNumbersList []client.PhoneNumber
				for _, v := range phoneNumbers {
					phoneNumber := v.(map[string]interface{})
					phoneNumberStruct := client.PhoneNumber{
						Country: phoneNumber["country"].(string),
						Code:    phoneNumber["code"].(string),
						Number:  phoneNumber["number"].(string),
						Label:   phoneNumber["label"].(string),
					}
					phoneNumbersList = append(phoneNumbersList, phoneNumberStruct)
				}
				user.PhoneNumbers = phoneNumbersList
			}
		}
		retryErr := resource.Retry(time.Duration(apiClient.TimeoutMinutes)*time.Minute, func() *resource.RetryError {
			if err := apiClient.UpdateUser(d.Id(), &user); err != nil {
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
	}
	return diags
}

func resourceUserDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	apiClient := m.(*client.Client)
	email := d.Id()
	retryErr := resource.Retry(time.Duration(apiClient.TimeoutMinutes)*time.Minute, func() *resource.RetryError {
		if err := apiClient.DeleteUser(email, d.Get("status").(string)); err != nil {
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
	d.SetId("")
	return diags
}

func resourceUserImporter(ctx context.Context, d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	apiClient := m.(*client.Client)
	email := d.Id()
	retryErr := resource.Retry(time.Duration(apiClient.TimeoutMinutes)*time.Minute, func() *resource.RetryError {
		user, err := apiClient.GetUser(email)
		if err != nil {
			if apiClient.IsRetry(err) {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		d.Set("email", user.Email)
		d.Set("first_name", user.FirstName)
		d.Set("last_name", user.LastName)
		d.Set("license_type", user.Type)
		d.Set("pmi", user.Pmi)
		d.Set("use_pmi", user.UsePmi)
		d.Set("timezone", user.Timezone)
		d.Set("vanity_name", user.VanityName)
		d.Set("host_key", user.HostKey)
		d.Set("cms_user_id", user.CmsUserId)
		d.Set("language", user.Language)
		d.Set("company", user.Company)
		d.Set("group_id", user.GroupId)
		d.Set("manager", user.Manager)
		d.Set("pronouns", user.Pronouns)
		d.Set("pronouns_option", user.PronounsOption)
		d.Set("status", user.Status)
		var phoneNumbers []interface{}
		for _, v := range user.PhoneNumbers {
			phoneNumber := make(map[string]string)
			phoneNumber["country"] = v.Country
			phoneNumber["code"] = v.Code
			phoneNumber["number"] = v.Number
			phoneNumber["label"] = v.Label
			phoneNumbers = append(phoneNumbers, phoneNumber)
		}
		d.Set("role_name", user.RoleName)
		d.Set("department", user.Department)
		d.Set("job_title", user.JobTitle)
		d.Set("location", user.Location)
		return nil
	})
	if retryErr != nil {
		time.Sleep(2 * time.Second)
		return nil, retryErr
	}
	return []*schema.ResourceData{d}, nil
}
