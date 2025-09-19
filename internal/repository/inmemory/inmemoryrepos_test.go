package inmemory_test

import (
	"OzonTestovoe/internal/dataStructs"
	"OzonTestovoe/internal/repository/inmemory"

	"testing"
)

func setupRepo() *inmemory.InMemoryDataAllign {
	repo := &inmemory.InMemoryDataAllign{}
	repo.Posts = []dataStructs.Post{
		{Postid: 1, Posttitle: "Post 1", Postcontent: "Content 1", Commentsareallowed: true},
		{Postid: 2, Posttitle: "Post 2", Postcontent: "Content 2", Commentsareallowed: false},
	}
	repo.Comments = []dataStructs.Comment{
		{Commentid: 1, Commentcontent: "Comment 1", Postid: 1},
		{Commentid: 2, Commentcontent: "Reply to 1", Postid: 1, ParentId: &[]int{1}[0]},
	}
	return repo
}

func TestGetPostById(t *testing.T) {
	repo := setupRepo()
	post, err := repo.GetPostById(1)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if post.Postid != 1 {
		t.Fatalf("expected post id 1, got %d", post.Postid)
	}
	_, err = repo.GetPostById(999)
	if err == nil {
		t.Fatalf("expected error for non-existent post")
	}
}

func TestGetCommByPostId(t *testing.T) {
	repo := setupRepo()
	comments, _ := repo.GetCommByPostId(1)
	if len(comments) != 1 {
		t.Fatalf("expected 1 top-level comment, got %d", len(comments))
	}
}

func TestCreateComment(t *testing.T) {
	repo := setupRepo()
	newComment := dataStructs.Comment{Commentcontent: "New", Postid: 1}
	created, err := repo.CreateComment(newComment)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if created.Commentid == 0 {
		t.Fatalf("expected valid comment id")
	}

	// Проверка комментариев на пост, где запрещены комментарии
	newComment.Postid = 2
	_, err = repo.CreateComment(newComment)
	if err == nil {
		t.Fatalf("expected error, comments not allowed")
	}
}

func TestGetCommentReplies(t *testing.T) {
	repo := setupRepo()
	replies, _ := repo.GetCommentReplies(1)
	if len(replies) != 1 {
		t.Fatalf("expected 1 reply, got %d", len(replies))
	}
}

func TestGetPostsById(t *testing.T) {
	repo := setupRepo()
	posts, err := repo.GetPostsById([]int{1, 2})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(posts) != 2 {
		t.Fatalf("expected 2 posts, got %d", len(posts))
	}

	_, err = repo.GetPostsById([]int{999})
	if err == nil {
		t.Fatalf("expected error for invalid ids")
	}
}

func TestCreatePost(t *testing.T) {
	repo := setupRepo()
	newPost := dataStructs.Post{Posttitle: "New Post", Postcontent: "Content", Commentsareallowed: true}
	created, err := repo.CreatePost(newPost)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if created.Postid == 0 {
		t.Fatalf("expected valid post id")
	}
}
