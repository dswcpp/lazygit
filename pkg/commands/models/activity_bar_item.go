package models

// ActivityBarItem represents an item in the activity bar
type ActivityBarItem struct {
	Icon      IconConfig       // Icon configuration
	Name      string           // Internal name
	Type      ActivityItemType // Item type
	Tooltip   string           // Tooltip text
	Action    string           // Action identifier
	CustomCmd string           // Custom command (optional)
	Shortcut  string           // Keyboard shortcut display (optional)
}

// ActivityItemType defines the type of activity bar item
type ActivityItemType int

const (
	ActivityTypeNavigation ActivityItemType = iota // Navigation item (switches context)
	ActivityTypeAction                             // Git action item (pull, push, etc.)
	ActivityTypeTool                               // Tool item (settings, help)
	ActivityTypeSeparator                          // Visual separator
	ActivityTypeCustom                             // Custom command
)

// IconConfig holds icon configuration with fallback options
type IconConfig struct {
	NerdFont string // Nerd Font icon (preferred)
	Emoji    string // Emoji fallback
	ASCII    string // ASCII fallback
	Color    string // Color (optional, hex format #RRGGBB)
}

// ID returns the unique identifier for the activity bar item
func (self *ActivityBarItem) ID() string {
	return self.Name
}

// IsNavigationItem returns true if this is a navigation item
func (self *ActivityBarItem) IsNavigationItem() bool {
	return self.Type == ActivityTypeNavigation
}

// IsActionItem returns true if this is a git action item
func (self *ActivityBarItem) IsActionItem() bool {
	return self.Type == ActivityTypeAction
}

// IsSeparator returns true if this is a separator
func (self *ActivityBarItem) IsSeparator() bool {
	return self.Type == ActivityTypeSeparator
}
