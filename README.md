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

