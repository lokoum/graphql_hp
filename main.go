package main

import (
	"log"
	"net/http"
	"fmt"

	_ "graphql_hp/db"

	"graphql_hp/resolvers"
	"github.com/neelance/graphql-go"
	"github.com/neelance/graphql-go/relay"
	"io/ioutil"
	"strings"
	"context"
)

var schema *graphql.Schema

func formatRequest(r *http.Request) string {
	// Create return string
	var request []string
	// Add the request string
	url := fmt.Sprintf("%v %v %v", r.Method, r.URL, r.Proto)
	request = append(request, url)
	// Add the host
	request = append(request, fmt.Sprintf("Host: %v", r.Host))
	// Loop through headers
	for name, headers := range r.Header {
	  name = strings.ToLower(name)
	  for _, h := range headers {
		request = append(request, fmt.Sprintf("%v: %v", name, h))
	  }
	}
	
	// If this is a POST, add post data
	if r.Method == "POST" {
	   r.ParseForm()
	   request = append(request, "\n")
	   request = append(request, r.Form.Encode())
	} 
	 // Return the request as a string
	 return strings.Join(request, "\n")
   }

func init() {
	schemaFile, err := ioutil.ReadFile("schema.graphqls")
	if err != nil {
		panic(err)
	}

	schema = graphql.MustParseSchema(string(schemaFile), &resolvers.Resolver{})
}

func main() {
	http.Handle("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write(page)
	}))

	http.Handle("/query", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//log.Println(r.)
		formatRequest(r)
		next := &relay.Handler{Schema: schema}
		authorization := r.Header.Get("Authorization")
		token := strings.Replace(authorization, "Bearer ", "", 1)
		ctx := context.WithValue(r.Context(), "AuthorizationToken", token)
		next.ServeHTTP(w, r.WithContext(ctx))
	}))

	log.Println("server is running on port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

var page = []byte(`
<!DOCTYPE html>
<html>
	<head>
		<link rel="stylesheet" href="https://cdnjs.cloudflare.com/ajax/libs/graphiql/0.10.2/graphiql.css" />
		<script src="https://cdnjs.cloudflare.com/ajax/libs/fetch/1.1.0/fetch.min.js"></script>
		<script src="https://cdnjs.cloudflare.com/ajax/libs/react/15.5.4/react.min.js"></script>
		<script src="https://cdnjs.cloudflare.com/ajax/libs/react/15.5.4/react-dom.min.js"></script>
		<script src="https://cdnjs.cloudflare.com/ajax/libs/graphiql/0.10.2/graphiql.js"></script>
	</head>
	<body style="width: 100%; height: 100%; margin: 0; overflow: hidden;">
		<div id="graphiql" style="height: 100vh;">Loading...</div>
		<script>
			function graphQLFetcher(graphQLParams) {
				return fetch("/query", {
					method: "post",
					headers: {
						'Accept': 'application/json',
						'Content-Type': 'application/json',
						'Authorization': 'Bearer f66aa7b3-add4-4fa2-88c8-534b2905c6ca'
					},
					body: JSON.stringify(graphQLParams),
					credentials: "include",
				}).then(function (response) {
					return response.text();
				}).then(function (responseBody) {
					try {
						return JSON.parse(responseBody);
					} catch (error) {
						return responseBody;
					}
				});
			}
			ReactDOM.render(
				React.createElement(GraphiQL, {fetcher: graphQLFetcher}),
				document.getElementById("graphiql")
			);
		</script>
	</body>
</html>
`)
