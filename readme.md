# 🚂 Locomotive
Locomotive is a tiny library build on top of GO stdlib and Gorilla.Mux to simplify web project creation. Things that you get using locomotive:
* Simple dependency injector.
* Glue code for controllers and views.

It's recommended to use locomotive-cli to automate and boost 🚀 project creation and structure: https://github.com/deltegui/locomotive-cli.

## Configuration
You can configure Locomotive to fit your needs.
Locomotive by default have disabled:

* HTML template. You can't use HTMLPresenter.
* Logo File.
* Static server.

And the project name and version will be "Locomotive" and "v0.1.0" by default.

If you want to configure these options you can get a LocomotiveConfig struct, then you can change your project name and version and enable the other features:

```go
locomotive.
	Configure().
	SetProjectVersion("your project name", "your project version").
	EnableLogoFile(). // It will search a file named "logo" inside your project root.
	EnableStaticServer(). // It will serve static files that are inside /static.
	EnableTemplates() // It will enable templates.
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
injector.Add(createHuman)
```

And then you can check that your builder has been registered with this method:

```go
injector.ShowAvailableBuilders()
```

That line should print something like this:

```
Builder for type: main.Human
```

And now you can build your type from everywhere using injector.Get and passing an empty Human struct or using reflect.Type

```go
injector.Get(Human{})
injector.GetByType(reflect.TypeOf(Human{}))
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
injector.Add(createHuman)
injector.Add(makeMorty)
injector.Add(createPet)
```

And then:

```go
injector.ShowAvailableBuilders()
fmt.Println(injector.Get(Morty{}))
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

## Handler Builder
A handler builder is a wrapper over a http.HanlderFunc that provides dependencies to your handler. For example:

```go
func Hello(name string) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		r := locomotive.JSONRenderer{w}
		r.Render(struct{ Name string }{name})
	}
}
```

As you can see, it's just a wrapper for your simple handler. With that wrapper you can tell to the locomotive's DI system to inject anything you want.
Now you can simply map that builder:

```go
locomotive.Map(locomotive.Mapping{locomotive.Get, "/hello", Hello})
```

That use a Mapping struct, that looks like this:

```go
type Mapping struct {
	Method   HTTPMethod
	Endpoint string
	Builder  injector.Builder
}
```

Or you can map in http GET and / endpoint using MapRoot (an alias of locomotive.Map(locomotive.Mapping{locomotive.Get, "", ...})):

```go
locomotive.MapRoot(Hello)
```

Then you can run your server

```go
locomotive.Run("localhost:3000")
```

If you want a bunch of handlers that have a part of the endpoint in common you can use MapGroup:

```go
locomotive.MapGroup("/greetings", func(m locomotive.Mapper) {
	m.MapController("/intendente", NewErrorController)
	m.MapRoot(Hello)
	m.Map(locomotive.Mapping{locomotive.Get, "/diego", HelloDiego})
})
```

That will create the endpoints "/greetings", "/greetings/intendente" and "/greetings/diego"

## Controllers
A Controller struct is anything that implements this interface:

```go
type Controller interface {
	GetMappings() []CMapping
}
```

Simply we have a method called GetMappings that returns mappings to your endpoints. A mapping is a struct that looks like this:

```go
type CMapping struct {
	Method   HTTPMethod
	Endpoint string
	Handler  http.HandlerFunc
}
```

So, let's see an example of use. Let's see a example controller that have a bunch of endpoints:

```go
type SensorController struct {...}

func NewSensorController(...) SensorController {...}

func (controller SensorController) SaveSensor(w http.ResponseWriter, req *http.Request) {...}
func (controller SensorController) GetSensorByName(w http.ResponseWriter, req *http.Request) {...}
func (controller SensorController) DeleteSensorByName(w http.ResponseWriter, req *http.Request) {...}
func (controller SensorController) UpdateSensor(w http.ResponseWriter, req *http.Request) {...}
func (controller SensorController) SensorNow(w http.ResponseWriter, req *http.Request) {...}
func (controller SensorController) GetMappings() []locomotive.CMapping {
	return []locomotive.CMapping{
		{Method: locomotive.Post, Handler: controller.SaveSensor, Endpoint: ""},
		{Method: locomotive.Get, Handler: controller.GetSensorByName, Endpoint: "/{name}"},
		{Method: locomotive.Delete, Handler: controller.DeleteSensorByName, Endpoint: "/{name}"},
		{Method: locomotive.Post, Handler: controller.UpdateSensor, Endpoint: "/update"},
		{Method: locomotive.Get, Handler: controller.SensorNow, Endpoint: "/{name}/now"},
	}
}
```
It's implementing the controller interface and returning mappings in it's GetMappings method. Then you can map that controller using locomotive's MapController:

```go
locomotive.MapController("/sensor", NewSensorController)
```

Notice that you map a builder for your controller. If the builder takes parameters, ensure you have added builders for all of them to the injector.

You can use MapRootController as an alias of Map("/", ...):

```go
locomotive.MapRootController(NewSensorController)
```

After all, you can run locomotive:

```go
locomotive.Run("localhost:3000")
```

And access your endpoints like *http://localhost:3000/sensor/hello/now*

**NOTE:** Map, MapRoot, MapController and MapRootController always adds a trailing slash to your route. So the route "/enpoint" will create two mappings, one is
"/endpoint" and the other "/endpoint/"

## Renderers
So, now we have a way to wire all up and we have controllers too. How do we return something from our system, like HTML or JSON? With renderers.
Renderers are an abstraction over one way to render data to users. Here we have the Renderer interface:

```go
type Renderer interface {
	Render(data interface{})
	RenderWithMeta(data interface{}, metadata RenderMetadata)
	RenderError(data error)
}
```

As you can see, a rednerer can render an error (passing the error), data (using a struct as ModelView) to the user or data with some metadata used to determine how to render your data. RenderMetadata looks like this:


```go
type RenderMetadata struct {
	ViewName string
}
```

Locomotive have two implementations of renderers: JSONRenderer and HTMLRenderer.

### JSONRenderer
Simply it renders your modelview as JSON. Here you have an example of use from a controller:

```go
func (controller Controller) JSON(w http.ResponseWriter, req *http.Request) {
	renderer := locomotive.JSONRenderer{w}
	renderer.Render(struct{Name string}{"Locomotive"})
}
```

### HTMLRenderer
The HTMLRenderer takes your ModelView and renders an HTML using go's templates.

Firstly, you will need a place to put your templates. Well, go's templates in locomotive will search for your html templates here: "./templates/\*/\*.html". That means you must create in your project root a folder named "templates". Inside that folder it will expect a bunch of folders (you can have as many folders as you want, and with names you like), that must contain your template. So, if you want to render "userindex.html", you must have a folder inside templates which can have any name (let's take "user"), and inside it, your template "userindex.html" like this:

\<project root\>/templates/user/userindex.html

Be careful naming your templates. If you create two templates with the same name in two distinct folders it will always render the first it finds. That's because it'll look for templates inside all subfolders.

But how can I tell to the HTMLRenderer that I want to render userindex.html? There's two ways:
* Calling Renderer.Render it'll automatically take the name of your method and downcase it. So, if the name of the method that calls the renderer is "UserIndex", HTMLRenderer will search for a template named "userindex.html". What happens if it found no template? It will try with the name of the callee of the method that called the renderer. If no template was found, it will try once again with the previous callee. Finally, if no template was found, it will fail. Why three times? To adapt the renderer to the project structure generated by locomotive-cli.

* Calling Renderer.RenderWithMeta and passing a RenderMetadata with the view's name.
