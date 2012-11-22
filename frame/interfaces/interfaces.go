// The interfaces package aims to describe all interfaces available to modules.
package interfaces

import (
	"reflect"
)

// This contains everything modules can access.
type Context interface {
	Conducting() Conducting
	FileSys() FileSys
	User() User
	Client() Client
	Db() Db
	Channels() Channels
	ViewContext() ViewContext
	NonPortable() NonPortable
	Display() Display
	Options() Options
}

type Conducting interface {
	Hooks() Hooks
	Events() Events
}

// Hooks are entry points for other modules.
type Hooks interface {
	Select(string) Hook
	Module(string) Module
}

type Subscriber interface {
	Name() string
	Method() string
}

// The Fire method is simply runs all modules subscribed to the given hook.
// With Iterate, one can provide a stop function as a first parameter, which will recieve the output of all hooks.
// If the stop function returns true, the execution of subscribers stop.
type Hook interface {
	Subscribers() []Subscriber
	HasSubscribers() bool
	SubscriberCount() int
	Fire(params ...interface{})
	Iterate(stopfunc interface{}, params ...interface{})
}

// With events one can trigger events and send messages cross-process, cross-machine or cross-network.
type Events interface {
	Select(string) Event
}

type Event interface {
	Publish([]byte) error
	Subscribe() error
	Unsubscribe() error
	Read() ([]byte, error)
}

type Method interface {
	Call(interface{}, ...interface{}) error // First param is the return reciever function, others are input arguments for the Method.
	Matches(interface{}) bool               // Returns true if the Method matches the signature of the supplied function.
	InputTypes() []reflect.Type
	OutputTypes() []reflect.Type
}

// An instance is an empty instance of a given type.
// To convert an instance of any type to an Instance, see /frame/mod.ToInstance(interface{})
type Instance interface {
	HasMethod(string) bool // Returns true if the Instance has a method with the name supplied.
	MethodNames() []string // Returns all public method names of Instance.
	Method(string) Method
}

// Module is a package/module selector.
// Used to get around the lack of dynamic code loading.
type Module interface {
	Instance() Instance
	Exists() bool
}

// Speaker is used to decipher the request path - it can tell if a string is a registered noun,
// and if that noun has a given verb defined on it.
type Speaker interface {
	IsNoun(string) bool
	NounHasVerb(string, string) bool
}

type QueryMod interface {
	Skip(int)
	Limit(int)
	Sort(...string)
}

// Set represents a collection of data, mainly data coming from a database.
// See interface Filter for further information.
type Set interface {
	QueryMod
	Count(map[string]interface{}) (int, error)
	FindOne(map[string]interface{}) (map[string]interface{}, error)
	Find(map[string]interface{}) ([]map[string]interface{}, error)
	Insert(map[string]interface{}) error
	// InsertAll([]map[string]interface{}) errors
	Update(map[string]interface{}, map[string]interface{}) error
	UpdateAll(map[string]interface{}, map[string]interface{}) (int, error)
	Remove(map[string]interface{}) error
	RemoveAll(map[string]interface{}) (int, error)
	Name() string
}

// A filter is a collection of data, where you only have access to a subset of the whole set.
// See interface Set for further information.
type Filter interface {
	// Returns the Ids of the documents in the filtered set.
	Ids() ([]Id, error)
	// AddQuery will further filter the already filtered set. A filter should never refer to a larger subset of the set after the
	// AddQuery operation than before it.
	AddQuery(map[string]interface{}) Filter
	// Cloning a filter allows you to "branch" your filters.
	// Eg: (pseudocode)
	// fa := Filter{a:2}
	// fb := fa.Clone().AddQuery({b:3}) (equals: Filter{a:2,b:3})
	// fc := fa.Clone().AddQuery({d:4})	(equals: Filter{a:2,d:4})
	Clone() Filter
	// Reducing filters is a way to query foreign key relationships.
	Reduce(...Filter) (Filter, error)
	// Subject returns you the collectionname the filter refers to.
	Subject() string
	AddParents(string, []Id)
	Modifiers() Modifiers
	Count() (int, error)
	Iterate(func(Document) error) error
	// --
	FindOne() (map[string]interface{}, error)
	Find() ([]map[string]interface{}, error)
	SelectOne() (Document, error)
	Insert(map[string]interface{}) error
	// InsertAll([]map[string]interface{}) errors
	Update(map[string]interface{}) error
	UpdateAll(map[string]interface{}) (int, error)
	Remove() error
	RemoveAll() (int, error)
}

// Id represents a document Id.
type Id interface {
	String() string
}

// Document represents a database document.
type Document interface {
	Data() map[string]interface{}
	Id() Id
	Update(map[string]interface{}) error
	Remove() error
}

// Unused atm,it's just a draft.
type DocumentPointer interface {
	Get() (Document, error)
	Update(map[string]interface{}) error
	Remove() error
	Id() Id
}

// Modifiers are used to modify the querying of a collection, see interfaces Fiter and Set.
type Modifiers interface {
	Sort() []string
	Limit() int
	Skip() int
}

// Channels are used as a last resort of passing data trough the system.
// Caution, nothing like the channels in Go.
type Channels interface {
	Select(string) Channel
}

type Channel interface {
	Send(interface{})
	Get() []interface{}
	GetFirst() interface{}
	// Returns the Xth element in the channel.
	GetX(int) interface{}
	HasData() bool
	// Returns true if the channel has an Xth element in it.
	HasX(int) bool
}

// Database session.
// Ugly circularity between Session and Db.
type Session interface {
	// Select database.
	Db(string) (Db, error)
}

// Db represents a database.
// The NewFilter gives you filtered access to a certain collection.
// (Except if you have no rights at all, then gives error)
type Db interface {
	NewFilter(string, map[string]interface{}) (Filter, error)
	ToId(string) (Id, error)
	NewId() Id
	Session() (Session, error)
	// Resolve(map[string]interface{}, ...string) error
}

// User represents the user interacting with the application.
type User interface {
	Document
	Level() int
	Languages() []string
}

// The client allows you to store and retrieve data at the client.
type Client interface {
	StoreEncrypted(string, interface{}) error
	Store(string, interface{}) error
	Get(string) (interface{}, error)
	GetDecrypted(string) (interface{}, error)
	Unstore(string) error
	Languages() []string
}

// ViewContext contains the data the views has access to.
type ViewContext interface {
	Publish(string, interface{}) ViewContext
	Get() map[string]interface{}
}

// FileInfo is a generalized interface to obtain information about a file.
type FileInfo interface {
	Name() string
	IsDir() bool
	File() File
	Directory() Directory
}

// A directory will never allow you to access a directory or file residing in any of its parent directory.
// For accessing directories see the SelectPlace method of interface FileSys.
type Directory interface {
	File(string) File
	Directory(...string) Directory
	Remove() error
	Exists() (bool, error)
	List() ([]FileInfo, error)
	Create() error
	Rename(string) error
}

// We will have to adjust this to allow handling large files.
type File interface {
	Create() error
	Exists() (bool, error)
	Write([]byte) error
	Read() ([]byte, error)
	Remove() error
	Rename(string) error
	Name() string
}

// Readable file provides a subset of the functionality provided by the interace File: it can be read only.
// (Used at uploaded temporary files)
type ReadableFile interface {
	Read() ([]byte, error)
	Name() string
}

// A way of accessing uploaded files without getting too involved in implementation details.
type Temporaries interface {
	Select(string) []ReadableFile
	Exists(string) bool
	Keys() []string
}

// FileSys is everything you can with your filesystem.
type FileSys interface {
	// With the help of the SelectPlace one can (only) access
	// predefined places in the filesystem, eg: "modules" is mapped to /modules,
	// "template" is mapped to the current template in use, etc...
	SelectPlace(string) (Directory, error)
	Temporaries() Temporaries
}

// Maybe this could be called request.
type NonPortable interface {
	Resource() string
	Params() map[string]interface{}
	Redirect(string) // Not sure about this.
	ComingFrom() string
	View() bool
	RawParams() string
}

type Writer interface {
	Write([]byte) (int, error)
}

// Interface Display sends visible data to the client.
type Display interface {
	Writer() Writer
	Write([]byte) error
	Type(string) Display
}

// Helper interface to handle multiply nested data structures (eg. JSON).
// Should support dot notation ("subscribers.Johnny.age").
type NestedData interface {
	Get(...string) (interface{}, bool)
	GetI(...string) (int64, bool)
	GetF(...string) (float64, bool)
	GetM(...string) (map[string]interface{}, bool)
	GetS(...string) ([]interface{}, bool)
	GetStr(...string) (string, bool)
	Exists(...string) bool
	All() interface{}
}

// Document() gives you back the option document,
// Modifiers() returns parameters specified by the client (eg. json=true).
type Options interface {
	Document() NestedData
	Modifiers() NestedData
}
