package looker

import (
	"log"
	"strconv"
	"strings"

	"github.com/billtrust/looker-go-sdk/client/space"

	apiclient "github.com/billtrust/looker-go-sdk/client"
	"github.com/billtrust/looker-go-sdk/models"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceChildSpace() *schema.Resource {
	return &schema.Resource{
		Create: resourceChildSpaceCreate,
		Read:   resourceChildSpaceRead,
		Update: resourceChildSpaceUpdate,
		Delete: resourceChildSpaceDelete,
		Exists: resourceChildSpaceExists,
		Importer: &schema.ResourceImporter{
			State: resourceChildSpaceImport,
		},

		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"parent_id": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"content_metadata_id": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func getChildSpaceByID(d *schema.ResourceData, m interface{}, id int64) (*models.Space, error) {
	client := m.(*apiclient.Looker)

	params := space.NewSpaceParams()
	params.SpaceID = strconv.FormatInt(id, 10)

	result, err := client.Space.Space(params)
	if err != nil {
		log.Printf("[ERROR] Error while getting space by id, %s", err.Error())
		return nil, err
	}

	return result.Payload, nil
}

func resourceChildSpaceCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(*apiclient.Looker)

	parentID, err := getIDFromString(d.Get("parent_id").(string))
	if err != nil {
		return err
	}

	var parentIDstr = strconv.FormatInt(parentID, 10)

	params := space.NewCreateSpaceParams()
	params.Body = &models.CreateSpace{}
	params.Body.Name = d.Get("name").(*string)
	params.Body.ParentID = &parentIDstr

	result, err := client.Space.CreateSpace(params)
	if err != nil {
		return err
	}

	d.SetId(result.Payload.ID)

	return resourceChildSpaceRead(d, m)
}

func resourceChildSpaceRead(d *schema.ResourceData, m interface{}) error {
	ID, err := getIDFromString(d.Id())
	if err != nil {
		return err
	}

	space, err := getChildSpaceByID(d, m, ID)
	if err != nil {
		if strings.Contains(err.Error(), "Not found") {
			d.SetId("")
			return nil
		}
		return err
	}

	d.Set("name", space.Name)
	d.Set("content_metadata_id", getStringFromID(space.ContentMetadataID))
	d.Set("parent_id", space.ParentID)

	return nil
}

func resourceChildSpaceUpdate(d *schema.ResourceData, m interface{}) error {
	client := m.(*apiclient.Looker)

	ID, err := getIDFromString(d.Id())
	if err != nil {
		return err
	}

	parentID, err := getIDFromString(d.Get("parent_id").(string))
	if err != nil {
		return err
	}

	params := space.NewUpdateSpaceParams()
	params.SpaceID = strconv.FormatInt(ID, 10)
	params.Body = &models.UpdateSpace{}
	params.Body.Name = d.Get("name").(string)
	params.Body.ParentID = strconv.FormatInt(parentID, 10)

	_, err = client.Space.UpdateSpace(params)
	if err != nil {
		return err
	}

	return resourceChildSpaceRead(d, m)
}

func resourceChildSpaceDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(*apiclient.Looker)

	ID, err := getIDFromString(d.Id())
	if err != nil {
		return err
	}

	params := space.NewDeleteSpaceParams()
	params.SpaceID = strconv.FormatInt(ID, 10)

	_, err = client.Space.DeleteSpace(params)
	if err != nil {
		return err
	}

	return nil
}

func resourceChildSpaceExists(d *schema.ResourceData, m interface{}) (b bool, e error) {
	// Exists - This is called to verify a resource still exists. It is called prior to Read,
	// and lowers the burden of Read to be able to assume the resource exists.
	ID, err := getIDFromString(d.Id())
	if err != nil {
		return false, err
	}

	_, err = getChildSpaceByID(d, m, ID)
	if err != nil {
		if strings.Contains(err.Error(), "Not found") {
			return false, nil
		}

		return false, err
	}

	return true, nil
}

func resourceChildSpaceImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	if err := resourceChildSpaceRead(d, m); err != nil {
		return nil, err
	}
	return []*schema.ResourceData{d}, nil
}
