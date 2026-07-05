package convert

import (
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// SetReferenceField keeps the user's reference (identity, slug, or name) when it still matches
// the API value; otherwise it writes back the API identity so out-of-band changes are detected.
func SetReferenceField(d *schema.ResourceData, field, identity, slug, name string) {
	current := d.Get(field).(string)
	switch {
	case current == "":
		_ = d.Set(field, identity)
	case current == identity:
		_ = d.Set(field, current)
	case slug != "" && current == slug:
		_ = d.Set(field, current)
	case name != "" && strings.EqualFold(current, name):
		_ = d.Set(field, current)
	default:
		_ = d.Set(field, identity)
	}
}
