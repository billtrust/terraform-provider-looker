package looker

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/billtrust/looker-go-sdk/client/content"

	"github.com/billtrust/looker-go-sdk/client/space"

	apiclient "github.com/billtrust/looker-go-sdk/client"
	"github.com/billtrust/looker-go-sdk/models"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceMainSpace() *schema.Resource {
	return &schema.Resource{
		Create: resourceMainSpaceCreate,
		Read:   resourceMainSpaceRead,
		Update: resourceMainSpaceUpdate,
		Delete: resourceMainSpaceDelete,
		Exists: resourceMainSpaceExists,
		Importer: &schema.ResourceImporter{
			State: resourceMainSpaceImport,
		},

		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"parent_space_name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"parent_id": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"parent_content_metadata_id": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"content_metadata_id": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"content_metadata_inherits": &schema.Schema{
				Type:     schema.TypeBool,
				Required: true,
			},
		},
	}
}

func getRootSpace(d *schema.ResourceData, m interface{}, name string) (*models.Space, error) {
	client := m.(*apiclient.Looker)

	params := space.NewSearchSpacesParams()
	params.Name = &name

	result, err := client.Space.SearchSpaces(params)
	if err != nil {
		log.Printf("[ERROR] Error while searching spaces with name '%s', %s", name, err.Error())
		return nil, err
	}

	for _, item := range result.Payload {
		if *item.Name == name && &item.ParentID == nil {
			return item, nil
		}
	}

	if name == "Embed Groups" {
		return nil, fmt.Errorf("[ERROR] 'Embed Groups' does not exist. Goto https://[domain].looker.com/admin/embed and set 'Embed Authentication' to Enabled")
	}

	return nil, fmt.Errorf("No root space with name '%s'", name)
}

func getSpaceByID(d *schema.ResourceData, m interface{}, id int64) (*models.Space, error) {
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

func resourceMainSpaceCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(*apiclient.Looker)

	rootSpace, err := getRootSpace(d, m, d.Get("parent_space_name").(string))
	if err != nil {
		return err
	}

	params := space.NewCreateSpaceParams()
	params.Body = &models.CreateSpace{}
	params.Body.Name = d.Get("name").(*string)
	params.Body.ParentID = &rootSpace.ID

	result, err := client.Space.CreateSpace(params)
	if err != nil {
		return err
	}

	d.SetId(result.Payload.ID)

	// TODO: should SetId happen before or after logic to update content_metadata?
	contentMetadataParams := content.NewUpdateContentMetadataParams()
	contentMetadataParams.ContentMetadataID = result.Payload.ContentMetadataID
	contentMetadataParams.Body = &models.ContentMeta{}
	inherits := d.Get("content_metadata_inherits").(bool)
	contentMetadataParams.Body.Inherits = inherits

	_, err = client.Content.UpdateContentMetadata(contentMetadataParams)
	if err != nil {
		return err
	}

	return resourceMainSpaceRead(d, m)
}

func resourceMainSpaceRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*apiclient.Looker)

	ID, err := getIDFromString(d.Id())
	if err != nil {
		return err
	}

	space, err := getSpaceByID(d, m, ID)
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
	var parentIDint, _ = strconv.ParseInt(space.ParentID, 10, 64)
	parentSpace, err := getSpaceByID(d, m, parentIDint)
	if err != nil {
		return err
	}

	d.Set("parent_space_name", parentSpace.Name)
	d.Set("parent_content_metadata_id", getStringFromID(parentSpace.ContentMetadataID))

	contentMetadataParams := content.NewContentMetadataParams()
	contentMetadataParams.ContentMetadataID = space.ContentMetadataID

	contentMetadataResult, err := client.Content.ContentMetadata(contentMetadataParams)
	if err != nil {
		return err
	}

	d.Set("content_metadata_inherits", contentMetadataResult.Payload.Inherits)

	return nil
}

func resourceMainSpaceUpdate(d *schema.ResourceData, m interface{}) error {
	client := m.(*apiclient.Looker)

	ID, err := getIDFromString(d.Id())
	if err != nil {
		return err
	}

	rootSpace, err := getRootSpace(d, m, d.Get("parent_space_name").(string))
	if err != nil {
		return err
	}

	params := space.NewUpdateSpaceParams()
	params.SpaceID = strconv.FormatInt(ID, 10)
	params.Body = &models.UpdateSpace{}
	params.Body.Name = d.Get("name").(string)
	params.Body.ParentID = rootSpace.ID

	_, err = client.Space.UpdateSpace(params)
	if err != nil {
		return err
	}

	if d.HasChange("content_metadata_inherits") {
		contentMetadataID, err := getIDFromString(d.Get("content_metadata_id").(string))
		if err != nil {
			return err
		}

		inherits := false

		contentMetadataInherts := d.Get("content_metadata_inherits").(bool)
		contentMetadataParams := content.NewUpdateContentMetadataParams()
		contentMetadataParams.ContentMetadataID = contentMetadataID
		contentMetadataParams.Body = &models.ContentMeta{}
		contentMetadataParams.Body.Inherits = inherits
		contentMetadataParams.Body.Name = d.Get("name").(string)
		contentMetadataParams.Body.InheritingID = 0

		log.Printf("[DEBUG] metadata id: %d, inherits %t", contentMetadataID, contentMetadataInherts)

		_, err = client.Content.UpdateContentMetadata(contentMetadataParams)
		if err != nil {
			return err
		}
	}

	return resourceMainSpaceRead(d, m)
}

func resourceMainSpaceDelete(d *schema.ResourceData, m interface{}) error {
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

func resourceMainSpaceExists(d *schema.ResourceData, m interface{}) (b bool, e error) {
	// Exists - This is called to verify a resource still exists. It is called prior to Read,
	// and lowers the burden of Read to be able to assume the resource exists.
	ID, err := getIDFromString(d.Id())
	if err != nil {
		return false, err
	}

	_, err = getSpaceByID(d, m, ID)
	if err != nil {
		if strings.Contains(err.Error(), "Not found") {
			return false, nil
		}

		return false, err
	}

	return true, nil
}

func resourceMainSpaceImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	if err := resourceMainSpaceRead(d, m); err != nil {
		return nil, err
	}
	return []*schema.ResourceData{d}, nil
}
