package helpers

func StringToBool(s string) *bool {

	if s == "" {
		return nil
	}

	b := s == "true"

	return &b
}

func GetFilterBtnClass(currentFilter string, buttonFilter string) string {
	if currentFilter == buttonFilter {
		return "btn btn-sm glass-primary"
	}

	return "btn btn-sm btn-outline border-white/20 text-white hover:bg-white/10 bg-transparent"
}
