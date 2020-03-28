package models

type Project struct {
	ID		int64		`json:"id"`
	Title	string		`json:"title"`
	Tags	[]string	`json:"tags"`
}
