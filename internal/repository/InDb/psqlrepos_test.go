package InDb_test

import (
	data "OzonTestovoe/internal/dataStructs"
	"OzonTestovoe/internal/repository/InDb"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/require"
)

func newMock(t *testing.T) (*InDb.PsqlRepos, sqlmock.Sqlmock, func()) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)

	repo := InDb.NewPsqlRepos(db)
	cleanup := func() {
		db.Close()
	}
	return repo, mock, cleanup
}

func TestGetPostById(t *testing.T) {
	repo, mock, cleanup := newMock(t)
	defer cleanup()

	rows := sqlmock.NewRows([]string{"id", "title", "content", "comments_are_allowed"}).
		AddRow(1, "title1", "content1", true)

	mock.ExpectQuery(`select id, title, content, comments_are_allowed from posts where id = \$1`).
		WithArgs(1).
		WillReturnRows(rows)

	post, err := repo.GetPostById(1)
	require.NoError(t, err)
	require.Equal(t, 1, post.Postid)
	require.Equal(t, "title1", post.Posttitle)
	require.Equal(t, "content1", post.Postcontent)
	require.True(t, post.Commentsareallowed)

	require.NoError(t, mock.ExpectationsWereMet())
}

func TestCreatePost(t *testing.T) {
	repo, mock, cleanup := newMock(t)
	defer cleanup()

	post := data.Post{
		Posttitle:          "new title",
		Postcontent:        "new content",
		Commentsareallowed: true,
		PostAuthor:         "author",
	}

	mock.ExpectQuery(`INSERT INTO posts`).
		WithArgs(post.Posttitle, post.Postcontent, post.Commentsareallowed, post.PostAuthor).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(10))

	created, err := repo.CreatePost(post)
	require.NoError(t, err)
	require.Equal(t, 10, created.Postid)
	require.Equal(t, "author", created.PostAuthor)

	require.NoError(t, mock.ExpectationsWereMet())
}

func TestCreateComment(t *testing.T) {
	repo, mock, cleanup := newMock(t)
	defer cleanup()

	comment := data.Comment{
		Commentcontent: "hi",
		Postid:         1,
		ParentId:       nil,
		CommentAuthor:  "tester",
	}

	mock.ExpectQuery(`INSERT INTO comments`).
		WithArgs(comment.Commentcontent, comment.Postid, comment.ParentId, comment.CommentAuthor).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(5))

	created, err := repo.CreateComment(comment)
	require.NoError(t, err)
	require.Equal(t, 5, created.Commentid)
	require.Equal(t, "tester", created.CommentAuthor)

	require.NoError(t, mock.ExpectationsWereMet())
}
