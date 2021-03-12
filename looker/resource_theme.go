package looker

import (
	"log"
	"strconv"
	"strings"

	apiclient "github.com/Foxtel-DnA/looker-go-sdk/client"
	"github.com/Foxtel-DnA/looker-go-sdk/client/theme"
	"github.com/Foxtel-DnA/looker-go-sdk/models"
	"github.com/go-openapi/strfmt"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceTheme() *schema.Resource {
	return &schema.Resource{
		Create: resourceThemeCreate,
		Read:   resourceThemeRead,
		Update: resourceThemeUpdate,
		Delete: resourceThemeDelete,
		Exists: resourceThemeExists,
		Importer: &schema.ResourceImporter{
			State: resourceThemeImport,
		},

		Schema: map[string]*schema.Schema{
			"theme_id": &schema.Schema{
				Type:     schema.TypeInt,
				Computed: true,
			},
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"begin_at": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"end_at": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"background_color": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"color_collection_id": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"font_color": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"font_family": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"font_source": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"info_button_color": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"primary_button_color": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"show_filters_bar": &schema.Schema{
				Type:     schema.TypeBool,
				Required: true,
			},
			"show_title": &schema.Schema{
				Type:     schema.TypeBool,
				Required: true,
			},
			"text_tile_text_color": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"tile_background_color": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"tile_text_color": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"title_color": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"warn_button_color": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"tile_title_alignment": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"tile_shadow": &schema.Schema{
				Type:     schema.TypeBool,
				Required: true,
			},
		},
	}
}

func getThemeByID(d *schema.ResourceData, m interface{}, id int64) (*models.Theme, error) {
	client := m.(*apiclient.Looker)

	params := theme.NewThemeParams()
	params.ThemeID = strconv.FormatInt(id, 10)

	result, err := client.Theme.Theme(params)
	if err != nil {
		log.Printf("[ERROR] Error while getting theme by id, %s", err.Error())
		return nil, err
	}

	return result.Payload, nil
}

func resourceThemeCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(*apiclient.Looker)

	// Set theme params
	params := theme.NewCreateThemeParams()
	params.Body = &models.Theme{}
	params.Body.Name = d.Get("name").(string)

	beginAt, err := strfmt.ParseDateTime(d.Get("begin_at").(string))
	if err != nil {
		return err
	}

	params.Body.BeginAt = beginAt

	endAt, err := strfmt.ParseDateTime(d.Get("end_at").(string))
	if err != nil {
		return err
	}

	params.Body.EndAt = endAt

	params.Body.Settings.BackgroundColor = d.Get("background_color").(string)
	params.Body.Settings.ColorCollectionID = d.Get("color_collection_id").(string)
	params.Body.Settings.FontColor = d.Get("font_color").(string)
	params.Body.Settings.FontFamily = d.Get("font_family").(string)
	params.Body.Settings.FontSource = d.Get("font_source").(string)
	params.Body.Settings.InfoButtonColor = d.Get("info_button_color").(string)
	params.Body.Settings.PrimaryButtonColor = d.Get("primary_button_color").(string)
	params.Body.Settings.ShowFiltersBar = d.Get("show_filters_bar").(bool)
	params.Body.Settings.ShowTitle = d.Get("show_title").(bool)
	params.Body.Settings.TextTileTextColor = d.Get("text_tile_text_color").(string)
	params.Body.Settings.TileBackgroundColor = d.Get("tile_background_color").(string)
	params.Body.Settings.TileTextColor = d.Get("tile_text_color").(string)
	params.Body.Settings.TileShadow = d.Get("tile_shadow").(bool)
	params.Body.Settings.TileTitleAlignment = d.Get("tile_title_alignment").(string)
	params.Body.Settings.TitleColor = d.Get("title_color").(string)
	params.Body.Settings.WarnButtonColor = d.Get("warn_button_color").(string)

	result, err := client.Theme.CreateTheme(params)
	if err != nil {
		return err
	}

	d.SetId(strconv.FormatInt(result.Payload.ID, 10))

	return resourceThemeRead(d, m)
}

func resourceThemeRead(d *schema.ResourceData, m interface{}) error {
	ID, err := getIDFromString(d.Id())
	if err != nil {
		return err
	}

	theme, err := getThemeByID(d, m, ID)
	if err != nil {
		if strings.Contains(err.Error(), "Not found") {
			d.SetId("")
			return nil
		}
		return err
	}

	d.Set("theme_id", theme.ID)
	d.Set("name", theme.Name)
	d.Set("begin_at", theme.BeginAt)
	d.Set("end_at", theme.EndAt)

	d.Set("background_color", theme.Settings.BackgroundColor)
	d.Set("color_collection_id", theme.Settings.ColorCollectionID)
	d.Set("font_color", theme.Settings.FontColor)
	d.Set("font_family", theme.Settings.FontFamily)
	d.Set("font_source", theme.Settings.FontSource)
	d.Set("info_button_color", theme.Settings.InfoButtonColor)
	d.Set("primary_button_color", theme.Settings.PrimaryButtonColor)
	d.Set("show_filters_bar", theme.Settings.ShowFiltersBar)
	d.Set("show_title", theme.Settings.ShowTitle)
	d.Set("text_tile_text_color", theme.Settings.TextTileTextColor)
	d.Set("tile_background_color", theme.Settings.TileBackgroundColor)
	d.Set("tile_text_color", theme.Settings.TileTextColor)
	d.Set("tile_shadow", theme.Settings.TileShadow)
	d.Set("tile_title_alignment", theme.Settings.TileTitleAlignment)
	d.Set("title_color", theme.Settings.TitleColor)
	d.Set("warn_button_color", theme.Settings.WarnButtonColor)

	return nil
}

func resourceThemeUpdate(d *schema.ResourceData, m interface{}) error {
	client := m.(*apiclient.Looker)

	ID, err := getIDFromString(d.Id())
	if err != nil {
		return err
	}

	params := theme.NewUpdateThemeParams()
	params.ThemeID = strconv.FormatInt(ID, 10)
	params.Body = &models.Theme{}
	params.Body.Name = d.Get("name").(string)

	beginAt, err := strfmt.ParseDateTime(d.Get("begin_at").(string))
	if err != nil {
		return err
	}

	params.Body.BeginAt = beginAt

	endAt, err := strfmt.ParseDateTime(d.Get("end_at").(string))
	if err != nil {
		return err
	}

	params.Body.EndAt = endAt

	params.Body.Settings.BackgroundColor = d.Get("background_color").(string)
	params.Body.Settings.ColorCollectionID = d.Get("color_collection_id").(string)
	params.Body.Settings.FontColor = d.Get("font_color").(string)
	params.Body.Settings.FontFamily = d.Get("font_family").(string)
	params.Body.Settings.FontSource = d.Get("font_source").(string)
	params.Body.Settings.InfoButtonColor = d.Get("info_button_color").(string)
	params.Body.Settings.PrimaryButtonColor = d.Get("primary_button_color").(string)
	params.Body.Settings.ShowFiltersBar = d.Get("show_filters_bar").(bool)
	params.Body.Settings.ShowTitle = d.Get("show_title").(bool)
	params.Body.Settings.TextTileTextColor = d.Get("text_tile_text_color").(string)
	params.Body.Settings.TileBackgroundColor = d.Get("tile_background_color").(string)
	params.Body.Settings.TileTextColor = d.Get("tile_text_color").(string)
	params.Body.Settings.TileShadow = d.Get("tile_shadow").(bool)
	params.Body.Settings.TileTitleAlignment = d.Get("tile_title_alignment").(string)
	params.Body.Settings.TitleColor = d.Get("title_color").(string)
	params.Body.Settings.WarnButtonColor = d.Get("warn_button_color").(string)

	_, err = client.Theme.UpdateTheme(params)
	if err != nil {
		return err
	}

	return resourceThemeRead(d, m)
}

func resourceThemeDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(*apiclient.Looker)

	ID, err := getIDFromString(d.Id())
	if err != nil {
		return err
	}

	params := theme.NewDeleteThemeParams()
	params.ThemeID = strconv.FormatInt(ID, 10)

	_, err = client.Theme.DeleteTheme(params)
	if err != nil {
		return err
	}

	return nil
}

func resourceThemeExists(d *schema.ResourceData, m interface{}) (b bool, e error) {
	// Exists - This is called to verify a resource still exists. It is called prior to Read,
	// and lowers the burden of Read to be able to assume the resource exists.
	ID, err := getIDFromString(d.Id())
	if err != nil {
		return false, err
	}

	_, err = getThemeByID(d, m, ID)
	if err != nil {
		if strings.Contains(err.Error(), "Not found") {
			return false, nil
		}

		return false, err
	}

	return true, nil
}

func resourceThemeImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	if err := resourceThemeRead(d, m); err != nil {
		return nil, err
	}
	return []*schema.ResourceData{d}, nil
}
