package fputil

import "reflect"

func getFieldTagValue(field reflect.StructField) string {
	return field.Tag.Get("filterparams")
}

func getName(field reflect.StructField) string {
	tag := getFieldTagValue(field)

	if tag == "-" {
		return ""
	}

	name := tag

	// Assume the field name if not set.
	if name == "" {
		name = field.Name
	}

	return name
}
