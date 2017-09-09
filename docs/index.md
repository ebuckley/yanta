# Yanta's yet another note taking application
![](https://media.giphy.com/media/wlEPdDuMQzSkE/giphy.gif)

Like todo app's note taking software is highly at risk of [NIH](https://en.wikipedia.org/wiki/Not_invented_here) syndrome. 
I made the decision to start this project for 3 reasons.

1. I needed some weekend code escapism. This project runs solely on the Minimum Viable, doesn't annoy me method of product management.  I'm using it to explore technology I find interesting.
2. I love bear writer, but I'm not ready to pay for the premium service. Instead, I'll implement the premium parts in Yanta.
3. I want a tool that works really well with existing repo's of markdown (think documentation/gitbooks)

# features
- index page (autorefresh)
- auto commit
- pdf download
- view page
- sync notes with git
- create new files

# Roadmap
See [/pages/docs/todo.md](/pages/docs/todo.md)

# usage 
## Pdf export
This is a total hack right now and needs some TLC. Basically, you're gonna need to have https://github.com/ebuckley/pdf-now running. Then you're gonna need to hack the code in the pdf package so that it has the correct path as it's currently hard coded to where that project is running on my system :joy:

This will probably be fixed as soon as one other person wants to use this project :joy:

## development
Assuming you have a correctly configured `GOPATH` with the required dependencies.

```
cd src/code/github.com/ebuckley/yanta && go run main.go
```

The development instance will be running on port 1337, everything under the 

## production
```
./yanta -dir <pathname>
```

Provide a directory for serving the site. The service will run on `0.0.0.0:1337`

The way pull/push works by default is by making the assumption that the target directory has a git repo and a "publish" remote. It will push master to "publish", and sync master from "publish".

# bugs/features/requests/planning
[TODO](/page/docs/todo.md)

Put content in the todo here


# Development dependencies
A common pattern in golang projects is dependency vendoring. I'm not doing that yet as the project is in it's initial stages and moving at a pace which warrents keeping dependenceis on the bleeding edge,

-	"github.com/gorilla/schema" is used for deserialization of form requests to structs
-	"github.com/gorilla/mux" is used for routing
-	"github.com/golang-commonmark/markdown" is used for markdown parsing
