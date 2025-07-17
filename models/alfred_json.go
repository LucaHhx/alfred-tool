package models

type ModName string

const (
	Mod_Cmd   ModName = "cmd"
	Mod_Alt   ModName = "alt"
	Mod_Ctrl  ModName = "ctrl"
	Mod_Shift ModName = "shift"
	Mod_Fn    ModName = "fn"
)

type AlfredData struct {
	Items []AlfredItem `json:"items"`
}

type AlfredItem struct {
	Uid       string                `json:"uid"`
	Title     string                `json:"title"`
	Subtitle  string                `json:"subtitle"`
	Arg       []string              `json:"arg"`
	Mods      map[ModName]AlfredMod `json:"mods"`
	Variables map[string]string     `json:"variables"`
}

type AlfredMod struct {
	Valid    bool     `json:"valid"`
	Arg      []string `json:"arg"`
	Subtitle string   `json:"subtitle"`
}

func NewAlfredMod(sub string, args ...string) AlfredMod {
	return AlfredMod{
		true,
		args,
		sub,
	}
}
