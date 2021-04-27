package looker

import (
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"

	apiclient "github.com/Foxtel-DnA/looker-go-sdk/client"
	cc "github.com/Foxtel-DnA/looker-go-sdk/client/color_collection"
	"github.com/Foxtel-DnA/looker-go-sdk/models"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceColorCollection() *schema.Resource {
	return &schema.Resource{
		Create: resourceColorCollectionCreate,
		Read:   resourceColorCollectionRead,
		Update: resourceColorCollectionUpdate,
		Delete: resourceColorCollectionDelete,
		Exists: resourceColorCollectionExists,
		Importer: &schema.ResourceImporter{
			State: resourceColorCollectionImport,
		},

		Schema: map[string]*schema.Schema{
			"color_collection_id": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"label": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"categorical_palettes": &schema.Schema{
				Type:     schema.TypeList,
				Required: true,
				Elem:     &schema.Schema{
					Type: schema.TypeList,
					Elem: &schema.Schema{
						Type: schema.TypeString,
					},
				},
			},
			"sequential_palettes": &schema.Schema{
				Type:     schema.TypeList,
				Required: true,
				Elem:     &schema.Schema{
					Type: schema.TypeList,
					Elem: &schema.Schema{
						Type: schema.TypeMap,
						Elem: &schema.Schema{
							Type: schema.TypeString,
						},
					},
				},
			},
			"diverging_palettes": &schema.Schema{
				Type:     schema.TypeList,
				Required: true,
				Elem:     &schema.Schema{
					Type: schema.TypeList,
					Elem: &schema.Schema{
						Type: schema.TypeMap,
						Elem: &schema.Schema{
							Type: schema.TypeString,
						},
					},
				},
			},
		},
	}
}

func getColorCollectionByID(d *schema.ResourceData, m interface{}, id int64) (*models.ColorCollection, error) {
	client := m.(*apiclient.Looker)

	params := cc.NewColorCollectionParams()
	params.CollectionID = strconv.FormatInt(id, 10)

	result, err := client.ColorCollection.ColorCollection(params)
	if err != nil {
		log.Printf("[ERROR] Error while getting color_collection by id, %s", err.Error())
		return nil, err
	}

	return result.Payload, nil
}

func resourceColorCollectionCreate(d *schema.ResourceData, m interface{}) error {
	var err error
	client := m.(*apiclient.Looker)

	// Set color_collection params
	params := cc.NewCreateColorCollectionParams()
	params.Body, err = parseColorCollection(d)
	if err != nil {
		return err
	}
	result, err := client.ColorCollection.CreateColorCollection(params)
	if err != nil {
		return err
	}

	d.SetId(result.Payload.ID)

	return resourceColorCollectionRead(d, m)
}

func parseColorCollection(d *schema.ResourceData) (*models.ColorCollection, error) {
	var err error
	collection := new(models.ColorCollection)
	collection.Label = d.Get("label").(string)

	cps := d.Get("categorical_palettes").([][]string)
	collection.CategoricalPalettes = make([]*models.DiscretePalette, len(cps))
	for i, cp := range cps {
		r := new(models.DiscretePalette)
		r.Type = "Categorical"
		r.Label = fmt.Sprintf("categorical_%d", i)
		nc := len(cp)
		r.Colors = make([]string, nc)
		for j, color := range cp {
			r.Colors[j] = color
		}
		collection.CategoricalPalettes[i] = r
	}

	sps := d.Get("sequential_palettes").([][]map[string]string)
	collection.SequentialPalettes = make([]*models.ContinuousPalette, len(sps))
	for i, sp := range sps {
		r := new(models.ContinuousPalette)
		r.Type = "Sequential"
		r.Label = fmt.Sprintf("sequential_%d", i)
		r.Stops, err = parseColorStops(sp)
		if err != nil {
			return nil, err
		}
		collection.SequentialPalettes[i] = r
	}


	dps := d.Get("diverging_palettes").([][]map[string]string)
	collection.DivergingPalettes = make([]*models.ContinuousPalette, len(dps))
	for i, dp := range dps {
		r := new(models.ContinuousPalette)
		r.Type = "Diverging"
		r.Label = fmt.Sprintf("diverging_%d", i)
		r.Stops, err = parseColorStops(dp)
		if err != nil {
			return nil, err
		}
		collection.DivergingPalettes[i] = r
	}
	return collection, nil
}


func parseColorStops(ms []map[string]string) ([]*models.ColorStop, error) {
	var err error
	stops := make([]*models.ColorStop, len(ms))
	for i, s := range ms {
		stops[i], err = parseColorStop(s)
		if err != nil {
			return nil, err
		}
	}
	return stops, nil
}


func parseColorStop(m map[string]string) (*models.ColorStop, error) {
	var err error
	stop := new(models.ColorStop)
	stop.Offset, err = colorStopOffset(m)
	if err != nil {
		return nil, err
	}
	stop.Color, err = colorStopColor(m)
	if err != nil {
		return nil, err
	}
	return stop, nil
}

func colorStopColor(m map[string]string) (string, error) {
	v, ok := m["color"]
	// FIXME: maybe validate color string here
	if !ok {
		return "", errors.New("invalid color stop: missing required color value")
	}
	return v, nil
}

func colorStopOffset(m map[string]string) (int64, error) {
	v, ok := m["offset"]
	if !ok {
		return 0, errors.New("invalid color stop: missing required offset value")
	}
	offset, err := strconv.ParseInt(v, 10, 64)
	if err != nil {
		return 0, err
	}
	if offset < 0 || offset > 100 {
		return 0, fmt.Errorf("invalid color stop: offset must be between 0 and 100 inclusive (got %d)", offset)
	}
	return offset, nil
}

func resourceColorCollectionRead(d *schema.ResourceData, m interface{}) error {
	ID, err := getIDFromString(d.Id())
	if err != nil {
		return err
	}

	cc, err := getColorCollectionByID(d, m, ID)
	if err != nil {
		if strings.Contains(err.Error(), "Not found") {
			d.SetId("")
			return nil
		}
		return err
	}

	d.Set("color_collection_id", cc.ID)
	d.Set("label", cc.Label)
	d.Set("categorical_palettes", resourceFromDiscretePalettes(cc.CategoricalPalettes))
	d.Set("sequential_palettes", resourceFromContinuousPalettes(cc.SequentialPalettes))
	d.Set("diverging_palettes", resourceFromContinuousPalettes(cc.DivergingPalettes))

	return nil
}

func resourceFromDiscretePalettes(ps []*models.DiscretePalette) [][]string {
	v := make([][]string, len(ps))
	for i, dp := range ps {
		v[i] = make([]string, len(dp.Colors))
		for j, color := range dp.Colors {
			v[i][j] = color
		}
	}
	return v
}


func resourceFromContinuousPalettes(ps []*models.ContinuousPalette) [][]map[string]string {
	v := make([][]map[string]string, len(ps))
	for i, cp := range ps {
		v[i] = make([]map[string]string, len(cp.Stops))
		for j, stop := range cp.Stops {
			v[i][j] = resourceFromColorStop(stop)
		}
	}
	return v
}

func resourceFromColorStop(s *models.ColorStop) map[string]string {
	m := make(map[string]string)
	m["color"] = s.Color
	m["offset"] = strconv.FormatInt(s.Offset, 10)
	return m
}

func resourceColorCollectionUpdate(d *schema.ResourceData, m interface{}) error {
	var err error
	client := m.(*apiclient.Looker)

	ID, err := getIDFromString(d.Id())
	if err != nil {
		return err
	}

	params := cc.NewUpdateColorCollectionParams()
	params.CollectionID = strconv.FormatInt(ID, 10)
	params.Body, err = parseColorCollection(d)
	if err != nil {
		return err
	}
	_, err = client.ColorCollection.UpdateColorCollection(params)
	if err != nil {
		return err
	}

	return resourceColorCollectionRead(d, m)
}

func resourceColorCollectionDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(*apiclient.Looker)

	ID, err := getIDFromString(d.Id())
	if err != nil {
		return err
	}

	params := cc.NewDeleteColorCollectionParams()
	params.CollectionID = strconv.FormatInt(ID, 10)

	_, _, err = client.ColorCollection.DeleteColorCollection(params)
	if err != nil {
		return err
	}

	return nil
}

func resourceColorCollectionExists(d *schema.ResourceData, m interface{}) (b bool, e error) {
	ID, err := getIDFromString(d.Id())
	if err != nil {
		return false, err
	}

	_, err = getColorCollectionByID(d, m, ID)
	if err != nil {
		if strings.Contains(err.Error(), "Not found") {
			return false, nil
		}

		return false, err
	}

	return true, nil
}

func resourceColorCollectionImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	if err := resourceColorCollectionRead(d, m); err != nil {
		return nil, err
	}
	return []*schema.ResourceData{d}, nil
}
