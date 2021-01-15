package covid

type CountyResponse struct {
	County  string `json:"county"`
	State   string `json:"state"`
	Cases   int    `json:"cases"`
	Deaths  int    `json:"deaths"`
	Updated int    `json:"updated"`
}
