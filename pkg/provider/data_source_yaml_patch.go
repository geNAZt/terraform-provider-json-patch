package provider

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	jsonpatch "github.com/evanphx/json-patch/v5"
	convert "github.com/icza/dyno"
	yaml "gopkg.in/yaml.v3"
)

func dataSourceYamlPatch() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceYamlPatchRead,
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

func dataSourceYamlPatchRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	document := d.Get("document").(string)

	// We need to convert document to JSON first
	var body interface{}
	err := yaml.Unmarshal([]byte(document), &body)
	if err != nil {
		return diag.FromErr(err)
	}

	body = convert.ConvertMapI2MapS(body)
	jsonDocument, err := json.Marshal(body)
	if err != nil {
		return diag.FromErr(err)
	}

	document = string(jsonDocument)
	patches := d.Get("patches").([]any)

	tflog.Info(ctx, fmt.Sprintf("Document: %s", document))

	for _, patch := range patches {
		tflog.Info(ctx, fmt.Sprintf("Patch: %v", patch))

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

	// Convert back to yaml
	err = json.Unmarshal([]byte(document), &body)
	if err != nil {
		return diag.FromErr(err)
	}

	body = convert.ConvertMapI2MapS(body)
	yamlDocument, err := yaml.Marshal(body)
	if err != nil {
		return diag.FromErr(err)
	}

	err = d.Set("patched", string(yamlDocument))
	if err != nil {
		return diag.FromErr(err)
	}

	tflog.Info(ctx, fmt.Sprintf("Patched: %s", yamlDocument))

	d.SetId(MakeId(yamlDocument))
	return nil
}
