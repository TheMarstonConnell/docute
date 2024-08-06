package gen

type Colors struct {
	Text       string `yaml:"text"`
	Primary    string `yaml:"primary"`
	Secondary  string `yaml:"secondary"`
	Background string `yaml:"background"`
	TitleBar   string `yaml:"titleBar"`
}

func DefaultColors() Colors {
	return Colors{
		Text:       "#040805",
		Primary:    "#809DBC",
		Secondary:  "#9BC1CA",
		Background: "#F4F9F5",
		TitleBar:   "#809DBC",
	}
}
