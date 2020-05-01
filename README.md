# Cronos: Dependency Injection for Go

Cronos is a reflection-based dependency injection library designed for Go projects.
Dependencies between components are represented in Cronos as constructor function 
parameters, encouraging explicit initialization rather than global variables.
Some features of Cronos were based on the Spring Framework, such as non-singleton 
dependency injection and interface-based dependency injection with the ability to specify the implementation.
