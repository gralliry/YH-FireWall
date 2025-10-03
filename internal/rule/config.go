package rule

type Config struct {
	Id          uint32 `json:"id"`
	Description string `json:"description"`
	SrcNet      string `json:"src_net"`
	SrcPort     string `json:"src_port"`
	TarNet      string `json:"tar_net"`
	TarPort     string `json:"tar_port"`
	InDev       string `json:"in_dev"`
	OutDev      string `json:"out_dev"`
	Protocol    string `json:"protocol"`
	Accept      bool   `json:"accept"`
	Priority    uint32 `json:"priority"`
	Enable      bool   `json:"enable"`
	Group       uint16 `json:"group"`
}
