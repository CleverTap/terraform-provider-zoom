package zoom

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"terraform-provider-zoom/client"
	"time"
)

func dataSourceUser() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceUserRead,
		Schema: map[string]*schema.Schema{
			"email": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validateEmail,
			},
			"first_name": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"last_name": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"status": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"license_type": &schema.Schema{
				Type:     schema.TypeInt,
				Computed: true,
			},
			"pmi": &schema.Schema{
				Type:     schema.TypeInt,
				Computed: true,
			},
			"use_pmi": &schema.Schema{
				Type:     schema.TypeBool,
				Computed: true,
			},
			"timezone": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"language": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"vanity_name": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"host_key": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"cms_user_id": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"company": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"group_id": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"manager": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"pronouns": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"pronouns_option": &schema.Schema{
				Type:     schema.TypeInt,
				Computed: true,
			},
			"phone_numbers": &schema.Schema{
				Type:     schema.TypeSet,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"country": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						"code": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						"number": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						"label": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"role_name": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"department": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"job_title": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"location": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceUserRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	apiClient := m.(*client.Client)
	var email string
	if v, ok := d.GetOk("email"); ok {
		email = v.(string)
	}
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
		d.SetId(email)
		return nil
	})
	if retryErr != nil {
		time.Sleep(2 * time.Second)
		return diag.FromErr(retryErr)
	}
	return diags
}
