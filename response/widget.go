package response

type Widget struct {
	Widgetid      int    `json:"widgetid"`
	Widgetname    string `json:"widgetname"`
	Widgettitle   string `json:"widgettitle"`
	Widgetversion string `json:"widgetversion"`
	Widgetby      string `json:"widgetby"`
	Description   string `json:"description"`
	Widgetform    string `json:"widgetform"`
	Moduleid      int    `json:"moduleid"`
	Modulename    string `json:"modulename"`
	Dashgroup     int    `json:"dashgroup"`
	Position      int    `json:"position"`
}
