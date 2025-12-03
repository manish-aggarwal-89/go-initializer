**{{SERVICE_TITLE}}**

# Table Of Content

- [Table Of Content](#table-of-content)
    - [Dependency Injection](#dependency-injection)
    - [init](#init)

## Dependency Injection

We are using dig uber library for a reflection based dependency injection toolkit for Go. On this package we need to
create new dig (Container) on init.go

## init

Use for register our DI (Dependency Injection) to our init. So we can use the dependenncy we want

```go
package sample

var Container = dig.New()

func init() {
  // - config
  if err := Container.Provide(config.New); err != nil {
    panic(fmt.Sprintf("failed to provide config %s", err))
  }

  // register layer you need to add on DI on here

}
```
