package FlareModel

// Generic Bookmark Data Model
type Bookmark struct {
	Name     string `yaml:"name" json:"name"`
	URL      string `yaml:"link" json:"link"`
	Icon     string `yaml:"icon,omitempty" json:"icon,omitempty"`
	Desc     string `yaml:"desc,omitempty" json:"desc,omitempty"`
	Private  bool   `yaml:"private,omitempty" json:"private,omitempty"`
	Category string `yaml:"category,omitempty" json:"category,omitempty"`
}

// Generic Category Data Model
type Category struct {
	ID   string `yaml:"id" json:"id"`
	Name string `yaml:"title" json:"title"`
}

// Generic Bookmarks Data Model
type Bookmarks struct {
	Categories []Category `yaml:"categories,omitempty" json:"categories,omitempty"`
	Items      []Bookmark `yaml:"links" json:"links"`
}
