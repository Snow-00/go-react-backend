package graph

import (
	"strings"

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
	// object definition
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

	// specify the actions
	fields := graphql.Fields{
		"list": &graphql.Field{
			Type:        graphql.NewList(movieType),
			Description: "Get all movies",
			Resolve: func(params graphql.ResolveParams) (any, error) {
				return movies, nil
			},
		},

		"search": &graphql.Field{
			Type:        graphql.NewList(movieType),
			Description: "search movies by title",
			Args: graphql.FieldConfigArgument{
				"titleContains": &graphql.ArgumentConfig{
					Type: graphql.String,
				},
			},
			Resolve: func(params graphql.ResolveParams) (any, error) {
				var theList []*models.Movie
				search, ok := params.Args["titleContains"].(string)
				if ok {
					for _, currentMovie := range movies {
						if strings.Contains(strings.ToLower(currentMovie.Title), strings.ToLower(search)) {
							theList = append(theList, currentMovie)
						}
					}
				}
				return theList, nil
			},
		},

		"get": &graphql.Field{
			Type:        movieType,
			Description: "Get movie by id",
			Args: graphql.FieldConfigArgument{
				"id": &graphql.ArgumentConfig{
					Type: graphql.Int,
				},
			},
			Resolve: func(p graphql.ResolveParams) (any, error) {
				id, ok := p.Args["id"].(int)
				if ok {
					for _, movie := range movies {
						if movie.ID == id {
							return movie, nil
						}
					}
				}
				return nil, nil
			},
		},
	}

	return &Graph{
		Movies:    movies,
		Fields:    fields,
		MovieType: movieType,
	}
}
