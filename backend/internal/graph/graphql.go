package graph

import (
	"github.com/Snow-00/go-react-movies-backend/internal/models"
	"github.com/graphql-go/graphql"
)

type Graph struct {
	Movies      []*models.Movie
	QueryString string
	Config      graphql.SchemaConfig
	Fields      graphql.Fields
	MovieType   *graphql.Object
}

func New(movies []*models.Movie) *Graph {
	movieType := graphql.NewObject(
		graphql.ObjectConfig{
			Name: "Movie",
			Fields: graphql.Fields{
				"id": &graphql.Field{
					Type: graphql.Int,
				},
				"title": &graphql.Field{
					Type: graphql.String,
				},
				"description": &graphql.Field{
					Type: graphql.String,
				},
				"release_date": &graphql.Field{
					Type: graphql.DateTime,
				},
				"runtime": &graphql.Field{
					Type: graphql.Int,
				},
				"mpaa_rating": &graphql.Field{
					Type: graphql.String,
				},
				"created_at": &graphql.Field{
					Type: graphql.DateTime,
				},
				"updated_at": &graphql.Field{
					Type: graphql.DateTime,
				},
				"image": &graphql.Field{
					Type: graphql.String,
				},
			},
		},
	)
}
