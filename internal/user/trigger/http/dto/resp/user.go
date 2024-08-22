package resp

type UserLoginResp struct {
	Token string `json:"token"`
}

type UserResp struct {
	Id       int64  `json:"id"`
	Username string `json:"username"`
	RealName string `json:"real_name"`
	Phone    string `json:"phone"`
	Mail     string `json:"mail"`
}

type UserActualResp struct {
	Id       int64  `json:"id"`
	Username string `json:"username"`
	RealName string `json:"real_name"`
	Phone    string `json:"phone"`
	Mail     string `json:"mail"`
}

type ShortLinkGroupResp []ShortLinkGroupDTO

type ShortLinkGroupDTO struct {
	Gid            string `json:"gid"`
	Name           string `json:"name"`
	SortOrder      int    `json:"sort_order"`
	ShortLinkCount int    `json:"short_link_count"`
}
