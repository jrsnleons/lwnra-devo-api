package models

// FBPost represents a Facebook post from the API
type FBPost struct {
	Message     string `json:"message"`
	CreatedTime string `json:"created_time"`
}

// FBPosts represents a collection of Facebook posts
type FBPosts struct {
	Data []FBPost `json:"data"`
}

// FBMeResponse represents the Facebook API response for the /me endpoint
type FBMeResponse struct {
	ID    string  `json:"id"`
	Name  string  `json:"name"`
	Posts FBPosts `json:"posts"`
}
