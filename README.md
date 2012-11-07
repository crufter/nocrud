noCrud
===

noCrud is a framework supporting both traditional and real time web development.  
It is inspired by REST, RPC, MVC, a couple of other buzzwords, and reggae music.  
As a basic philosophy it employs "Configuration over code."  

Instead of coding by hand or code generation, noCrud feeds on configuration data, it acts like a DSL execution engine.  
This makes interchanging site setups easy and safe, since your whole application is drived by the Options document (optdoc for short), which is basically a JSON map.

The framework itself is a collection of interfaces (see /frame/interfaces), which can be accessed by the modules (/modules).

Expect the website soon at www.nocrud.com

Getting started
---

### Install

- get get github.com/opesun/nocrud
- go install github.com/opesun/nocrud
- install MongoDb

### Run

- have a MongoDB instance running at 127.0.0.1:27017
- issue
```
sudo nocrud
```
Or, to show some example command line arguments:
```
sudo nocrud -p=6060 -db_pass="my secret pass" -db_name=admin -db
```

For descriptions of command line arguments see the package /frame/config or type issue
```
nocrud -help
```

Basics
---

### Fundamental concepts

You can make your way trough noCrud without fancy GUI navigation by understanding two of the main concepts: **resources** and **verbs** acting on them.

#### Resources

A resource is, most commonly, a collection or table in a database.  
A resource can be either a full collection, or a subset of it, or even a single element of it.  
When you start an application with a empty optdoc, the only resource defined is the *options* collection itself.

#### Verbs and modules

A verb is a method which can be called on a resource. It is an exported method of a module.  
The verbs available on a resource is defined by which modules are assigned to a given resource.  
You can think of modules as a collection of verbs.

### Routing and method dispatch

By default, the module **jsonedit** is assigned to the resource **options**.  
The module jsonedit has 7 verbs: Get, GetSingle, Edit, New, Insert, Update, Delete.  
(You can read and edit it's code in the folder /modules/jsonedit)

By typing
```
/options
```
In the address bar of the browser, you already ran the verb Get on the resource options.  
You will see the listing of the options collection. By clicking on a record, you will follow a link similar to
```
/options/UIFZ-2-Q6QK8AAAB
```

This triggers the verb GetSingle. Get and GetSingle are the only two verbs, where you dont have to explicitly specifiy the verb itself.  
When you access a full or filtered collection (more on filters later) without any verb specified, you issue a Get command.  
When you access an element of a collection specified by an Id, without any verb specified, you issue a GetSingle command.

To get the hang of it, here is this table:

URL 								| Verb issued | Effect
----------------------------------- | ----------- | -----------
/options							| Get		  | You read a list of elements.
/options/UIFZ-2-Q6QK8AAAB 			| GetSingle	  | You read the element specified by the Id.
/options/new 						| New		  | You see a form where you can input a new element of options.
/options/delete						| Delete  	  | You delete every element of the collection options.
/options/UIFZ-2-Q6QK8AAAB/delete 	| Delete	  | You delete the element specified by the Id.
/options/UIFZ-2-Q6QK8AAAB/edit		| Edit		  | You edit the element specified by the Id.
/options/insert						| Insert	  | You insert a new element to the collection options. (Needs POST data)
/options/UIFZ-2-Q6QK8AAAB/update	| Update	  | You update a specific element of the collection options. (Needs POST data)
/options/update						| Update	  | You update all elements of the collection options. (Needs POST data)

### Filters

#### The URLs

So far, the Verbs we issued acted on resources wich where one of the two extremes: the full collection, or a single element of it specified by an Id.  
There is a middle ground, and they are called Filters.

Some possible examples of filters:
```
/cars?make=renault
/customers?vip=true
```

Both of those issue the implicit verb Get.  If one wants to use an explicit verb with filters, must do:
```
/cars/delete?make=renault
/customers/delete?vip=true
```

We have just deleted all cars with the make renault, and all customers who are VIPs.

#### The Codes

The nice thing about noCrud and it's modules that they work with filters, so once we write our code, it will be **collection independent**,  
and we dont have to write different methods like ***DeleteById***, ***DeleteByName***, etc... The Verb issued is entirely decoupled from the subject (the resources) it
acts on.

Lets see the code for a hypothetical Delete method and its module.

```
package mymodule

import(
	iface "github.com/opesun/nocrud/frame/interfaces"
)

type C struct{}

func (c *C) Delete(a iface.Filter) error {
	// ... We can of course insert custom bussiness logic here.
	_, err := a.RemoveAll()
	return err
}
```

Modules
---

### A simple example module

Create the file /modules/mymodule/mymodule.go

```go
package mymodule

type C struct{}

import "fmt"

func (c *C) ActionName() {
	fmt.Println("Hello")
}
```

### Accessing the Context in a module

If you want access to the Context in your application, create an Init function.
The Init function will be called with the Context as its parameter:

```go
package mymodule

import(
	"fmt"
	iface "github.com/opesun/nocrud/interfaces"
)
type C struct {
	ctx		iface.Context
}

func (c *Context) Init(ctx iface.Context) {
	c.ctx = ctx
}

func (c *C) ActionName() {
	fmt.Println("I can access the context now.")
}
```

### Exporting your module

You can export one object from your module, by creating /frame/mods/mymodule.go with the following contents:

```go
package mod

import "github.com/opesun/nocrud/modules/mymodule"

func init() {
	mods.register("mymodule", mymodule.C{})
}
```

Exporting may get automated away later.

### Assigning your modules to resources

After you exported your module, you can assign your module to resources (nouns) in the optdoc.

```js
{
	"nouns": {
		"my_noun": {
			"composed_of": [
				"mymodule"
			]
		}
	}
}
```

Multiple modules may be assigned to the resources. You can think of modules as a collection of verbs which can act on a given resource.  
When multiple modules are assigned to a resource, and more than one of them contains a given verb, the control will be routed to the first one in the list.

### Hooks

The methods of your exported object can act as hooks.  
One can call a hook by:
```go
a := 12
b := "x"
// ctx is iface.Context
ctx.Hooks().Select("aPlace").Fire(a, b)
```

Then, if you have an exported method of mymodule:
```go
func (c *C) MethodNameNotImportant(i int, s string) {
}
```

You can glue the two together in the options document:
```js
{
	"Hooks": {
		"aPlace": [ [ "mymodule", "MethodNameNotImportant"] ]
	}
}
```

As we can see method name is not important, but, when the method name is identical to the name of the hook, we don't have to specify it:
```js
{
	"Hooks": {
		"aPlace": [ "mymodule" ]
	}
}
```

The options document.
---

The options document (optdoc for short) drives your whole application.  

An application (in most cases, a website) is able to start with an nonexistent options document, to allow bootstrapping.
The default public template will be loaded, and a default optdoc will be loaded,
where the only resource defined will be the "options" itself, so you can edit it.

Optdoc examples are coming soon.

