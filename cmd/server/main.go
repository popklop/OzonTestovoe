package main

import (
	cfg "OzonTestovoe/config"
	data "OzonTestovoe/internal/dataStructs"
	"OzonTestovoe/internal/database"
	gql "OzonTestovoe/internal/graphql"
	repository2 "OzonTestovoe/internal/repository/InDb"
	"OzonTestovoe/internal/repository/inmemory"
	"fmt"
	"github.com/graphql-go/graphql"
	"github.com/graphql-go/handler"
	"log"
	"net/http"
)

func main() {

	var repo inmemory.GetData
	var schema graphql.Schema
	conf, err := cfg.LoadConfig()
	if err != nil {
		log.Fatal(err)
	}
	if conf.Storage == "postgres" {
		db := database.Connect()
		fmt.Println("Database connection:", db)
		repopsql := repository2.NewPsqlRepos(db)
		schema, err = gql.Schema(repopsql)
		if err != nil {
			log.Fatal(err)
		}

	} else {
		repo = &inmemory.InMemoryDataAllign{}
		schema, err = gql.Schema(repo)
		if err != nil {
			log.Fatal(err)
		}
		initializeTestData(repo.(*inmemory.InMemoryDataAllign))
	}
	handle := handler.New(&handler.Config{
		Schema:   &schema,
		Pretty:   true,
		GraphiQL: true,
	})
	http.Handle("/graphql", handle)
	log.Println("Server started on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func initializeTestData(repo *inmemory.InMemoryDataAllign) {
	repo.Posts = []data.Post{
		{Postid: 1, Posttitle: "Go и GraphQL", Postcontent: "Как связать Go с GraphQL и зачем это нужно", Commentsareallowed: true},
		{Postid: 2, Posttitle: "Docker для новичков", Postcontent: "Базовые команды и workflow", Commentsareallowed: true},
		{Postid: 3, Posttitle: "CI/CD в GitHub Actions", Postcontent: "Настраиваем автоматический деплой", Commentsareallowed: true},
		{Postid: 4, Posttitle: "Микросервисы на Go", Postcontent: "Подходы к проектированию и примеры", Commentsareallowed: true},
		{Postid: 5, Posttitle: "PostgreSQL советы", Postcontent: "Оптимизация запросов и индексов", Commentsareallowed: false},
	}

	repo.Comments = []data.Comment{
		{Commentid: 1, Commentcontent: "Отличное введение!", Postid: 1},
		{Commentid: 2, Commentcontent: "А как насчёт аутентификации?", Postid: 1},
		{Commentid: 3, Commentcontent: "Жду продолжения серии.", Postid: 1},

		{Commentid: 4, Commentcontent: "Очень помогло с настройкой Docker!", Postid: 2},
		{Commentid: 5, Commentcontent: "Добавьте раздел про docker-compose.", Postid: 2},
		{Commentid: 6, Commentcontent: "Спасибо за понятные примеры.", Postid: 2},

		{Commentid: 7, Commentcontent: "А можно ли деплоить на Heroku?", Postid: 3},
		{Commentid: 8, Commentcontent: "Работает даже на Windows runner'ах!", Postid: 3},
		{Commentid: 9, Commentcontent: "CI/CD — must have.", Postid: 3},

		{Commentid: 10, Commentcontent: "Поясните про gRPC.", Postid: 4},
		{Commentid: 11, Commentcontent: "Хотелось бы пример с Kubernetes.", Postid: 4},
		{Commentid: 12, Commentcontent: "Сложно, но интересно.", Postid: 4},

		{Commentid: 13, Commentcontent: "Оптимизация запросов — супер.", Postid: 5},
		{Commentid: 14, Commentcontent: "Что насчёт partitioning?", Postid: 5},
		{Commentid: 15, Commentcontent: "Добавьте про jsonb индексы.", Postid: 5},

		{Commentid: 16, Commentcontent: "Классная статья!", Postid: 1},
		{Commentid: 17, Commentcontent: "Хочу больше примеров кода.", Postid: 2},
		{Commentid: 18, Commentcontent: "GraphQL — топ.", Postid: 1},
		{Commentid: 19, Commentcontent: "CI/CD спасает жизнь.", Postid: 3},
		{Commentid: 20, Commentcontent: "Go forever!", Postid: 4},
	}
}
