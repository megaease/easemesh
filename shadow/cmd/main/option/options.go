package option

import (
	"fmt"
)

type Options struct {
	MeshServer string
}

// Validate Options
func (o *Options) Validate() []error {
	var errors []error
	if (o.MeshServer == "") {
		errors = append(errors, fmt.Errorf("MeshServer is required."))
	}
	return errors
}

