package models

type Image struct {
	ID      string `bson:"_id,omitempty" json:"id"`
	Name    string `bson:"name" json:"name"`
	Url     string `bson:"url" json:"url"`
	UserID  string `bson:"user_id" json:"user_id"`
	PostID  string `bson:"post_id" json:"post_id"`
	StoryID string `bson:"story_id" json:"story_id"`
}
