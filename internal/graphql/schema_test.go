package graphql_test

import (
	"OzonTestovoe/internal/dataStructs"
	mygql "OzonTestovoe/internal/graphql"
	"OzonTestovoe/internal/repository/inmemory"
	"github.com/graphql-go/graphql"
	"testing"
)

func setupSchema() graphql.Schema {
	repo := &inmemory.InMemoryDataAllign{}
	repo.Posts = []dataStructs.Post{
		{Postid: 1, Posttitle: "Post 1", Postcontent: "Content 1", Commentsareallowed: true},
		{Postid: 2, Posttitle: "Post 2", Postcontent: "Content 2", Commentsareallowed: false},
	}
	schema, err := mygql.Schema(repo)
	if err != nil {
		panic(err)
	}
	return schema
}

func TestGraphQLQueries(t *testing.T) {
	schema := setupSchema()
	query := `{ posts { id title content commentsareallowed } }`
	params := graphql.Params{Schema: schema, RequestString: query}
	res := graphql.Do(params)
	if len(res.Errors) > 0 {
		t.Fatalf("unexpected errors: %v", res.Errors)
	}
	query = `{ post(id: 1) { id title content } }`
	params = graphql.Params{Schema: schema, RequestString: query}
	res = graphql.Do(params)
	if res.Data.(map[string]interface{})["post"] == nil {
		t.Fatalf("expected post data, got nil")
	}
	query = `{ comment(postId: 1) { id content } }`
	params = graphql.Params{Schema: schema, RequestString: query}
	res = graphql.Do(params)
	if len(res.Errors) > 0 {
		t.Fatalf("unexpected errors: %v", res.Errors)
	}
}

func TestGraphQLMutations(t *testing.T) {
	schema := setupSchema()
	mutation := `
	mutation {
		createComment(content: "Hello", postId: 1, commentAuthor: "Tester") {
			content postId commentAuthor
		}
	}`
	params := graphql.Params{Schema: schema, RequestString: mutation}
	res := graphql.Do(params)
	if len(res.Errors) > 0 {
		t.Fatalf("unexpected errors: %v", res.Errors)
	}
	mutation = `
	mutation {
		createComment(content: "Hello", postId: 2, commentAuthor: "Tester") {
			content
		}
	}`
	params = graphql.Params{Schema: schema, RequestString: mutation}
	res = graphql.Do(params)
	if len(res.Errors) == 0 {
		t.Fatalf("expected error, got none")
	}
	mutation = `
	mutation {
		createPost(title: "New", content: "New Content", commentsareallowed: true, postAuthor: "Author") {
			title content postAuthor
		}
	}`
	params = graphql.Params{Schema: schema, RequestString: mutation}
	res = graphql.Do(params)
	if len(res.Errors) > 0 {
		t.Fatalf("unexpected errors: %v", res.Errors)
	}
}
