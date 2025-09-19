package graphql

import (
	data "OzonTestovoe/internal/dataStructs"
	"OzonTestovoe/internal/repository/inmemory"
	"errors"
	"github.com/graphql-go/graphql"
)

func RootMutation(repos inmemory.GetData) *graphql.Object {
	obj := graphql.NewObject(graphql.ObjectConfig{
		Name: "Mutation",
		Fields: graphql.Fields{
			"createComment": &graphql.Field{
				Type: NewCommentType(repos),
				Args: graphql.FieldConfigArgument{
					"content": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.String),
					},
					"postId": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.Int),
					},
					"parentId": &graphql.ArgumentConfig{
						Type: graphql.Int,
					},
					"commentAuthor": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.String),
					},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					content := p.Args["content"].(string)
					if len(content) == 0 || len(content) > 2000 {
						return nil, errors.New("invalid input, check your comment length")
					}
					postId := p.Args["postId"].(int)
					var comParentId *int
					if parentVal, ok := p.Args["parentId"].(int); ok {
						comParentId = &parentVal
						parentComment, err := repos.GetCommById(*comParentId)
						if err != nil || parentComment.Postid != postId {
							return nil, errors.New("parent comment not found or belongs to another post")
						}
					}
					post, err := repos.GetPostById(postId)
					if err != nil {
						return nil, errors.New("invalid postId")
					}
					if !post.Commentsareallowed {
						return nil, errors.New("comments are not allowed")
					}
					newComm := data.Comment{
						Commentcontent: content,
						Postid:         postId,
						ParentId:       comParentId,
						CommentAuthor:  p.Args["commentAuthor"].(string),
					}
					return repos.CreateComment(newComm)
				},
			},
			"createPost": &graphql.Field{
				Type: NewPostType(repos),
				Args: graphql.FieldConfigArgument{
					"title": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.String),
					},
					"content": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.String),
					},
					"commentsareallowed": &graphql.ArgumentConfig{
						Type:         graphql.NewNonNull(graphql.Boolean),
						DefaultValue: true,
					},
					"postAuthor": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.String),
					},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					title := p.Args["title"].(string)
					content := p.Args["content"].(string)
					commentsareallowed := p.Args["commentsareallowed"].(bool)
					postAuthor := p.Args["postAuthor"].(string)
					post := &data.Post{
						Posttitle:          title,
						Postcontent:        content,
						Commentsareallowed: commentsareallowed,
						PostAuthor:         postAuthor,
					}
					return repos.CreatePost(*post)
				},
			},
		},
	})
	return obj
}
