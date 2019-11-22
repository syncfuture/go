package u

type DatatableModel struct {
	Draw    int    `url:"draw" form:"draw" json:"draw"`
	Start   int    `url:"start" form:"start" json:"start"`
	Length  int    `url:"length" form:"length" json:"length"`
	OrderBy int    `url:"order[0][column]" form:"order[0][column]" json:"order[0][column]"`
	Sort    string `url:"order[0][dir]" form:"order[0][dir]" json:"order[0][dir]"`
	Keyword string `url:"search[value]" form:"search[value]" json:"search[value]"`
}

func (x *DatatableModel) GetPageIndex() int32 {
	if x.Length <= 0 {
		return 1
	} else {
		return int32((x.Start / x.Length) + 1)
	}
}
