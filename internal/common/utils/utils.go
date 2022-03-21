package utils

func ConvertBoolToValue(b bool) float64 {
	if b {
		return 1
	} else {
		return 0
	}
}
