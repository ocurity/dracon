# New Go SDK Design

This document highlights the work to be done to create a new Go SDK to write Dracon components.

**Table of contents:**

* [Objectives](#objectives)
* [Callouts](#callouts)
* [Proposal](#proposal)
  * [Structure and contributions guidelines](#proposal-structure)
  * [Component API](#proposal-component-api)
    * [First pass - simple components](#proposal-component-api-first-pass)
      * [Creating a new component](#proposal-component-api-first-pass-create)
        * [Running a component](#proposal-component-api-first-pass-running)
        * [Defaults](#proposal-component-api-first-pass-defaults)
        * [Example](#proposal-component-api-first-pass-example)
    * [Second pass - environment aware components](#proposal-component-api-second-pass)
      * [Why the first pass is not enough](#proposal-component-api-second-pass-why)
      * [A new environment aware runner](#proposal-component-api-second-pass-env-aware-runner)
      * [Creating a new component](#proposal-component-api-second-pass-create-new-component)
        * [Example](#proposal-component-api-second-pass-example)
        * [Wrapping it up](#proposal-component-api-second-pass-wrap)
  * [Components catalog](#proposal-components-catalog)
  * [Default readers/writers](#proposal-default-readers-writers)
    * [Local](#proposal-default-readers-writers-local)
    * [SAAS](#proposal-default-readers-writers-saas)

<a id="objectives"></a>

## Objectives

* Create an SDK to write components in a predictable way
* The SDK should take care of how a component runs and should handle:
  * logging, monitoring and tracing
  * panics
  * graceful shutdown and cancellations
  * health checks
* Provide a set of reusable building blocks to perform common operations, such as reading from a source
* Provision a simple default storage (overridable) to allow transferring information between different components
* Standardise environment aware components: initialisation, execution and shutdown
* Represent information shared between components using [ocsf](https://github.com/ocsf/examples/tree/main) format
  leveraging protobuf and generated code
* Make SDK and components easy to test
* Work towards being able to run components with zero dependencies locally (single binary)

<a id="callouts"></a>

## Callouts

* This document covers ONLY how to create new components via SDK
* The advertised proposal will be worked on incrementally
* Some more technical bits about retries, metrics and so on have been left out for the moment.
  We'll discuss them as a follow-up in RFCs.
* The implementation will likely differ a bit from the proposal as we'll tweak the final the SDK
  while migrating components. We'll track the various changes to the SDK with RFCs.
* The proposed environment variables will be heavily documented and validated
* A default `context.Context` is being currently passed around in the examples for the SDK.
  This will most likely change with an internal `component.Context` which still preserves `context.Context`
  semantics like cancellation etc but also provides stronger typing for context values which will be
  most likely used to get re-usable debugging/monitoring bits around the components. I.E. `component.GetLoggerFromContext(ctx)`.
  This is similar to how lots of workflow engines execute on this, for example [Temporal](https://typescript.temporal.io/api/classes/activity.Context).

<a id="proposal"></a>

## Proposal

<a id="proposal-structure"></a>

### Structure and contributions guidelines

The code for the new SDK will live in this repository in the `sdk/component` directory.

The SDK will have its own go module `github.com/ocurity/dracon/sdk/component` with its own set of dependencies.

> Why?
>
> * having a separate directory helps keeping responsibilities segregated
> * having a separate module helps keeping dependencies separate, we want very few dependencies on the core SDK to boost adoption
> * it's trivial to move the SDK in its own repository, if needed
> * independent versioning from dracon-oss

The structure of the files will follow these guidelines:

```
dracon/
├── ...
├── sdk/
│   └── component/
│       ├── go.mod
│       ├── component.go
│       ├── reader.go
│       ├── writer.go
│       ├── transformer.go
│       ├── ...
│       └── examples/
└── ...

** The file names are an example
```

All the files belonging to the SDK will leave on the root level of the package to simplify import paths.

Unexported code MUST be unexported or live inside the `internal/` directory.

<a id="proposal-component-api"></a>

### Component API

<a id="proposal-component-api-first-pass"></a>

#### First pass - simple components

At its core, regardless of the type of the component, each one of them has to perform 3 actions:

* `Read`: read raw information from some source - filesystem, database, api, ...
* `Transform`: perform operations to enrich the raw data based on the context - adding annotations, call an api, ...
* `Write`: write the enriched results - filesystem, database, api ...

Different components implement these actions in one way or another. Some of them can be NOOP.

The components share information between different actions. This information follows a predictable format of a vulnerability report.

This means that we can provide a minimum set of API to guide the creation of new components. The Component API will be represented by a `Runner` interface which can be satisfied to implement a new component:

```go
// Runner advertises the actions that need to be implemented by a Smithy component.
type Runner interface {
    // Read should be implemented to read Smithy's raw input from an underlying data source.
    Read(ctx context.Context) (*v1.VulnerabilityFinding, error)
    // Transform should be implemented to process and enrich the raw input into an enriched output.
    Transform(ctx context.Context, in *v1.VulnerabilityFinding) (*v1.VulnerabilityFinding, error)
    // Write should be implemented to store the enriched results into an underlying data destination.
    Write(ctx context.Context, out *v1.VulnerabilityFinding) error
}
```

[`v1.VulnerabilityFinding`](https://schema.ocsf.io/1.3.0/classes/vulnerability_finding?extensions=) is used to represent raw and enriched vulnerability reports.

This will help us optimising for compatibility with other tools.

<a id="proposal-component-api-first-pass-create"></a>

##### Creating a new component

In order to create a new component, we can define the following configuration type and its constructor:

```go
// Config contains the component configuration.
type (
    Config struct {}
	
    // ConfigOption can be used to override configuration defaults.
    // For example overriding the default logger.
    ConfigOption func(*Config) error
)

// NewConfig returns a new component configuration with overridable defaults.
func NewConfig(opts ConfigOption...) (*Config, error) { 
    ...
    cfg := &Config{
        // Logger: defaultLogger{},
    }
	
    for _, opt := range opts {
        if err := opt(cfg); err != nil { ... }
    }   
	
    return cfg, nil
}
```

`Config` will contain different dependencies that are used by the component for shared functionalities:

* logging
* panic handling
* metrics
* ...

We'll cover it more in depth later.

Now we can define a basic concrete private component type and its constructor:

```go
type component struct {
    config *Config
}

func New(config *Config) (*component, error) { ... }
```

> Why is this private?
>
> * We make sure to not allow manipulating an instance of a component outside the supplied methods
> * We keep our public SDK API small
> * The only way to get a new component is via constructor so we can ensure that some internals are
>   \>   initialised with defaults

<a id="proposal-component-api-first-pass-running"></a>

##### Running a component

After instantiating a component, this can be executed using the `Run` method.

`Run` will prepare the execution context for a component by running the run actions and providing:

* signal/cancellation termination handling + graceful shutdown
* panic handling
* [setting GOMAXPROCS to match Linux container CPU quota](https://github.com/uber-go/automaxprocs)
* exposing/pushing metrics + registering default metrics
* exposing health checks for long-lived components
* emitting traces
* profiling server
* ...
* supporting any other reliability capability

```go
func (c *component) Run(ctx context.Context, runner Runner) error {
	// set GOMAXPROCS correctly for containers
	// handle metrics/logs/traces/profiler
	// handle panics
	// handle health checks for long-lived components
	
	...

	in, err := runner.Read(ctx)
	if err != nil {
	    return fmt.Errorf("failed to execute Read step: %w", err)
	}

	out, err := runner.Transform(ctx, in)
	if err != nil {
	    return fmt.Errorf("failed to execute Transform step: %w", err)
	}

	if err := runner.Write(ctx, out); err != nil {
	    return fmt.Errorf("failed to execute Write step: %w", err)
	}
	
	return nil
}
```

The code above is just an example; the idea is that `Run` will execute the `Runner` actions
and provide for free a lot of functionalities that are needed to build reliable, debuggable and efficient components.

The actions executed by run will be configurable and have notions of retry policies and custom error handling.

*Not all of these functionalities will be provided from day one. They'll be added one at a time.*

<a id="proposal-component-api-first-pass-defaults"></a>

##### Defaults

The component SDK is designed to provide overridable smart defaults to provide reliability, monitoring and debugging
capabilities for free to adopters.

These functionalities will be shared across actions and configured in `Config`.

We won't cover them all here. Just to give an example, considering logging and panic handling:

```go
type (
    // PanicHandler defines a generic contract for handling panics following the recover semantics.
    PanicHandler interface {
        HandlePanic()
    }

    // Logger exposes an slog.Logger compatible logger contract.
    Logger interface {
        Debug(msg string, keyvals ...any)
        Info(msg string, keyvals ...any)
        Warn(msg string, keyvals ...any)
        Error(msg string, keyvals ...any)
    }
    
    Config struct {
        Logger       Logger
        PanicHandler PanicHandler
    }

    // ConfigOption can be used to override configuration defaults.
    // For example overriding the default logger.
    ConfigOption func(*Config) error
	
    defaultPanicHandler struct {
        logger Logger
    }
)

func (dph defaultPanicHandler) HandlePanic() {
    if err := recover(); err != nil {
        dph.logger.Error("recovered from panic: %v", err)
    }
}

// NewConfig returns a new configuration with overridable defaults.
func NewConfig(opts ConfigOption...) (*Config, error) {
    logger := Logger: log.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
        Level: slog.LevelDebug,
    }))
	
    cfg := &Config{
        Logger: logger,
        PanicHandler: defaultPanicHandler{logger: logger},
    }

    for _, opt := range opts {
        if err := opt(cfg); err != nil { ... }
    }
	
    return cfg, nil
}
```

<a id="proposal-component-api-first-pass-example"></a>

##### Example

Essentially, one has to just implement a `Runner` to create a new component and that will be it.

<details>
  <summary>Simple Example</summary>

```go
package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/ocurity/dracon/sdk/component"
	"github.com/ocurity/dracon/api/gen/com/github/ocsf/ocsf-schema/v1"
)

// simpleComponent parses a Vulnerability Findings report from JSON
// enriches it and stores the processed output to SQLite.
type simpleComponent struct{}

func (s simpleComponent) Read(ctx context.Context) (*v1.VulnerabilityFinding, error) {
	b, err := os.ReadFile("input.json")
	if err != nil {
		return nil, fmt.Errorf("failed to read input: %w", err)
	}

	// Parse to *v1.VulnerabilityFinding

	return &v1.VulnerabilityFinding{}, nil
}

func (s simpleComponent) Write(ctx context.Context, out *v1.VulnerabilityFinding) error {
	// Write finding to SQLite
	return nil
}

func (s simpleComponent) Transform(ctx context.Context, in *v1.VulnerabilityFinding) (output *v1.VulnerabilityFinding, err error) {
	// Enrich and transform Input in Output
	return &v1.VulnerabilityFinding{}, nil
}

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cfg, err := component.NewConfig()
	if err != nil {
		log.Fatalf("could not create configuration: %v", err)
	}

	c, err := component.New(cfg)
	if err != nil {
		log.Fatalf("could not create component: %v", err)
	}

	if err := c.Run(ctx, &simpleComponent{}); err != nil {
		log.Fatalf("could not run component: %v", err)
	}
}
```

</details>

<a id="proposal-component-api-second-pass"></a>

### Second pass - environment aware components

<a id="proposal-component-api-second-pass-why"></a>

#### Why the first pass is not enough

The first pass takes us closer to our objectives, but when we think about maintaining multiple components it can make
our lives more complicated, why so?

* Code configuration for a component is required. One has to create (or re-use) a `Runner` and initialise its various bits:
  Reader, Writer and Transformer. Often this is done based on environment variables lookup. 💡 *Ideally this is done by the SDK.*
* Updating a component to read from, for example, `S3` to `MongoDB` would require a code change. A new
  reader should be initialised correctly and swapped in the implementation of `Runner.` 💡 *Ideally this can be done automatically by the SDK.*
* Customisation for components could be simplified and made more dynamic. 💡 *Perhaps the SDK core doesn't change but we can dynamically use a reader, writer or transformer*

> What if we re-think the runner to be more environment aware and support automatically bootstrapping and serving
> different readers, writers and transformers while not requiring significant code changes? Perhaps just a dependency bump
> will do.

<a id="proposal-component-api-second-pass-env-aware-runner"></a>

#### A new environment aware runner

To introduce a notion of environment aware runner, we need to think about what environment we depend on.

We need to be:

* aware of which readers, writers and transformers to use
* which extra parameters and steps are needed to bring the latter up

It's clear that to build a platform that makes sense, we need to define at least three basic environment variables:

| Environment Variable | Required | Type   | Default | Possible Values                        |
|----------------------|----------|--------|---------|----------------------------------------|
| READER               | yes      | string | none    | \[s3, mongodb, ...]                    |
| WRITER               | yes      | string | none    | \[s3, mongodb, ...]                    |
| TRANSFORMER          | yes      | string | none    | \[codeowners, deduplication, ...]      |

And we need an association between Reader, Writer and Transformer names to which actions to bootstrap and execute.

Once the association has been figured out, we can boostrap the actions for the component. To do so, we need to
have a boostrap step on each action.

We could revisit our `Run` functions to do exactly that.

```go
type (
    Reader interface {
        Construct(ctx context.Context) (Reader, error)
        Read(ctx context.Context) (*v1.VulnerabilityFinding, error)	
    }
    Writer interface {
        Construct(ctx context.Context) (Writer, error)
        Write(ctx context.Context, out *v1.VulnerabilityFinding) error
    }
    Transformer interface {
        Construct(ctx context.Context) (Transformer, error)
        Transform(ctx context.Context, in *v1.VulnerabilityFinding) (*v1.VulnerabilityFinding, error)
    }
	
    Config struct {
        ...
        // extra fields for configuration like showed in the first pass
        ...
        Reader Reader
        Writer Writer
        Transformer Transformer
    }
)

func Run(ctx context.Context, cfg *Config...) error {
    // If cfg not passed, get a default configuration.
    // In this step also look for READER, WRITER and TRANSFORMER configurations.
    cfg, err := NewConfig()
    if err != nil {
        return fmt.Errorf("failed to create new configuration: %w", err)
    }

    // 1. What's the value of READER|WRITER|TRANSFORMER environment variable?
    // 2. Does it match a known reader?
    ok, err := checkEnv(ctx)
    if err != nil {
        return fmt.Errorf("could not construct reader: %w", err)
    }
	
    r, err := reader.Construct(ctx)
    if err != nil {
        return fmt.Errorf("could not construct reader: %w", err)
    }
	
    w, err := writer.Construct(ctx)
    if err != nil {
        return fmt.Errorf("could not construct writer: %w", err)
    }
	
    t, err := transformer.Construct(ctx)
    if err != nil {
        return fmt.Errorf("could not construct writer: %w", err)
    }
	
    // set GOMAXPROCS correctly for containers
    // handle metrics/logs/traces/profiler
    // handle panics
    // handle health checks for long-lived components
    
    ...
    
    // The following steps most likely have to be run concurrently and indefinitely
    // for long-lived actions e.g. read from a Kafka topic.
	
    in, err := r.Read(ctx)
    if err != nil {
        return fmt.Errorf("failed to execute Read step: %w", err)
    }
    
    out, err := t.Transform(ctx, in)
    if err != nil {
        return fmt.Errorf("failed to execute Transform step: %w", err)
    }
    
    if err := w.Write(ctx, out); err != nil {
        return fmt.Errorf("failed to execute Write step: %w", err)
    }
    
    return nil
}
```

This is overly simplified but the idea is to automatically create actions and perform transformations
like in the first pass and preserving all the niceties.

<a id="proposal-component-api-second-pass-create-new-component"></a>

#### Creating a new component

Now creating a new component would require the following code:

```go
func main() {
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()
    
    if err := component.Run(ctx); err != nil {
        log.Fatalf("could not run component: %v", err)
    }
}
```

Assuming that the environment is set and the actions are already existing, that would be it.

<a id="proposal-component-api-second-pass-example"></a>

##### Example

But how do we connect the dots if we want to write a custom action?

Assuming a scenario where we want to read from MongoDB, then enrich with CODEOWNERS tags and then store the results back to MongoDB.

We can hypothetically create a MongoDB reader/writer:

```go
package mongodb

import "context"

type ReadWriter struct {
    // db *mongo.DB
}

func (rw *ReadWriter) Construct(ctx context.Context) (*ReadWriter, error) {
    // 1. Lookup for environment variables to create a new MongoDB client:
    //  - collection name
    //  - url
    //  - ...
    // 2. Create the client.
    return &ReadWriter{}, nil
}

func (rw *ReadWriter) Read(ctx context.Context) (*v1.VulnerabilityFinding, error) {
    // 1. Read document(s) from collection
    // 2. Parse to *v1.VulnerabilityFinding
    return v1.VulnerabilityFinding{}, nil
}

func (rw *ReadWriter) Write(ctx context.Context, out *v1.VulnerabilityFinding) error {
    // 1. Parse output to MongoDB document
    // 2. Write to collection
    return nil
}
```

Assuming that the environment is set:

| Environment Variable          | Required | Type   | Default | Possible Values |
|-------------------------------|----------|--------|---------|-----------------|
| RW\_MONGODB\_DB\_URL          | yes      | string | none    | -               |
| RW\_MONGODB\_COLLECTION\_NAME | yes      | string | none    | -               |

And a CODEOWNERS enrichment transformation:

```go
package codeowners

import "context"

type Transformer struct {}

func (t *Transformer) Construct(ctx context.Context) (*Transformer, error) {
	// 1. Lookup for environment variables to create a new CODEOWNERS transformer:
	//  - repo name
	//  - url
	//  - ...
	return &Transformer{}, nil
}

func (t *Transformer) Transform(ctx context.Context, in *v1.VulnerabilityFinding) (*v1.VulnerabilityFinding, error) {
	// perform transformation
	return &v1.VulnerabilityFinding{}, nil
}
```

Assuming that the environment is set:

| Environment Variable                | Required | Type   | Default | Possible Values |
|-------------------------------------|----------|--------|---------|-----------------|
| TRANSFORMER\_CODEOWNERS\_REPO\_NAME | yes      | string | none    | -               |
| TRANSFORMER\_CODEOWNERS\_VCS\_URL   | yes      | string | none    | -               |

<a id="proposal-component-api-second-pass-wrap"></a>

##### Wrapping it up

Going back to the previous example on how to run our component, we'd basically need:

```go
func main() {
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()
    
    if err := component.Run(ctx); err != nil {
        log.Fatalf("could not run component: %v", err)
    }
}
```

And the following environment:

| Environment Variable                | Value                             |
|-------------------------------------|-----------------------------------|
| READER                              | mongodb                           |
| WRITER                              | mongodb                           |
| TRANSFORMER                         | codeowners                        |
| RW\_MONGODB\_DB\_URL                | mongodb://localhost:27017/reports |
| RW\_MONGODB\_COLLECTION\_NAME       | vuln-reports                      |
| TRANSFORMER\_CODEOWNERS\_REPO\_NAME | dracon                            |
| TRANSFORMER\_CODEOWNERS\_VCS\_URL   | github.com                        |

<a id="proposal-components-catalog"></a>

### Components catalog

Having this simple, transparent and solid foundation for creating new components means that we can
now provide more out-of-the box ways to get re-usable components and bits.

We can create a separate directory where we can contribute new:

* readers
* writers
* transformers
* components

Each one of these components and actions will have their own go module and dependencies
to minimise pollution and making individual maintenance easier.

For instance, we can provide re-usable readers, writers and transformers to read/write/transform
data from/to different source/destinations. Then, such bits can be re-used to create components.

For example:

```
dracon/
├── readwriter/
│   ├── sqlite/
│   │   ├── go.mod
│   │   ├── reader.go
│   │   └── writer.go
│   ├── remote-rpc/
│   │   ├── go.mod
│   │   ├── reader.go
│   │   └── writer.go
│   └── elasticsearch/
│       ├── go.mod
│       ├── reader.go
│       └── writer.go
│   └── transformer/
│       ├── atom-enricher/
│       │   └── go.mod
│       ├── tagger/
│       │   └── go.mod
│       └── prioritiser/
│           └── go.mod
```

The components will follow the known structure and just re-use actions or define their own
to implement the goals of the component.

One can see a scenario where an SQLite reader/writer can be used between multiple components
to read/write Vulnerability Reports from/to.

<a id="proposal-default-readers-writers"></a>

### Default readers/writers

We want to provide as soon as possible two reader/writer solutions for both local and saas environments.

The reason is to speed up development as well as going closer to having a way to run components more easily.

<a id="proposal-default-readers-writers-local"></a>

#### Local

For local settings, we want to use an [Embedded Database](https://en.wikipedia.org/wiki/Embedded_database) which
is lightweight and doesn't require a standalone application or container to run.

This should drastically help running components locally and driving adoption given the benefits in terms
of development experience.

Ultimately, this will be a new `Reader/Writer`.

Our decision is to experiment with [SQLite](https://www.sqlite.org/).

Some of the reasons while we decided to experiment with it are:

* Self-contained: SQLite is a single library that doesn't require separate server software or configuration, making it easy to integrate directly into applications.
* Lightweight: With a small footprint, ideal for embedded systems and applications with limited resources.
* Zero-configuration: SQLite doesn't require a server to set up or maintain, simplifying deployment and reducing maintenance costs.
* Cross-platform: It works on multiple operating systems, ensuring broad compatibility for embedded systems.
* SQL compatible API: Knowing SQL is all it takes to use SQLite. There are several packages to interact with SQLite like [go-sqlite3](https://github.com/mattn/go-sqlite3).
* Reliable: SQLite supports atomic transactions, ensuring data integrity, even in case of crashes or power failures.
* Fast for small-scale applications: It is optimized for read-heavy workloads and provides high performance for smaller datasets, making it ideal for lightweight use cases.
* Public domain: SQLite is free and open-source, with no licensing fees, which makes it attractive for commercial and non-commercial use.

It's great in environments where a database service is not necessary.

We did a small experiment [here](https://github.com/ocurity/sqlite-go-spike).

<a id="proposal-default-readers-writers-saas"></a>

#### SAAS

For `SAAS` we decided to invest in a more sophisticated system to abstract reads/writes to different places.

We'll create a [gRPC](https://grpc.io/) API Gateway which will provide an RPC to store Vulnerability Reports into different data storages.
This RPC will capture the flags needed to understand on which data target to write.

This will boost massively productivity as all components will be able to re-use this single RPC
and we'll be able to hide the complexity away and add extra niceties.

Ultimately, this will be a new `Reader/Writer` part of the `Runner` contract.

We'll discuss this further in a separate RFC.
