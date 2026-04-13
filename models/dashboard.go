package models

// Widget position in the dashboard grid
type WidgetPosition struct {
	Row   int32 `json:"row"`
	Col   int32 `json:"col"`
	SizeX int32 `json:"size_x"`
	SizeY int32 `json:"size_y"`
}

type VisualizationWidgetOutput struct {
	WidgetID        int64          `json:"widget_id"`
	VisualizationID int64          `json:"visualization_id"`
	Position        WidgetPosition `json:"position"`
}

type TextWidgetOutput struct {
	WidgetID int64          `json:"widget_id"`
	Text     string         `json:"text"`
	Position WidgetPosition `json:"position"`
}

type VisualizationWidgetInput struct {
	VisualizationID int64           `json:"visualization_id"`
	Position        *WidgetPosition `json:"position,omitempty"`
}

type TextWidgetInput struct {
	Text     string          `json:"text"`
	Position *WidgetPosition `json:"position,omitempty"`
}

type CreateDashboardRequest struct {
	Name             string            `json:"name"`
	IsPrivate        *bool             `json:"is_private,omitempty"`
	VisualizationIDs []int64           `json:"visualization_ids,omitempty"`
	TextWidgets      []TextWidgetInput `json:"text_widgets,omitempty"`
	ColumnsPerRow    *int32            `json:"columns_per_row,omitempty"`
}

type ParamWidgetOutput struct {
	WidgetID              string         `json:"widget_id"`
	Key                   string         `json:"key"`
	QueryID               int64          `json:"query_id"`
	VisualizationWidgetID int64          `json:"visualization_widget_id"`
	Position              WidgetPosition `json:"position"`
}

type ParamWidgetInput struct {
	Key                   string          `json:"key"`
	QueryID               int64           `json:"query_id"`
	VisualizationWidgetID int64           `json:"visualization_widget_id"`
	Position              *WidgetPosition `json:"position,omitempty"`
}

type DashboardResponse struct {
	DashboardID          int64                       `json:"dashboard_id"`
	Name                 string                      `json:"name"`
	Slug                 string                      `json:"slug"`
	IsPrivate            bool                        `json:"is_private"`
	Tags                 []string                    `json:"tags"`
	DashboardURL         string                      `json:"dashboard_url"`
	VisualizationWidgets []VisualizationWidgetOutput `json:"visualization_widgets"`
	TextWidgets          []TextWidgetOutput          `json:"text_widgets"`
	ParamWidgets         []ParamWidgetOutput         `json:"param_widgets"`
}

type UpdateDashboardRequest struct {
	Name                 *string                     `json:"name,omitempty"`
	Slug                 *string                     `json:"slug,omitempty"`
	IsPrivate            *bool                       `json:"is_private,omitempty"`
	Tags                 *[]string                   `json:"tags,omitempty"`
	VisualizationWidgets *[]VisualizationWidgetInput `json:"visualization_widgets,omitempty"`
	TextWidgets          *[]TextWidgetInput          `json:"text_widgets,omitempty"`
	ParamWidgets         *[]ParamWidgetInput         `json:"param_widgets,omitempty"`
	ColumnsPerRow        *int32                      `json:"columns_per_row,omitempty"`
}

type ArchiveDashboardResponse struct {
	OK          bool  `json:"ok"`
	DashboardID int64 `json:"dashboard_id"`
}
