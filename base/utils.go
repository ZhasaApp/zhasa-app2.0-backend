package base

func Contains[T comparable](slice []T, element T) bool {
	for _, v := range slice {
		if v == element {
			return true
		}
	}
	return false
}

type ArrayResponse[T any] struct {
	Result []T `json:"result"`
}
