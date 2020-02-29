# gobuilder

## Installation
`go get -u github.com/deff7/gobuilder/cmd/gobuilder`

## Examples

### Sample usage
File:
```go
// foo.go
package foo

type Name struct {
	First string
	Last  string
}

type User struct {
	ID   int
	Name Name
	Age  int
}

type UserData struct {
	Data map[string]interface{}
}
```

Generate builders for structures that have `User` substring and skip fields with `ID` substring

`gobuilder -f foo.go -s User -fields-filter=ID`

```go
// UserBuilder is builder for type User
type UserBuilder struct {
	instance *main.User
}

// User creates new builder
func User() *UserBuilder {
	return &UserBuilder{
		instance: &main.User{},
	}
}

// Name sets field with type Name
func (b *UserBuilder) Name(v Name) *UserBuilder {
	b.instance.Name = v
	return b
}

// Age sets field with type int
func (b *UserBuilder) Age(v int) *UserBuilder {
	b.instance.Age = v
	return b
}

// P returns pointer to User instance
func (b *UserBuilder) P() *main.User {
	return b.instance
}

// V returns value of User instance
func (b *UserBuilder) V() main.User {
	return *b.instance
}

// UserDataBuilder is builder for type UserData
type UserDataBuilder struct {
	instance *main.UserData
}

// UserData creates new builder
func UserData() *UserDataBuilder {
	return &UserDataBuilder{
		instance: &main.UserData{},
	}
}

// P returns pointer to UserData instance
func (b *UserDataBuilder) P() *main.UserData {
	return b.instance
}

// V returns value of UserData instance
func (b *UserDataBuilder) V() main.UserData {
	return *b.instance
}
```

### Generate recursively for specific directory
`gobuilder -R -d ./tmp`
