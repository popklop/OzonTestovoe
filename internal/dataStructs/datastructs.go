package dataStructs

type Post struct {
	Postid             int    `json:"id"`
	Posttitle          string `json:"title"`
	Postcontent        string `json:"content"`
	Commentsareallowed bool   `json:"commentsareallowed"`
	PostAuthor         string `json:"postauthor"`
}

type Comment struct {
	Commentid      int    `json:"id"`
	Commentcontent string `json:"content"`
	Postid         int    `json:"postId"`
	ParentId       *int   `json:"parentId"`
	CommentAuthor  string `json:"commentAuthor"`
}
