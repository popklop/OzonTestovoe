package InDb

import (
	data "OzonTestovoe/internal/dataStructs"
	"database/sql"
	"errors"
	"fmt"
	"strings"
)

type PsqlRepos struct {
	psqlDb *sql.DB
}

func NewPsqlRepos(psqlDb *sql.DB) *PsqlRepos {
	return &PsqlRepos{psqlDb}
}

func (d *PsqlRepos) GetPostById(id int) (data.Post, error) {
	var post data.Post
	err := d.psqlDb.QueryRow(
		`select id, title, content, comments_are_allowed from posts where id = $1`, id).
		Scan(&post.Postid, &post.Posttitle, &post.Postcontent, &post.Commentsareallowed)
	return post, err
}

func (d *PsqlRepos) GetAllPosts() ([]*data.Post, error) {
	postasrr := []*data.Post{}
	rows, err := d.psqlDb.Query(` select id, title, content, comments_are_allowed from posts`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var post data.Post
		if err := rows.Scan(&post.Postid, &post.Posttitle, &post.Postcontent, &post.Commentsareallowed); err != nil {
			return nil, err
		}
		postasrr = append(postasrr, &post)
	}
	return postasrr, nil
}

func (d *PsqlRepos) GetCommByPostId(postId int) ([]data.Comment, error) {
	rows, err := d.psqlDb.Query(`SELECT id, content, post_id, parent_id FROM comments WHERE post_id=$1`, postId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	commmas := []data.Comment{}
	for rows.Next() {
		var comm data.Comment
		if err := rows.Scan(&comm.Commentid, &comm.Commentcontent, &comm.Postid, &comm.ParentId); err != nil {
			return nil, err
		}
		commmas = append(commmas, comm)
	}
	return commmas, nil
}

func (d *PsqlRepos) GetCommById(commId int) (*data.Comment, error) {
	var comm data.Comment
	err := d.psqlDb.QueryRow(`SELECT id, content, post_id, parent_id FROM comments WHERE id=$1`, commId).
		Scan(&comm.Commentid, &comm.Commentcontent, &comm.Postid, &comm.ParentId)
	if err != nil {
		return nil, err
	}
	return &comm, nil
}

func (d *PsqlRepos) GetCommentReplies(parentId int) ([]data.Comment, error) {
	rows, err := d.psqlDb.Query(
		`SELECT id, content, post_id, parent_id FROM comments WHERE parent_id=$1`, parentId,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	replies := []data.Comment{}
	for rows.Next() {
		comm := data.Comment{}
		if err := rows.Scan(&comm.Commentid, &comm.Commentcontent, &comm.Postid, &comm.ParentId); err != nil {
			return nil, err
		}
		replies = append(replies, comm)
	}
	return replies, nil
}

func (d *PsqlRepos) CreatePost(post data.Post) (*data.Post, error) {
	err := d.psqlDb.QueryRow(`INSERT INTO posts (title, content,comments_are_allowed, post_author)
         VALUES ($1,$2,$3, $4) RETURNING id`, post.Posttitle, post.Postcontent, post.Commentsareallowed, post.PostAuthor).Scan(&post.Postid)
	if err != nil {
		return nil, err
	}
	return &post, nil
}

func (d *PsqlRepos) CreateComment(comment data.Comment) (*data.Comment, error) {
	err := d.psqlDb.QueryRow(`
        INSERT INTO comments (content, post_id, parent_id, comment_author)
        VALUES ($1, $2, $3, $4)
        RETURNING id`, comment.Commentcontent, comment.Postid, comment.ParentId, comment.CommentAuthor,
	).Scan(&comment.Commentid)

	if err != nil {
		return nil, err
	}
	return &comment, nil
}
func (r *PsqlRepos) GetPostsById(ids []int) ([]*data.Post, error) {
	if len(ids) == 0 {
		return nil, errors.New("empty ids")
	}
	params := []interface{}{}
	placeholders := []string{}
	for i, id := range ids {
		params = append(params, id)
		placeholders = append(placeholders, fmt.Sprintf("$%d", i+1))
	}
	query := fmt.Sprintf(
		`SELECT id, title, content
		 FROM posts
		 WHERE id IN (%s)`,
		strings.Join(placeholders, ","),
	)
	rows, err := r.psqlDb.Query(query, params...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var res []*data.Post
	for rows.Next() {
		var post data.Post
		if err := rows.Scan(&post); err != nil {
			return nil, err
		}
		res = append(res, &post)
	}
	if len(res) == 0 {
		return nil, errors.New("couldn't find any post by these ids")
	}
	return res, nil
}
