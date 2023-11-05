package brand

type Brand struct {
	Id    int32  `json:"id"`
	Title string `json:"title"`
}

type SelectedBrand struct {
	Brand
	Selected bool
}
