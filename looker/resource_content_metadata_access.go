package looker

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/billtrust/looker-go-sdk/client/content"

	apiclient "github.com/billtrust/looker-go-sdk/client"
	"github.com/billtrust/looker-go-sdk/models"
	"github.com/hashicorp/terraform/helper/schema"
)

// The GET resource for metadata content access does not exist, we need to search all the metadata content access for a specific contentmetadaid and than look for the group id.
// for this reason, i am setting the id to be "content_metadata_id:group_id", and have content_metadata_access_id be a computed field
// Since this is just access and nothing depends on it, I Am going to make this simpler by just implenting create, delete, read and skipping update
func resourceContentMetadataAccess() *schema.Resource {
	return &schema.Resource{
		Create: resourceContentMetadataAccessCreate,
		Read:   resourceContentMetadataAccessRead,
		Delete: resourceContentMetadataAccessDelete,
		Exists: resourceContentMetadataAccessExists,
		Importer: &schema.ResourceImporter{
			State: resourceContentMetadataAccessImport,
		},

		Schema: map[string]*schema.Schema{
			"content_metadata_access_id": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"group_id": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"content_metadata_id": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"permission_type": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func getContentMetadataAccess(m interface{}, contentMetadataID int64, groupID int64) (*models.ContentMetaGroupUser, error) {
	client := m.(*apiclient.Looker)

	params := content.NewAllContentMetadataAccessesParams()
	params.SetTimeout(time.Minute * 5)
	params.ContentMetadataID = contentMetadataID

	result, err := client.Content.AllContentMetadataAccesses(params)
	if err != nil {
		return nil, err
	}

	for _, contentMetaGroupUser := range result.Payload {
		if contentMetaGroupUser.GroupID == groupID {
			return contentMetaGroupUser, nil
		}
	}

	return nil, fmt.Errorf("Content Metadata Access Not found")
}

func resourceContentMetadataAccessCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(*apiclient.Looker)

	groupID, err := getIDFromString(d.Get("group_id").(string))
	if err != nil {
		return err
	}

	contentMetadataID, err := getIDFromString(d.Get("content_metadata_id").(string))
	if err != nil {
		return err
	}

	params := content.NewCreateContentMetadataAccessParams()
	params.SetTimeout(time.Minute * 5)
	params.Body = &models.ContentMetaGroupUser{}
	params.Body.ContentMetadataID = strconv.FormatInt(contentMetadataID, 10)
	params.Body.GroupID = groupID
	params.Body.PermissionType = d.Get("permission_type").(string)

	result, err := client.Content.CreateContentMetadataAccess(params)
	if err != nil {
		if !strings.Contains(err.Error(), "already has access on content") {
			return err
		}

		access, err := getContentMetadataAccess(m, contentMetadataID, groupID)
		if err != nil {
			return err
		}

		d.SetId(access.ContentMetadataID + ":" + getStringFromID(access.GroupID))
	} else {
		d.SetId(result.Payload.ContentMetadataID + ":" + getStringFromID(result.Payload.GroupID))
	}

	return resourceContentMetadataAccessRead(d, m)
}

func resourceContentMetadataAccessRead(d *schema.ResourceData, m interface{}) error {
	id := strings.Split(d.Id(), ":")
	if len(id) != 2 {
		return fmt.Errorf("ID Should be two strings separated by a colon (:)")
	}

	sContentMetadataID := id[0]
	sGroupID := id[1]

	groupID, err := getIDFromString(sGroupID)
	if err != nil {
		return err
	}

	contentMetadataID, err := getIDFromString(sContentMetadataID)
	if err != nil {
		return err
	}

	access, err := getContentMetadataAccess(m, contentMetadataID, groupID)
	if err != nil {
		if strings.Contains(err.Error(), "Not found") {
			d.SetId("")
			return nil
		}
		return err
	}

	d.Set("content_metadata_access_id", access.ID)
	d.Set("group_id", getStringFromID(access.GroupID))
	d.Set("content_metadata_id", access.ContentMetadataID)
	d.Set("permission_type", access.PermissionType)

	return nil
}

func resourceContentMetadataAccessDelete(d *schema.ResourceData, m interface{}) error {
	id := strings.Split(d.Id(), ":")
	if len(id) != 2 {
		return fmt.Errorf("ID Should be two strings separated by a colon (:)")
	}

	sContentMetadataID := id[0]
	sGroupID := id[1]

	groupID, err := getIDFromString(sGroupID)
	if err != nil {
		return err
	}

	contentMetadataID, err := getIDFromString(sContentMetadataID)
	if err != nil {
		return err
	}

	access, err := getContentMetadataAccess(m, contentMetadataID, groupID)
	if err != nil {
		// if attempting to delete and it is already deleted say it was succesful
		if strings.Contains(err.Error(), "Not found") {
			return nil
		}
		return err
	}

	client := m.(*apiclient.Looker)
	params := content.NewDeleteContentMetadataAccessParams()
	params.SetTimeout(time.Minute * 5)
	var accessIDStr, _ = strconv.ParseInt(access.ID, 10, 32)
	params.ContentMetadataAccessID = accessIDStr

	_, err = client.Content.DeleteContentMetadataAccess(params)
	if err != nil {
		// if the error is "Cannot remove access for [group_name] Group with edit on parent", I think the correct thing to do is ignore this error since the user already has edit on parent.
		// When/if parent access is deleted, it deletes access on child (verified)
		if strings.Contains(err.Error(), "with edit on parent") {
			log.Printf("[WARN] Deleting access from child does not work since it is inherited from parent., %s", err.Error())
			return nil
		}
		return err
	}

	return nil
}

func resourceContentMetadataAccessExists(d *schema.ResourceData, m interface{}) (b bool, e error) {
	// Exists - This is called to verify a resource still exists. It is called prior to Read,
	// and lowers the burden of Read to be able to assume the resource exists.
	id := strings.Split(d.Id(), ":")
	if len(id) != 2 {
		return false, fmt.Errorf("ID Should be two strings separated by a colon (:)")
	}
	sContentMetadataID := id[0]
	sGroupID := id[1]

	groupID, err := getIDFromString(sGroupID)
	if err != nil {
		return false, err
	}

	contentMetadataID, err := getIDFromString(sContentMetadataID)
	if err != nil {
		return false, err
	}

	_, err = getContentMetadataAccess(m, contentMetadataID, groupID)
	if err != nil {
		if strings.Contains(err.Error(), "Not found") {
			return false, nil
		}
		return false, err
	}

	return true, nil
}

func resourceContentMetadataAccessImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	if err := resourceContentMetadataAccessRead(d, m); err != nil {
		return nil, err
	}
	return []*schema.ResourceData{d}, nil
}
