package graphql

import (
	data "OzonTestovoe/internal/dataStructs"
	"OzonTestovoe/internal/repository/inmemory"
	"github.com/graphql-go/graphql"
	"sync"
)

var (
	commentType *graphql.Object
	postType    *graphql.Object
	once        sync.Once
)

func NewCommentType(repos inmemory.GetData) *graphql.Object {
	once.Do(func() {
		commentType = graphql.NewObject(graphql.ObjectConfig{
			Name: "Comment",
			Fields: graphql.Fields{
				"id": &graphql.Field{
					Type: graphql.Int,
				},
				"content": &graphql.Field{
					Type: graphql.String,
				},
				"postId": &graphql.Field{
					Type: graphql.Int,
				},
				"parentId": &graphql.Field{
					Type: graphql.Int,
				},
				"commentAuthor": &graphql.Field{
					Type: graphql.NewNonNull(graphql.String),
				},
			},
		})
		commentType.AddFieldConfig("replies", &graphql.Field{
			Type: graphql.NewList(commentType),
			Args: graphql.FieldConfigArgument{
				"limit":  &graphql.ArgumentConfig{Type: graphql.Int},
				"offset": &graphql.ArgumentConfig{Type: graphql.Int},
			},
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				comment := p.Source.(data.Comment)
				replies, err := repos.GetCommentReplies(comment.Commentid)
				if err != nil {
					return nil, err
				}
				limit, _ := p.Args["limit"].(int)
				offset, _ := p.Args["offset"].(int)
				if offset < 0 {
					offset = 0
				}
				if limit > 0 {
					end := offset + limit
					if end > len(replies) {
						end = len(replies)
					}
					if offset > len(replies) {
						return []data.Comment{}, nil
					}
					replies = replies[offset:end]
				}
				return replies, nil
			},
		})
	})
	return commentType
}
func NewPostType(repos inmemory.GetData) *graphql.Object {
	if postType == nil {
		postType = graphql.NewObject(graphql.ObjectConfig{
			Name: "Post",
			Fields: graphql.Fields{
				"id": &graphql.Field{
					Type: graphql.Int,
				},
				"title": &graphql.Field{
					Type: graphql.String,
				},
				"content": &graphql.Field{
					Type: graphql.String,
				},
				"commentsareallowed": &graphql.Field{
					Type: graphql.Boolean,
				},
				"postAuthor": &graphql.Field{
					Type: graphql.String,
				},
				"comments": &graphql.Field{
					Type: graphql.NewList(NewCommentType(repos)),
					Args: graphql.FieldConfigArgument{
						"limit":  &graphql.ArgumentConfig{Type: graphql.Int},
						"offset": &graphql.ArgumentConfig{Type: graphql.Int},
					},
					Resolve: func(p graphql.ResolveParams) (interface{}, error) {
						post := p.Source.(data.Post)
						limit, _ := p.Args["limit"].(int)
						offset, _ := p.Args["offset"].(int)
						comments, err := repos.GetCommByPostId(post.Postid)
						if err != nil {
							return nil, err
						}
						if offset < 0 {
							offset = 0
						}
						if limit > 0 {
							end := offset + limit
							if end > len(comments) {
								end = len(comments)
							}
							if offset > len(comments) {
								return []data.Comment{}, nil
							}
							comments = comments[offset:end]
						}
						return comments, nil
					},
				},
			},
		})
	}
	return postType
}
