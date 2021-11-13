package server

func strOrNil(v string) *string {
	if v == "" {
		return nil
	}
	return &v
}

func strOrNull(ptr *string) string {
	switch ptr {
	case nil:
		return "null"
	default:
		return *ptr
	}
}

// newOrCurr returns the new value if it is not empty, otherwise the current value.
func PtrNewOrCurr(new string, curr *string) string {
	if new == "" {
		return *curr
	}
	return new
}

func NewOrCurr(new string, curr string) string {
	if new == "" {
		return curr
	}
	return new
}

func ptrToStr(v *string) string {
	if v == nil {
		return ""
	}
	return *v
}
