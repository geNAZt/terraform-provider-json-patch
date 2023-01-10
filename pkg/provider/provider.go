package provider

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func Provider() *schema.Provider {
	return &schema.Provider{
		DataSourcesMap: map[string]*schema.Resource{
			"json-patch_json_patch": dataSourceJsonPatch(),
			"json-patch_yaml_patch": dataSourceYamlPatch(),
		},
	}
}
