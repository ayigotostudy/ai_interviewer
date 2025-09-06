package req

type CreateWikiRequest struct {
	UserID   uint   `json:"user_id" form:"user_id"`
	Title    string `json:"title" form:"title"`
	Content  string `json:"content" form:"content"`
	ParentID uint   `json:"parent_id" form:"parent_id"`
	WikiType int    `json:"wiki_type" form:"wiki_type"`
	Type     int    `json:"type" form:"type"`
	RootId   uint   `json:"root_id" form:"root_id"`
	Url      string `json:"url" form:"url"`
}

type GetWikiListRequest struct {
	UserID uint `json:"user_id" form:"user_id"`
	ParentID uint `json:"parent_id" form:"parent_id"`
}

type GetWikiRequest struct {
	ID     uint `json:"id" form:"id"`
	UserID uint `json:"user_id" form:"user_id"`
}

type DeleteWikiRequest struct {
	ID     uint `json:"id" form:"id"`
	UserID uint `json:"user_id" form:"user_id"`
}

type UpdateWikiRequest struct {
	ID      uint   `json:"id" form:"id" binding:"required"`
	Title   string `json:"title" form:"title"`
	Content string `json:"content" form:"content"`
}

type QueryWikiRequest struct {
	Query  string `json:"query" form:"query"`
	UserID uint   `json:"user_id" form:"user_id"`
	RootId uint   `json:"root_id" form:"root_id"`
}
