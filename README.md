# Factory Go

A fixtures replacement for GO, inspired by factory_bot in Ruby with the goal to be easy to setup and use

*This package is still in early development*

## Getting started

Install the dependency

    go get -u github.com/kejiro/factory
     
Use in the project:

    import "github.com/kejiro/factory"
    
## Usage

Use with static data
```go
	type User struct {
		FirstName string
		LastName  string
		Dob       time.Time
	}

	Define(User{}, Definition{
		"FirstName": "John",
		"LastName":  "Doe",
	})
	
    user := User{}
    Build(&user, Definition{})
    fmt.Println(user)

```

or with generated data
```go
	type User struct {
		FirstName string
		LastName  string
		Dob       time.Time
	}

	Define(User{}, Definition{
		"FirstName": "John",
		"LastName": func() string {
			return "Doe"
		},
		"Dob": func() time.Time {
			dob, _ := time.Parse(time.RFC3339, "1970-01-01T00:00:00Z")
			return dob
		},
	})

	user := User{}
	Build(&user, Definition{})
	fmt.Println(user)
```

it is also possible to override defined values on a per build request 
```go
	type User struct {
		FirstName string
		LastName  string
		Dob       time.Time
	}

	Define(User{}, Definition{
		"FirstName": "John",
		"LastName": func() string {
			return "Doe"
		},
		"Dob": func() time.Time {
			dob, _ := time.Parse(time.RFC3339, "1970-01-01T00:00:00Z")
			return dob
		},
	})

	user := User{}
	Build(&user, Definition{"FirstName": "Jack"})
	fmt.Println(user)
```

