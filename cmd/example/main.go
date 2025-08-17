package main

import (
	"context"
	"log"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humachi"
	"github.com/go-chi/chi/v5"
	"github.com/grafana/pyroscope-go"
	"github.com/j-dumbell/humascope"
)

type GreetingOutput struct {
	Body struct {
		Message string `json:"message" example:"Hello, world!" doc:"Greeting message"`
	}
}

type ReviewInput struct {
	Body struct {
		Author  string `json:"author" maxLength:"10" doc:"Author of the review"`
		Rating  int    `json:"rating" minimum:"1" maximum:"5" doc:"Rating from 1 to 5"`
		Message string `json:"message,omitempty" maxLength:"100" doc:"Review message"`
	}
}

func main() {
	router := chi.NewMux()
	api := humachi.New(router, huma.DefaultConfig("My API", "1.0.0"))

	// Register GET /greeting/{name} handler with profiling labels:
	//  method = GET
	//  path = /greeting/{name}
	humascope.Register(
		api,
		huma.Operation{
			OperationID: "get-greeting",
			Method:      http.MethodGet,
			Path:        "/greeting/{name}",
		},
		func(ctx context.Context, input *struct {
			Name string `path:"name" maxLength:"30" example:"world" doc:"Name to greet"`
		}) (*GreetingOutput, error) {
			resp := &GreetingOutput{}
			return resp, nil
		},
	)

	// Register handler with custom profiling labels via humascope.NewPyroscopeMW.
	huma.Register(api, huma.Operation{
		OperationID: "post-review",
		Method:      http.MethodPost,
		Path:        "/reviews",
		Middlewares: huma.Middlewares{humascope.NewPyroscopeMW("custom_label", "custom_value")},
	}, func(ctx context.Context, i *ReviewInput) (*struct{}, error) {
		// implementation omitted
		return nil, nil
	})

	log.Println("starting Pyroscope")
	_, err := pyroscope.Start(pyroscope.Config{
		ApplicationName: "example",
		ServerAddress:   "http://localhost:4040",
	})
	if err != nil {
		log.Fatal("error starting Pyroscope:", err)
	}

	log.Println("starting webserver on port 8888")
	http.ListenAndServe("localhost:8888", router)
}
