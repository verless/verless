package model

import "time"

type Page struct {
	ID          string
	Title       string
	Author      string
	Date        time.Time
	Tags        []string
	Img         string
	Credit      string
	Description string
	Content     string
	Related     []*Page
	RelatedFQNs []FQN
	Template    string
	Hide        bool
}

type IndexPage struct {
	Page
	Pages []*Page
}

type FQN string
