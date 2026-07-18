package domain

// Health representa o status de saúde do sistema
type Health struct {
	Status   string `json:"status"`
	Database string `json:"database"`
	Storage  string `json:"storage"`
	Version  string `json:"version"`
	Uptime   string `json:"uptime"`
}

// Version representa informações de versão do sistema
type Version struct {
	Version     string `json:"version"`
	Commit      string `json:"commit"`
	Build       string `json:"build"`
	Environment string `json:"environment"`
}

// Capabilities representa as capacidades do sistema
type Capabilities struct {
	Upload          bool `json:"upload"`
	SEO             bool `json:"seo"`
	Marketplace     bool `json:"marketplace"`
	IFood           bool `json:"ifood"`
	PIX             bool `json:"pix"`
	Fiscal          bool `json:"fiscal"`
	Delivery        bool `json:"delivery"`
	CardapioDigital bool `json:"cardapioDigital"`
}
