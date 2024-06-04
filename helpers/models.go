package helpers

type Response struct {
    Count    int      `json:"count"`
    Next     *string  `json:"next"`
    Previous *string  `json:"previous"`
    Results  []Result `json:"results"`
}

type Result struct {
    ID          int    `json:"id"`
    File        string `json:"file"`
    Mission     int    `json:"mission"`
    DisplayName string `json:"display_name"`
    Name        string `json:"name"`
}

type MissionsResponse struct {
	Results []Mission `json:"results"`
}

type Mission struct {
	Id   int    `json:"id"`
	Date string `json:"date"`
	Pointcloud []Pointcloud `json:"pointclouds"`
}

type Pointcloud struct {
	Id   int    `json:"id"`
}