# Smithy Go Style Guide

- [Introduction](#introduction)
- [Guidelines](#guidelines)
  - [Pointers to Interfaces](#pointers-to-interfaces)
  - [Don't Panic](#dont-panic)

## Introduction

Styles are the conventions that govern our code. The term style is a bit of a
misnomer, since these conventions cover far more than just source file
formatting—gofmt handles that for us.

The goal of this guide is to manage this complexity by describing in detail the
Dos and Don'ts of writing Go code at Uber. These rules exist to keep the code
base manageable while still allowing engineers to use Go language features
productively.

This documents idiomatic conventions in Go code that we follow at Smithy. A lot
of these are general guidelines for Go, while others extend upon external
resources:

1. [Effective Go](https://go.dev/doc/effective_go)
2. [Go Common Mistakes](https://go.dev/wiki/CommonMistakes)
3. [Go Code Review Comments](https://go.dev/wiki/CodeReviewComments)
4. [Uber Go Style Guide](https://github.com/uber-go/guide/blob/master/style.md)

## Guidelines

### Pointers to Interfaces

You almost never need a pointer to an interface. You should be passing
interfaces as values—the underlying data can still be a pointer.

An interface is two fields:

1. A pointer to some type-specific information. You can think of this as
   "type."
2. Data pointer. If the data stored is a pointer, it’s stored directly. If
   the data stored is a value, then a pointer to the value is stored.

If you want interface methods to modify the underlying data, you must use a
pointer.

Pointer to interfaces are quite tedious to dereference.

<table>
<thead><tr><th>Bad</th><th>Good</th></tr></thead>
<tbody>
<tr><td>

```go
type (
	Shape interface {
            Area() float64
    }

	Circle struct {
            Radius float64
	}
)

func (c *Circle) Area() float64 {
    return 3.14 * c.Radius * c.Radius
}

func printArea(s *Shape) {
    fmt.Println((*s).Area()) // Dereferencing pointer to interface
}

func main() {
    c := &Circle{Radius: 5}

    var s Shape = c // Assign Circle to Shape
    printArea(&s)   // Passing pointer to interface (bad)
}
```

</td><td>

```go
type (
    Shape interface {
            Area() float64
    }

    Circle struct {
            Radius float64
    }
)

func (c *Circle) Area() float64 {
    return 3.14 * c.Radius * c.Radius
}

func printArea(s Shape) {
    fmt.Println(s.Area()) // No need to dereference
}

func main() {
    c := &Circle{Radius: 5}

    var s Shape = c // Assign Circle to Shape
    printArea(s)    // Passing interface by value
}
```

</td></tr>
</tbody></table>

### Don't Panic

Code running in production must avoid panics. Panics are a major source of
[cascading failures](https://en.wikipedia.org/wiki/Cascading_failure). 
If an error occurs, the function must return an error and
allow the caller to decide how to handle it.

<table>
<thead><tr><th>Bad</th><th>Good</th></tr></thead>
<tbody>
<tr><td>

```go
func run(args []string) {
  if len(args) == 0 {
    panic("an argument is required")
  }
  // ...
}

func main() {
  run(os.Args[1:])
}
```

</td><td>

```go
func run(args []string) error {
  if len(args) == 0 {
    return errors.New("an argument is required")
  }
  // ...
  return nil
}

func main() {
  if err := run(os.Args[1:]); err != nil {
    fmt.Fprintln(os.Stderr, err)
    os.Exit(1)
  }
}
```

</td></tr>
</tbody></table>

Panic/recover is not an error handling strategy. 
A program must panic only when
something irrecoverable happens such as a nil dereference. 
An exception to this is
program initialization: bad things at program startup that 
should abort the program may cause panic.

```go
var _statusTemplate = template.Must(template.New("name").Parse("_statusHTML"))
```

Even in tests, prefer `t.Fatal` or `t.FailNow` over panics to ensure that the
test is marked as failed.

<table>
<thead><tr><th>Bad</th><th>Good</th></tr></thead>
<tbody>
<tr><td>

```go
// func TestFoo(t *testing.T)

f, err := os.CreateTemp("", "test")
if err != nil {
  panic("failed to set up test")
}
```

</td><td>

```go
// func TestFoo(t *testing.T)

f, err := os.CreateTemp("", "test")
if err != nil {
  t.Fatal("failed to set up test")
}
```

</td></tr>
</tbody></table>

Panics should always be reported in a way that
the team is aware that a service is having such
issue. A tool like [Sentry](https://sentry.io/welcome/)
is excellent to report such extreme issues and make
sure that the team is notified to resolve the root cause.

[It's useful to see the stacktrace at the moment
of a process panicking](https://yourbasic.org/golang/recover-from-panic/). 
Logs, Traces and Metrics should be enriched with the latter.
