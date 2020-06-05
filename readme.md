# 🐦 Phoenix
Phoenix is a tiny library build on top of GO stdlib, Gorilla.Mux, Gorilla.Sessions and autocert to simplify web project creation. Things that you get using phoenix:
* Simple dependency injector.
* Glue code for controllers.
* Ready to use HTML and JSON renderers.
* Available HTTPS connection using cert.pem and key.pem or using autocert.
* Sessions.
* Middlewares.
* CSRF middleware.
* Startup ASCII logo support.
* Gracefully server stop.

It's recommended to use phoenix-cli to automate and boost 🚀 project creation and structure: https://github.com/deltegui/phoenix-cli.

## When to use phoenix?

When you feel that you are repeating "glue code" everytime you want to create a server in go. Otherwise use go's stdlib.

## Little Example

```go
package main

import (
	"fmt"
	"net/http"

	"github.com/deltegui/phoenix"
	"github.com/gorilla/mux"
)

func main() {
	app := phoenix.NewApp()
	app.Get("/hello/{name}", func() http.HandlerFunc {
		return func(w http.ResponseWriter, req *http.Request) {
			vars := mux.Vars(req)
			fmt.Fprintf(w, "Hello %s", vars["name"])
		}
	})
	app.Run("localhost:8080")
}
```

## Application
Firstly, you can create an application.

```go
app := phoenix.NewApp()

```

Then you can access:

* Dependency Injector container throught ```app.Injector```
* App's config calling ```app.Configure()```
* Start server calling ```app.Run("localhost:8080")```

## Configuration
You can configure phoenix to fit your needs.
phoenix by default have disabled:

* Logo File.
* Static server.
* Sessions.
* HTTPS.

And the project name and version will be "phoenix" and "v0.1.0" by default.

If you want to configure these options you can get a PhoenixConfig struct, then you can change your project name and version and enable the other features:

```go
app.Configure().
	SetProjectVersion("your project name", "your project version").
	EnableLogoFile(). // It will search a file named "logo" inside your project root.
	EnableStaticServer(). // It will serve static files that are inside /static.
	EnableSessions(). // It will enable sessions. sessions.Store will be available to inject into your controllers.
	SetStopHandler(func() {...}). // Set the function to be called when server is stopping.
	UseHTTPS("cert", "key"). // Enable HTTPS using cert and key.
	UseAutoHTTPS("example.com", "www.example.com", ...). // Alternative to UseHTTPS that uses let's encrypt with autocert

```

## Dependency Injection API
The way that the dependency injection works it's using builders. A builder is a function that takes as parameters the types you need and returns your brand new struct. For instance let's use a builder to create a Human:

```go
func  createHuman() Human {
	return Human{
		Name:  "Rick",
		Age:  70,
	}
}
```
A Human struct looks like this:

```go
type Human struct {
	Name string
	Age int
}
```

You can register your builder by adding it to the injector:

```go
app.Injector.Add(createHuman)
```

And then you can check that your builder has been registered with this method:

```go
app.Injector.ShowAvailableBuilders()
```

That line should print something like this:

```
Builder for type: main.Human
```

And now you can build your type from everywhere using app.Injector.Get and passing an empty Human struct or using reflect.Type

```go
app.Injector.Get(Human{})
app.Injector.GetByType(reflect.TypeOf(Human{}))
```

 Well, for now things are simple. Let's try to complicate it. Now we have this:

```go
type Human struct {
	Name string
	Age int
}

type Morty struct {
	Human
	Pet string
}
```
And then we have these builders:

```go
func  makeMorty(human Human, pet string) Morty {
	return Morty{human, pet}
}

func  createHuman() Human {
	return Human{
		Name:  "Morty",
		Age:  14,
	}
}

func  createPet()  string  {
	return  "snuffles"
}
```

So, to create a Morty the injector must provide a string and an instance for a Human struct. Well, simply add all builders to the injector and let the injector do all the work:

```go
app.Injector.Add(createHuman)
app.Injector.Add(makeMorty)
app.Injector.Add(createPet)
```

And then:

```go
app.Injector.ShowAvailableBuilders()
fmt.Println(app.Injector.Get(Morty{}))
```

That should print something like this:

```
Builder for type: main.Human
Builder for type: main.Morty
Builder for type: string
{{Morty 14} snuffles}
```

The injector system it's used mainly to create Controllers that need external services to be injected to. For example here we have a builder for a real world controller:
```go
func NewSensorController(sensorRepo SensorRepo, validator Validator, reporter Reporter, reportTypeRepo ReportTypeRepo) SensorController {...}
```

**NOTE:** Be careful creating a builder that take the same type that returns. It will create a infinite builder call and it will crash.

## Middlewares

A middleware is simply a functions that takes a http.HandlerFunc and returns a http.HandlerFunc. For example, take a look to the log middleware:

```go
func logMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(writer http.ResponseWriter, req *http.Request) {
		log.Printf("Request [%s] %s from %s\n", req.Method, req.RequestURI, req.RemoteAddr)
		next.ServeHTTP(writer, req)
	}
}
```

## Handler Builder
A handler builder is a wrapper over a http.HanlderFunc that provides dependencies to your handler. For example:

```go
func Hello(name string) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		r := phoenix.JSONRenderer{w}
		r.Render(struct{ Name string }{name})
	}
}
```

As you can see, it's just a wrapper for your simple handler. With that wrapper you can tell to the phoenix's DI system to inject anything you want.
Now you can simply map that builder:

```go
app.Map(phoenix.Mapping{phoenix.Get, "/hello", Hello}, [...Middlewares])
```

That use a Mapping struct, that looks like this:

```go
type Mapping struct {
	Method   HTTPMethod
	Endpoint string
	Builder  injector.Builder
}
```

Or you can map in http GET and / endpoint using MapRoot (an alias of app.Map(phoenix.Mapping{phoenix.Get, "", ...})):

```go
app.MapRoot(Hello)
```

Also, if you are mapping using a HTTP GET, POST, DELETE or PUT you can change app.Map(...) to:

```go
app.Get("/hello", Hello, [...Middlewares])
app.Post("/hello", Hello, [...Middlewares])
app.Delete("/hello", Hello, [...Middlewares])
app.Put("/hello", Hello, [...Middlewares])
```

Then you can run your server

```go
phoenix.Run("localhost:3000")
```

If you want a bunch of handlers that have a part of the endpoint in common you can use MapGroup:

```go
phoenix.MapGroup("/greetings", func(m phoenix.Mapper) {
	m.MapController("/intendente", NewErrorController)
	m.MapRoot(Hello)
	m.Get("/diego", HelloDiego)
})
```

That will create the endpoints "/greetings", "/greetings/intendente" and "/greetings/diego"


**NOTE:** Map, MapRoot, Get, Delete, Post and Put always adds a trailing slash to your route. So the route "/enpoint" will create two mappings, one is "/endpoint" and the other "/endpoint/"

## Renderers
So, now we have a way to wire all up and we have controllers too. How do we return something from our system, like HTML or JSON? With renderers.

### JSONRenderer
Simply it renders your modelview as JSON. Here you have an example of use from a controller:

```go
func (controller Controller) JSON(w http.ResponseWriter, req *http.Request) {
	phoenix.NewJSONRenderer(w).Render(struct{Name string}{"phoenix"})
}
```

JSONRenderer have these methods:

```go
Render(interface{})
RenderError(error)
```

### HTMLRenderer
The HTMLRenderer takes your ModelView and renders an HTML using go's templates. You can use it like this:

```go
func (controller Controller) JSON(w http.ResponseWriter, req *http.Request) {
	phoenix.NewHTMLRenderer(w).RenderData("hello.html", struct{Name string}{"phoenix"})
}
```

HTMLRenderer have these methods:

```go
RenderData(view string, data interface{})
Render(view string)
RenderError(error) // Equals to RenderData("error.html", error)
```

Firstly, you will need a place to put your templates. Well, go's templates in phoenix will search for your html templates here: "./templates/\*/\*.html". That means you must create in your project root a folder named "templates". Inside that folder it will expect a bunch of folders (you can have as many folders as you want, and with names you like), that must contain your template. So, if you want to render "userindex.html", you must have a folder inside templates which can have any name (let's take "user"), and inside it, your template "userindex.html" like this:

\<project root\>/templates/user/userindex.html

Be careful naming your templates. If you create two templates with the same name in two distinct folders it will always render the first it finds. That's because it'll look for templates inside all subfolders.

## CSRF

To use the csrf middleware simply add it to the route

```go
app.Get("/hello", Hello, phoenix.NewCSRFMiddleware())
```

Then you can render your HTML view like this:

```go
...
phoenix.NewHTMLRenderer.RenderData("userindex.html", map[string]interface{}{
	csrf.TemplateTag: csrf.TemplateField(req),
})
...
```

And inside your view you can render the CSRF token using ```{{.csrfField}}```

Everytime you create a CSRF middleware a new secret key is created.

## Sessions

To use sessions you have to enable it:

```go
app.Configure().EnableSessions()

```

Then you can add sessions.Store (gorilla.Sessions) as a parameter inside your controller and it will be automatically injected:

```go
func Hello(session sessions.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		...
	}
}
```

The cookie key will be generated and stored inside a file called ```cookie.key``` the first time you enable sessions. Later the phoenix will try to read the key from that file. It's recommended not to commit ```cookie.key```
