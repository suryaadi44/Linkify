package entity

type Links struct {
	Username    string `bson:"_id"`
	Description string `bson:"description"`
	Links       []Link `bson:"links,omitempty"`
}

type Link struct {
	Title     string `bson:"title"`
	URL       string `bson:"url"`
	Thumbnail string `bson:"thumbnail,omitempty"`
}
