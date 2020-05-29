package kmetal

func s(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

func b(s *bool) bool {
	if s == nil {
		return false
	}
	return *s
}
