# Cool Language Features
## Brendan Horng

Website for Parallel Processing Course creating a website.

## Getting Started

### Running

You can run the code inside `cool-lang-features/`:

```
go run frontend.go [--listen {port} --backend {backendAddress}]
go run backend.go [--listen {port}]
```

### Testing

The GoLang code has pretty good test coverage, so you can run the tests that I wrote or manually test it by spinning up the server.
Currently the routes are as follows:
- `/`: Serves the listing of all the features. From here you can navigate to delete, modify, or add new features
- `/add`: Form to add a new feature
- `/modify/{id}`: Form to modify a feature identified by id
- `/api/feature`: REST API defined here for GET and POST commands to retrieve all features or to add a new one
- `/api/feature/{id}`: REST API defined her for GET, PATCH, and DELETE to retrieve, modify, or delete specific features based on id

## Other Notes

### Current Status

The front end code is a little more unclean than I would like it to be, but the overall functionality is there. The web server, written in Go, is capable of doing add, delete, read, and modify on "features".

#### Update for Part 2

The above is what I wrote for part 1 and still relevant so I left it as is. The frontend can be spun up separately from the backend (see **Running**).

I don't like how I ended up structuring the code base. I will probably refactor it in the near future. Particularly what I don't like is because both my frontend and my backend use **main** as their package, I can't do a `go build`.

To properly do this, I would probably have to separate out everything even further from the current folder hierarchy that I have. For simplicity, I left it as is for now.

### Additional Resources

- http://elm-lang.org/
    - I used Elm to create the front end of this application. I am still pretty new to it so the code code be improved.
    - You should not have to do anything with it, I built the one html file that the Go server will be hosting.

### Additional Thoughts

It took longer than I thought it would. I enjoy these kinds of projects though, and I thoroughly enjoyed it because I was using two languages I don't typically get to.

#### Update for Part 2

No additional insights.
