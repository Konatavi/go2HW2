package models

//Article models...
type Automobiles struct {
	ID       int    `json:"id"`
	Mark     string `json:"mark"`
	Maxspeed int    `json:"max_speed"`
	Distance int    `json:"distance"`
	Handler  string `json:"handler"`
	Stock    string `json:"stock"`
}
