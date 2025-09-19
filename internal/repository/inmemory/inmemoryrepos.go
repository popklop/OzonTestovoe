package inmemory

import (
	data "OzonTestovoe/internal/dataStructs"
	"errors"
	"fmt"
	"sync"
)


type GetData interface {
	GetPostsById(id []int) ([]*data.Post, error)
	GetCommByPostId(postid int) ([]data.Comment, error)
	GetAllPosts() ([]*data.Post, error)
	CreateComment(comment data.Comment) (*data.Comment, error)
	CreatePost(post data.Post) (*data.Post, error)
	GetCommentReplies(postid int) ([]data.Comment, error)
	GetCommById(commid int) (*data.Comment, error)
	GetPostById(postid int) (data.Post, error)
}

type InMemoryDataAllign struct {
	mu       sync.RWMutex
	Posts    []data.Post
	Comments []data.Comment
}

func (d *InMemoryDataAllign) GetPostsById(ids []int) ([]*data.Post, error) {
	d.mu.RLock()
	defer d.mu.RUnlock()
	res := []*data.Post{}
	for _, id := range ids {
		post, err := d.GetPostById(id)
		if err == nil {
			res = append(res, &post)
		}
	}
	if len(res) == 0 {
		return nil, errors.New("couldn't find any post by these ids")
	}
	return res, nil
}

func (d *InMemoryDataAllign) GetCommByPostId(postId int) ([]data.Comment, error) {
	d.mu.RLock()
	defer d.mu.RUnlock()

	var comments []data.Comment
	for _, postcomms := range d.Comments {
		if postId == postcomms.Postid && postcomms.ParentId == nil {
			comments = append(comments, postcomms)
		}
	}
	return comments, nil
}

func (d *InMemoryDataAllign) GetCommById(commentId int) (*data.Comment, error) {
	d.mu.RLock()
	defer d.mu.RUnlock()

	for _, commid := range d.Comments {
		if commid.Commentid == commentId {
			return &commid, nil
		}
	}
	return nil, errors.New("comment not found")
}

func (d *InMemoryDataAllign) GetPostById(id int) (data.Post, error) {
	d.mu.RLock()
	defer d.mu.RUnlock()

	for _, post := range d.Posts {
		if post.Postid == id {
			return post, nil
		}
	}
	return data.Post{}, errors.New("not found")
}

func (d *InMemoryDataAllign) GetAllPosts() ([]*data.Post, error) {
	d.mu.RLock()
	defer d.mu.RUnlock()

	if len(d.Posts) == 0 {
		return nil, errors.New("no posts found")
	}

	res := make([]*data.Post, 0, len(d.Posts))
	for i := range d.Posts {
		res = append(res, &d.Posts[i])
	}
	return res, nil
}

func (d *InMemoryDataAllign) GetCommentReplies(parentId int) ([]data.Comment, error) {
	d.mu.RLock()
	defer d.mu.RUnlock()
	var replies []data.Comment
	for _, comment := range d.Comments {
		if comment.ParentId != nil && *comment.ParentId == parentId {
			replies = append(replies, comment)
		}
	}
	return replies, nil
}

func (d *InMemoryDataAllign) CreateComment(comment data.Comment) (*data.Comment, error) {
	d.mu.Lock()
	defer d.mu.Unlock()
	post, err := d.GetPostById(comment.Postid)
	if err != nil {
		return nil, err
	}
	if !post.Commentsareallowed {
		return nil, errors.New("comments are not allowed")
	}
	if comment.ParentId != nil {
		parentExists := false
		for _, c := range d.Comments {
			if c.Commentid == *comment.ParentId {
				parentExists = true
				break
			}
		}
		if !parentExists {
			return nil, fmt.Errorf("parent comment not found")
		}
	}
	comment.Commentid = len(d.Comments) + 1
	d.Comments = append(d.Comments, comment)
	return &comment, nil
}

func (d *InMemoryDataAllign) CreatePost(post data.Post) (*data.Post, error) {
	d.mu.Lock()
	defer d.mu.Unlock()
	post.Postid = len(d.Posts) + 1
	d.Posts = append(d.Posts, post)
	return &post, nil
}
