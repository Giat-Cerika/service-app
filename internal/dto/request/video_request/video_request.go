package videorequest

type CreateVideoRequest struct {
	VideoPath   string `form:"video_path" json:"video_path"`
	Title       string `form:"title" json:"title"`
	Description string `form:"description" json:"description"`
}

type UpdateVideoRequest struct {
	VideoPath   string `form:"video_path" json:"video_path"`
	Title       string `form:"title" json:"title"`
	Description string `form:"description" json:"description"`
}
