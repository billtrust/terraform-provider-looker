package looker

import (
	"fmt"
	"strings"

	"github.com/Foxtel-DnA/looker-go-sdk/client/connection"

	apiclient "github.com/Foxtel-DnA/looker-go-sdk/client"
	"github.com/Foxtel-DnA/looker-go-sdk/models"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceConnection() *schema.Resource {
	return &schema.Resource{
		Create: resourceConnectionCreate,
		Read:   resourceConnectionRead,
		Update: resourceConnectionUpdate,
		Delete: resourceConnectionDelete,
		Exists: resourceConnectionExists,
		Importer: &schema.ResourceImporter{
			State: resourceConnectionImport,
		},

		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true, // todo: the ID is the name of the connection so if it changes i think it would require a new object be created.  I should verify this
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					if strings.ToLower(old) == strings.ToLower(new) {
						return true
					}
					return false
				},
				ValidateFunc: func(val interface{}, key string) (warns []string, errs []error) {
					v := val.(string)
					if strings.Contains(v, " ") {
						errs = append(errs, fmt.Errorf("%q must not contain any spaces, got: %q", key, v))
					}
					return
				},
			},
			"dialect_name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"host": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"port": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Default:  "443",
			},
			"database": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"username": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"password": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					if d.Id() == "" {
						return false
					}
					// TODO: not sure how to handle this scenario.  The terraform sets the password for the connection on create, but this value is not returned on get
					// TODO: to handle this, DiffSupressFunc always say this field is not changed to handle this.  And if the user changes the field, ForceNew should create a new one
					return true
				},
				Sensitive: true,
			},
			"certificate": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				// Same deal as `password` field above - the certificate
				// is not returned in the API response, so there's no
				// sensible way to diff.
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					if d.Id() == "" {
						return false
					}
					return true
				},
				Sensitive: true,
			},
			"certificate_file_type": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				// Certificate file type also not returned in GET response.
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					if d.Id() == "" {
						return false
					}
					return true
				},
			},
			"schema": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"tmp_db_name": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"jdbc_additional_params": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Default:  "",
			},
			"ssl": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"db_timezone": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"query_timezone": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func resourceConnectionCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(*apiclient.Looker)

	params := connection.NewCreateConnectionParams()
	params.Body = &models.DBConnection{}
	params.Body.Name = d.Get("name").(string)
	params.Body.DialectName = d.Get("dialect_name").(string)
	params.Body.Host = d.Get("host").(string)
	params.Body.Port = d.Get("port").(string)
	params.Body.Database = d.Get("database").(string)
	params.Body.Username = d.Get("username").(string)
	params.Body.Password = d.Get("password").(string)
	params.Body.Certificate = d.Get("certificate").(string)
	params.Body.FileType = d.Get("certificate_file_type").(string)
	params.Body.Schema = d.Get("schema").(string)
	params.Body.TmpDbName = d.Get("tmp_db_name").(string)
	params.Body.JdbcAdditionalParams = d.Get("jdbc_additional_params").(string)
	params.Body.Ssl = d.Get("ssl").(bool)
	params.Body.DbTimezone = d.Get("db_timezone").(string)
	params.Body.QueryTimezone = d.Get("query_timezone").(string)

	result, err := client.Connection.CreateConnection(params)
	if err != nil {
		return err
	}

	d.SetId(result.Payload.Name)

	return resourceConnectionRead(d, m)
}

func resourceConnectionRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*apiclient.Looker)

	params := connection.NewConnectionParams()
	params.ConnectionName = d.Id()

	result, err := client.Connection.Connection(params)
	if err != nil {
		if strings.Contains(err.Error(), "Not found") {
			d.SetId("")
			return nil
		}

		return err
	}

	d.Set("name", result.Payload.Name)
	d.Set("dialect_name", result.Payload.DialectName)
	d.Set("host", result.Payload.Host)
	d.Set("port", result.Payload.Port)
	d.Set("database", result.Payload.Database)
	d.Set("username", result.Payload.Username)
	d.Set("password", result.Payload.Password)
	d.Set("certificate", result.Payload.Certificate)
	d.Set("certificate_file_type", result.Payload.FileType)
	d.Set("schema", result.Payload.Schema)
	d.Set("tmp_db_name", result.Payload.TmpDbName)
	d.Set("jdbc_additional_params", result.Payload.JdbcAdditionalParams)
	d.Set("ssl", result.Payload.Ssl)
	d.Set("db_timezone", result.Payload.DbTimezone)
	d.Set("query_timezone", result.Payload.QueryTimezone)

	return nil
}

func resourceConnectionUpdate(d *schema.ResourceData, m interface{}) error {
	client := m.(*apiclient.Looker)

	params := connection.NewUpdateConnectionParams()
	params.ConnectionName = d.Get("name").(string)
	params.Body = &models.DBConnection{}
	params.Body.Name = d.Get("name").(string)
	params.Body.DialectName = d.Get("dialect_name").(string)
	params.Body.Host = d.Get("host").(string)
	params.Body.Port = d.Get("port").(string)
	params.Body.Database = d.Get("database").(string)
	params.Body.Username = d.Get("username").(string)
	params.Body.Password = d.Get("password").(string)
	params.Body.Certificate = d.Get("certificate").(string)
	params.Body.FileType = d.Get("certificate_file_type").(string)
	params.Body.Schema = d.Get("schema").(string)
	params.Body.TmpDbName = d.Get("tmp_db_name").(string)
	params.Body.JdbcAdditionalParams = d.Get("jdbc_additional_params").(string)
	params.Body.Ssl = d.Get("ssl").(bool)
	params.Body.DbTimezone = d.Get("db_timezone").(string)
	params.Body.QueryTimezone = d.Get("query_timezone").(string)

	_, err := client.Connection.UpdateConnection(params)
	if err != nil {
		return err
	}

	return resourceConnectionRead(d, m)
}

func resourceConnectionDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(*apiclient.Looker)

	params := connection.NewDeleteConnectionParams()
	params.ConnectionName = d.Id()

	_, err := client.Connection.DeleteConnection(params)
	if err != nil {
		return err
	}

	return nil
}

func resourceConnectionExists(d *schema.ResourceData, m interface{}) (b bool, e error) {
	// Exists - This is called to verify a resource still exists. It is called prior to Read,
	// and lowers the burden of Read to be able to assume the resource exists.
	client := m.(*apiclient.Looker)

	params := connection.NewConnectionParams()
	params.ConnectionName = d.Id()

	_, err := client.Connection.Connection(params)
	if err != nil {
		if strings.Contains(err.Error(), "Not found") {
			return false, nil
		}

		return false, err
	}

	return true, nil
}

func resourceConnectionImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	if err := resourceConnectionRead(d, m); err != nil {
		return nil, err
	}
	return []*schema.ResourceData{d}, nil
}
