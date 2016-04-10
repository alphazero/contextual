![image](./resources/goerror.png)

####stat
    star date         04 08 16

    api               stable 
    documentation     inlined godoc

##about

Package goerror provides for creation of semantic error types.

Using the stdlib error package, trying to discern the error type returned by a function, we can either (a) compare to a well-known error reference (for example, io.EOF), or (b) we have to parse the actual error message to distiguish between error flavors. The former approach (alal io.EOF) is fine for generic errors but obviously won't allow for call-site specific information in the error. And the latter approach is ad-hoc and fragile. (See motivating case below for elaboration.)


## usage

#### defining an error 'type'

    package example
    
    import (
    	"goerror"
    )
    
    // package errors
    var (
        FooError = goerror.Define("Foo error")
        BarError = goerror.Define("Bar error")
        FubarError = goerror.Define("Fubar error")
    )

#### using an error 'type'

    package something
    
    // note that we don't need to import goerror here.
    import (
        "example"
    )

    func Doit(…) error {
    
        // creating an instance of example.Foo error
        if whynot() {
            return example.Foo("whynot")
        }
    
        // keep in mind details are optional 
        if because() {
            return example.Foo()
        }
        
        // and of course we could be returning another 
        // error flavor
        if allFubar() {
            return example.Fubar("blame it on murphy")
        }
    }

#### handling the error

With `goerror` we can distinguish between error types and also have call-site specific error details.
   
    package elsewhere
    
    import (
       "example"
       "goerror"
       "something"
    )
    
    if e := something.Doit(…); e != nil {
       switch typ := goerror.TypeOf(e); {
       case typ.Is(example.Foo):
       case typ.Is(example.Bar):
       case typ.Is(example.FooBar):
       default: 
          /* generic error */
       }
    }

Of course, you can also treat them like ordinary errors.


## motivating case

For example, let's say we create an error for asserting input arg correctness, something like an `AssertionError`, or `IllegalArgument`. We want informative (read: useful) error messages, and we also want a consistent way of distinguishing between error types. 

#### stdlib approach (I): use a global reference and pass that around

This approach allows for distinguishing between *generic* error types at the call-site, comparing references. But can you tell which of the arguments here is the one that causes the error?

###### define the error

    // use a generic error
    var IllegalArgument = errors.New("illegal argument error")

###### using the error
   
    func foo(arg0 T, arg1 T2, …) error {
        if invalid(arg0) {
           return IllegalArgument
        }
        …
        if invalid(argn) {
           return IllegalArgument
        }
        …
        
        // and don't forget; may return other kinds of error here
        e := bar(…)
        if e != nil {
            return e
        }
        … 
    }
    
###### handling the error
Here we can distinguish between generic error flavors (by comparing references) but 
   
    func foo(arg0 T, arg1 T2, …) error {
    e := foo(a, b, c, … )
    if e != nil {
        if e == IllegalArgument {
           /* which argument? */
        }
    }
    
#### stdlib approach (II): parse the error message

###### define the error

    // define the error prefix
    const IllegalArgument = "illegal argument"

###### using the error
    func foo(arg0 T, arg1 T2, …) error {
        if invalid(arg0) {
           return fmt.Errorf("%s: %s", IllegalArgument, "arg0")
        }
        …
        if invalid(argn) {
           return fmt.Errorf("%s: %s", IllegalArgument, "argn")
        }
        …
    }
    
###### handling the error

Here we can no longer distinguish between generic error flavors (by comparing references) but can get some additional information.
   
    func foo(arg0 T, arg1 T2, …) error {
    e := foo(a, b, c, … )
    if e != nil {
        errmsg := e.Error()
        // now parse the error message 
        ...
    }

