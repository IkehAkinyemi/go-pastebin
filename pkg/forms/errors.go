package forms

// errors type, which hold the validation error messages
// for forms. The name of the form field will be used as the key in
// this map.
type errors map[string][]string

// Add, adds error message for a given field
func (e errors) Add(field, message string) {
	e[field] = append(e[field], message)
}

func (e errors) Get(field string) string {
	errs := e[field]
	if len(errs) == 0 {
		return ""
	}
	return errs[0]
}
