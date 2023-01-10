package provider

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func Provider() *schema.Provider {
	return &schema.Provider{
		DataSourcesMap: map[string]*schema.Resource{
			"json_patch": dataSourceJsonPatch(),
			"yaml_patch": dataSourceYamlPatch(),
		},
	}
}
