package classrequest

type CreateClassRequest struct {
	NameClass string `form:"name_class" json:"name_class"`
	Grade     string `form:"grade" json:"grade"`
	Teacher   string `form:"teacher" json:"teacher"`
}

type UpdateClassRequest struct {
	NameClass string `form:"name_class" json:"name_class"`
	Grade     string `form:"grade" json:"grade"`
	Teacher   string `form:"teacher" json:"teacher"`
}
