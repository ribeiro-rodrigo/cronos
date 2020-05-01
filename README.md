# Cronos: Dependency Injection for Go

Cronos is a reflection-based dependency injection library designed for Go projects.
Dependencies between components are represented in Cronos as constructor function 
parameters, encouraging explicit initialization rather than global variables.
Some features of Cronos were based on the Spring Framework, such as non-singleton 
dependency injection and interface-based dependency injection with the ability to specify the implementation.

## Installing

Install Cronos by running:

```shell
go get github.com/ribeiro-rodrigo/cronos
```
and ensuring that `$GOPATH/bin` is added to your `$PATH`.

## Project status

Cronps is currently in *beta*. During the beta period, we encourage you to use Cronos and provide feedback. 
We will focus on improving and evolving the library as the needs of the community.

## Usage

The Cronos dependency injection container can be created as follows.

```go
container := cronos.New()
```
Cronos is based on building dependencies through constructor functions. Below is an example of creation.

```go
type Person struct {
  name string 
}

func NewPerson()Person{
  return Person{"Bob"}
}
```
Constructor functions must be registered via the Register method of the container.

```go
container.Register(NewPerson)
```
The container can be initialized using the Init method.

```go
container.Init(func(person Person){
  fmt.Println(person.name)
})
```

