package forms

import (
	"fmt"
	"net/url"
	"strings"
	"unicode/utf8"
)

// Form struct, which anonymously embeds a url.Values object
// (to hold the form data) and an Errors field to hold any validation errors
// for the form data.
type Form struct {
	url.Values
	Errors errors
}

// New function initializes a custom Form struct
func New(data url.Values) *Form {
	return &Form{
		data,
		errors(map[string][]string{}),
	}
}

// Required checks that specific fields in the form 
// data are present and not blank, fails then add the 
// appropriate message to the form errors.
func (f *Form) Required(fields ...string)  {
	for _, field := range fields {
		value := f.Get(field)
		if strings.TrimSpace(value) == "" {
			f.Errors.Add(field, "This field cannot be blank")
		}
	}
}

// MaxLength checks that a specific field in the form 
// contains a maximum number of characters, fails then add the 
// appropriate message to the form errors.
func (f *Form) MaxLengthChars(field string, d int) {
	value := f.Get(field)
	if strings.TrimSpace(value) == "" {
		return
	}
	if utf8.RuneCountInString(value) > d {
		f.Errors.Add(field, fmt.Sprintf("This field is too long(maximumis %d characters", d))
	}
}

// PermittedValues checks that a specific field in the form 
// matches one of a set of specific permitted values, fails
// then add the appropriate message to the form errors.
func (f *Form) PermittedValues(field string, opts ...string) {
	value := f.Get(field)

	for _, opt := range opts {
		if value == opt {
			return
		}
	}
	f.Errors.Add(field, "This field is invalid")
}

// Valid returns true if there are no errors.
func (f *Form) Valid() bool {
	return len(f.Errors) == 0
}
