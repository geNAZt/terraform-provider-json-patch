package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-log/tflog"
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
	patches := d.Get("patches").([]interface{})

	tflog.Info(ctx, fmt.Sprintf("Document: %s", document))

	for _, patch := range patches {
		tflog.Info(ctx, fmt.Sprintf("Patch: %v", patch))

		if patch == nil {
			continue
		}

		patchStr := patch.(string)
		if patchStr == "" {
			continue
		}

		patch, err := jsonpatch.DecodePatch([]byte(patchStr))
		if err != nil {
			return diag.FromErr(err)
		}

		patchedDocument, err := patch.Apply([]byte(document))
		if err != nil {
			return diag.FromErr(err)
		}

		document = string(patchedDocument)
	}

	d.Set("patched", document)
	return nil
}
