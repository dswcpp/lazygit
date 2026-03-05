package constants

type Docs struct {
	CustomPagers      string
	CustomCommands    string
	CustomKeybindings string
	Keybindings       string
	Undoing           string
	Config            string
	Tutorial          string
	CustomPatchDemo   string
}

var Links = struct {
	Docs    Docs
	Issues  string
	RepoUrl string
	Releases string
}{
	RepoUrl:  "https://github.com/dswcpp/lazygit",
	Issues:   "https://github.com/dswcpp/lazygit/issues",
	Releases: "https://github.com/dswcpp/lazygit/releases",
	Docs: Docs{
		CustomPagers:      "https://github.com/dswcpp/lazygit/blob/master/docs/Custom_Pagers.md",
		CustomKeybindings: "https://github.com/dswcpp/lazygit/blob/master/docs/keybindings/Custom_Keybindings.md",
		CustomCommands:    "https://github.com/dswcpp/lazygit/wiki/Custom-Commands-Compendium",
		Keybindings:       "https://github.com/dswcpp/lazygit/blob/%s/docs/keybindings",
		Undoing:           "https://github.com/dswcpp/lazygit/blob/master/docs/Undoing.md",
		Config:            "https://github.com/dswcpp/lazygit/blob/%s/docs/Config.md",
		Tutorial:          "https://youtu.be/VDXvbHZYeKY",
		CustomPatchDemo:   "https://github.com/dswcpp/lazygit#rebase-magic-custom-patches",
	},
}
