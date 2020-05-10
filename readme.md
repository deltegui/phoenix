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

## Controllers
A Controller struct is anything that implements this interface:

```go
type Controller interface {
	GetMappings() []Mapping
}
```

Simply we have a method called GetMappings that returns mappings to your endpoints. A mapping is a struct that looks like this:

```go
type Mapping struct {
	Method HTTPMethod
	Handler http.HandlerFunc
	Endpoint string
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
func (controller SensorController) GetMappings() []locomotive.Mapping {
	return []locomotive.Mapping{
		{Method: locomotive.Post, Handler: controller.SaveSensor, Endpoint: ""},
		{Method: locomotive.Get, Handler: controller.GetSensorByName, Endpoint: "/{name}"},
		{Method: locomotive.Delete, Handler: controller.DeleteSensorByName, Endpoint: "/{name}"},
		{Method: locomotive.Post, Handler: controller.UpdateSensor, Endpoint: "/update"},
		{Method: locomotive.Get, Handler: controller.SensorNow, Endpoint: "/{name}/now"},
	}
}
```
It's implementing the controller interface and returning mappings in it's GetMappings method. Then you can map that controller using locomotive's Map:

```go
locomotive.Map("/sensor", NewSensorController)
```

Notice that you map a builder for your controller. If the builder takes parameters, ensure you have added builders for all of them to the injector.

You can use MapRoot as an alias of Map("/", ...):

```go
locomotive.MapRoot(NewSensorController)
```

After all, you can run locomotive:

```go
locomotive.Run("localhost:3000")
```

And access your endpoints like *http://localhost:3000/sensor/hello/now*

## Presenters
So, now we have a way to wire all up and we have controllers too. How do we return something from our system, like HTML or JSON? With presenters.
Presenters are an abstraction over one way to present data to users. Here we have the presenter interface:

```go
type Presenter interface {
	Present(data interface{})
	PresentError(data error)
}
```

As you can see, a presenter can 'present' an error (passing the error) or data (using a struct as ModelView) to the user.
Locomotive have two implementations of presenters: JSONPresenter and HTMLPresenter.

### JSONPresenter
Simply it renders your modelview as JSON. Here you have an example of use from a UseCase that deletes a router:

```go
func NewDeleteRouterCase(routerRepo RouterRepository, validator Validator) UseCase {
	return func(presenter Presenter, req UseCaseRequest) {
		if err := validator.Validate(req); err != nil {
			presenter.PresentError(MalformedRequest)
			return
		}
		deleteReq := req.(DeleteRouterRequest)
		if !routerRepo.ExistsWithName(deleteReq.RouterName) {
			presenter.PresentError(RouterNotFound)
			return
		}
		routerRepo.DeleteWithName(deleteReq.RouterName)
		presenter.PresentInformation(EmptyRequest)
	}
}
```

### HTMLPresenter
The HTMLPresenter takes your ModelView and renders an HTML using go's templates.

Firstly, you will need a place to put your templates. Well, go's templates in locomotive will search for your html templates here: "./templates/\*/\*.html". That means you must create in your project root a folder named "templates". Inside that folder it will expect a bunch of folders (you can have as many folders as you want, and with names you like), that must contain your template. So, if you want to render "userindex.html", you must have a folder inside templates which can have any name (let's take "user"), and inside it, your template "userindex.html" like this:

\<project root\>/templates/user/userindex.html

Be careful naming your templates. If you create two templates with the same name in two distinct folders it will always render the first it find. That's because it'll look for templates inside all subfolders.

But how can I tell to the presenter that I want to render userindex.html? Presenter's interfaces only takes a ModelView in it's present method. It will automatically take the name of your method and downcase it. So, if the name of the method that calls the presenter is "UserIndex", HTMLPresenter will search for a template named "userindex.html"

What happens if it found no template? It will try with the name of the callee of the method that called the presenter. If no template was found, It will try once again with the previous callee. Finally, if no template was found, it will fail.

Why three times? To adapt the presenter to the project structure generated by locomotive-cli.