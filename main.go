package main

import (
	"github.com/drlogout/trellocms"
	"net/http"
	"fmt"
	"text/template"
)


func main() {

	config, err := trellocms.ParseConfig()
	check(err)
	lists, err := trellocms.GetLists(config)
	check(err)

	http.Handle("/", &MyHandler{
		Config: config,
		Lists: lists,
	})
	http.ListenAndServe(":8888", nil)

}

type MyHandler struct {
	http.Handler
	Config trellocms.Config
	Lists  trellocms.Lists
}

type Context struct {
	Cards []trellocms.Card
	Title string
	Lists []trellocms.List
}
var routes = map[string]string{
	"/": "Blog",
	"/blog": "Blog",
	"/about": "About",
}

func (this *MyHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {

	var context Context
	if route, ok := routes[req.URL.Path]; ok {
		list := this.Lists.GetByName(route)
		cards, err := list.GetCards(this.Config)
		if err != nil {
			panic(err)
		}

		context = Context{
			Cards: cards,
			Title: list.Name,
			Lists: this.Lists.Lists,
		}

		fmt.Println("route defined", context)
	} else {
		fmt.Println("route not defined")
	}

	w.Header().Add("Content Type", "text/html")
	templates := template.New("trellocms")

	templates.New("trellocms").Parse(doc)
	templates.New("footer").Parse(footer)
	templates.New("header").Parse(header)
	templates.New("navbar").Parse(navbar)

	templates.Lookup("trellocms").Execute(w, context)

}

const header = `
<!DOCTYPE html>
<html>
	<head>
		<link rel="stylesheet" href="https://maxcdn.bootstrapcdn.com/bootstrap/3.3.5/css/bootstrap.min.css"/>
		<title>{{.}}</title>
	</head>
`
const footer = `
	<script src="https://code.jquery.com/jquery-2.1.4.min.js"></script>
	<script src="https://maxcdn.bootstrapcdn.com/bootstrap/3.3.5/js/bootstrap.min.js" integrity="sha256-Sk3nkD6mLTMOF0EOpNtsIry+s1CsaqQC1rVLTAy+0yc= sha512-K1qjQ+NcF2TYO/eI3M6v8EiNYZfA95pQumfvcVrTHtwQVDG+aHRqLi/ETn2uB+1JqwYqVG3LIvdm9lj6imS/pQ==" crossorigin="anonymous"></script>
</html>
`
const doc = `
	{{template "header" .Title}}
	<body>
		<div class="container">
			<div class="row">
				<div class="col-sm-8 col-sm-offset-2">
					{{template "navbar" .Lists}}
				</div>
				<div class="col-sm-8 col-sm-offset-2">
					<h1>{{.Title}}</h1>
					<ul class="list-unstyled">
						{{range .Cards}}
							<li>
								<h2>{{.Name}}</h2>
								<p>{{.Desc}}</p>
							</li>
						{{end}}
					</ul>
				</div>
			</div>
		</div>
	</body>
	{{template "footer"}}
`

const navbar = `
<nav class="navbar navbar-default">
  <div class="container-fluid">
    <!-- Brand and toggle get grouped for better mobile display -->
    <div class="navbar-header">
      <button type="button" class="navbar-toggle collapsed" data-toggle="collapse" data-target="#bs-example-navbar-collapse-1" aria-expanded="false">
        <span class="sr-only">Toggle navigation</span>
        <span class="icon-bar"></span>
        <span class="icon-bar"></span>
        <span class="icon-bar"></span>
      </button>
      <a class="navbar-brand" href="/">Trello-CMS</a>
    </div>

    <!-- Collect the nav links, forms, and other content for toggling -->
    <div class="collapse navbar-collapse" id="bs-example-navbar-collapse-1">
      <ul class="nav navbar-nav">
				{{range .}}
					<li><a href="/{{.Slug}}">{{.Name}}</a></li>
				{{end}}
      </ul>
    </div><!-- /.navbar-collapse -->
  </div><!-- /.container-fluid -->
</nav>
`


func check(err error) {
	if err != nil {
		panic(err)
	}
}

