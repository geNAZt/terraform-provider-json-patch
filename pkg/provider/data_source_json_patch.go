package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	jsonpatch "github.com/evanphx/json-patch/v5"
)

func dataSourceJsonPatch() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceJsonPatchRead,
		Schema: map[string]*schema.Schema{
			"document": {
				Type:     schema.TypeString,
				Required: true,
			},
			"patches": {
				Type:     schema.TypeList,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Required: true,
			},
			"patched": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceJsonPatchRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	document := d.Get("document").(string)
	patches := d.Get("patches").([]string)

	for _, patch := range patches {
		patchedDocument, err := jsonpatch.MergePatch([]byte(document), []byte(patch))
		if err != nil {
			return diag.FromErr(err)
		}

		document = string(patchedDocument)
	}

	d.Set("patched", document)
	return nil
}
