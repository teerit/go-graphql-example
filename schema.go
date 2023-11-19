package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/graphql-go/graphql"
)

// Graphql Schema
// type Beast {
// 	name: String!
// 	description String
// 	id: Int!
// 	imageUrl: String
// 	otherNames: [String!]
// }

// define schema, with our rootQuery
var BeastSchema, _ = graphql.NewSchema(graphql.SchemaConfig{
	Query:    rootQuery,
	Mutation: rootMutation,
})

var beastType = graphql.NewObject(graphql.ObjectConfig{
	Name: "Beast",
	Fields: graphql.Fields{
		"name": &graphql.Field{
			Type: graphql.String,
		},
		"description": &graphql.Field{
			Type: graphql.String,
		},
		"id": &graphql.Field{
			Type: graphql.Int,
		},
		"otherNames": &graphql.Field{
			Type: graphql.NewList(graphql.String),
		},
		"imageUrl": &graphql.Field{
			Type: graphql.String,
		},
	},
})

var rootQuery = graphql.NewObject(graphql.ObjectConfig{
	Name: "RootQuery",
	Fields: graphql.Fields{
		"beast": &graphql.Field{
			Type:        beastType,
			Description: "Get single beast",
			Args: graphql.FieldConfigArgument{
				"name": &graphql.ArgumentConfig{
					Type: graphql.String,
				},
			},
			Resolve: func(params graphql.ResolveParams) (interface{}, error) {

				nameQuery, isOk := params.Args["name"].(string)
				if isOk {
					// search for el with name
					for _, beast := range BeastList {
						if beast.Name == nameQuery {
							return beast, nil
						}
					}
				}
				return Beast{}, nil
			},
		},
		"beastList": &graphql.Field{
			Type:        graphql.NewList(beastType),
			Description: "List of beasts",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				return BeastList, nil
			},
		},
	},
})

var currentMaxId = 10
var rootMutation = graphql.NewObject(graphql.ObjectConfig{
	Name: "RootMutation",
	Fields: graphql.Fields{
		"addBeast": &graphql.Field{
			Type:        beastType,
			Description: "add new beast",
			Args: graphql.FieldConfigArgument{
				"name": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.String),
				},
				"description": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.String),
				},
				"otherNames": &graphql.ArgumentConfig{
					Type: graphql.NewList(graphql.String),
				},
				"imageUrl": &graphql.ArgumentConfig{
					Type: graphql.String,
				},
			},
			Resolve: func(params graphql.ResolveParams) (interface{}, error) {
				// marshall and cast the argument value
				name, _ := params.Args["name"].(string)
				description, _ := params.Args["description"].(string)
				otherNames, _ := params.Args["otherNames"].([]string)
				imageUrl, _ := params.Args["imageUrl"].(string)

				newID := currentMaxId + 1
				currentMaxId = currentMaxId + 1

				newBeast := Beast{
					ID:          newID,
					Name:        name,
					Description: description,
					OtherNames:  otherNames,
					ImageUrl:    imageUrl,
				}

				BeastList = append(BeastList, newBeast)
				return newBeast, nil
			},
		},
	},
})

func importJSONDataFromFile(fileName string, result interface{}) (isOk bool) {
	content, err := os.ReadFile(fileName)
	if err != nil {
		fmt.Println("Error: ", err)
		return false
	}

	err = json.Unmarshal(content, result)
	if err != nil {
		fmt.Println("Error: ", err)
		return false
	}
	return true
}

var BeastList []Beast
var _ = importJSONDataFromFile("./beastData.json", &BeastList)

type Beast struct {
	ID          int      `json:"id"`
	Name        string   `json:"name"`
	Description string   `json:"description"`
	OtherNames  []string `json:"otherNames"`
	ImageUrl    string   `json:"imageUrl"`
}
