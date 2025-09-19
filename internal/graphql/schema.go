package graphql

import (
	"OzonTestovoe/internal/repository/inmemory"
	"github.com/graphql-go/graphql"
)

func Schema(repos inmemory.GetData) (graphql.Schema, error) {
	post := NewPostType(repos)
	comment := NewCommentType(repos)
	rootQuery := graphql.NewObject(graphql.ObjectConfig{
		Name: "RootQuery",
		Fields: graphql.Fields{
			"posts": &graphql.Field{
				Type: graphql.NewList(post),
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					return repos.GetAllPosts()
				},
			},
			"post": &graphql.Field{
				Type: post,
				Args: graphql.FieldConfigArgument{
					"id": &graphql.ArgumentConfig{Type: graphql.Int},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					id := p.Args["id"].(int)
					return repos.GetPostById(id)
				},
			},
			"comment": &graphql.Field{
				Type: graphql.NewList(comment),
				Args: graphql.FieldConfigArgument{
					"postId": &graphql.ArgumentConfig{Type: graphql.Int},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					id := p.Args["postId"].(int)
					return repos.GetCommByPostId(id)
				},
			},
			"postsbyId": &graphql.Field{
				Type: graphql.NewList(post),
				Args: graphql.FieldConfigArgument{
					"postIds": &graphql.ArgumentConfig{Type: graphql.NewList(graphql.Int)},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					id := p.Args["postIds"].([]interface{})
					ids := make([]int, len(id))
					for i, v := range id {
						ids[i] = v.(int)
					}
					return repos.GetPostsById(ids)
				},
			},
		},
	})
	rootmut := RootMutation(repos)
	return graphql.NewSchema(graphql.SchemaConfig{
		Query:    rootQuery,
		Mutation: rootmut,
	})
}
