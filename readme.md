# 🚂 Locomotive
Locomotive is a tiny library build on top of GO stdlib and Gorilla.Mux to simplify web project creation. Things that you get using locomotive:
* Simple dependency injector.
* Glue code for controllers and views.

It's recommended to use locomotive-cli to automate and boost 🚀 project creation and structure: https://github.com/deltegui/locomotive-cli.

## Dependency Injection API
The way that the dependency injector works it's using builders. A builder is a function that takes as parameters the types you need in order to create an instance of your struct and returns your brand new struct. For instance let's use a builder to create a Human:

    func  createHuman() Human { 
	    return Human{
		    Name:  "Rick",
		    Age:  70,
	    }
	}
A Human struct looks like this:

    type Human struct {
	    Name string
	    Age int
	}
You can register your builder by adding it to the injector:

    injector.Add(createHuman)
And then you can check that your builder has been registered with this method:

    injector.ShowAvailableBuilders()
That line should print something like this:

    Builder for type: main.Human
And now you can build your type from everywhere using injector.Get and passing a empty Human struct or using reflect.Type

    injector.Get(Human{})
    injector.GetByType(reflect.TypeOf(Human{}))
 Well, for now things are simple. Let's try to complicate it. Now we have this:
 

	type Human struct {
	    Name string
	    Age int
	}
	type Morty struct {
		Me Human
		Pet string
	}
And then we have these builders:

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
So, to create a Morty the injector must provide a string and an instance for a Human struct. Well, simply add all builders to the injector and let the injector do all the work:

    injector.Add(createHuman)
    injector.Add(makeMorty)
    injector.Add(createPet)
And then:

    injector.ShowAvailableBuilders()
    fmt.Println(injector.Get(Morty{}))
That should print something like this:

    Builder for type: main.Human
    Builder for type: main.Morty
    Builder for type: string
    {{Morty 14} snuffles}
The injector system it's used mainly to create Controllers that need external services to be injected to. For example here we have a builder for a real world controller:

    func NewSensorController(sensorRepo SensorRepo, validator Validator, reporter Reporter, reportTypeRepo ReportTypeRepo) SensorController {...}
**NOTE:** Be careful creating a builder that take the same type that returns. It will create a infinite call to builders and it will crash.

## Controllers
A Controller struct is any struct that implements this interface:

    type Controller interface {
	    GetMappings() []Mapping
	}

Simply we have a method called GetMappings that returns mappings to your endpoints. A mapping is a struct that looks like this:

    type Mapping struct {
	    Method HTTPMethod
	    Handler http.HandlerFunc
	    Endpoint string
	}

Sooo, let's see an example of use:

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

Then you can map that controller struct like this:

    locomotive.Map("/sensor", NewSensorController)

Notice that you map a builder for your controller. If the builder takes parameters, ensure you have added all of them to the injector.

You can use MapRoot as an alias of Map("/", ...):

    locomotive.MapRoot(NewSensorController)

After all, you can run locomotive:

    locomotive.Run("localhost:3000")

And access your endpoints like *http://localhost:3000/sensor/hello/now*