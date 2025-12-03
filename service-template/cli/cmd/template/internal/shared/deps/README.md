**{{SERVICE_TITLE}}**

# Table Of Content

- [Table Of Content](#table-of-content)
    - [Deps](#deps)
    - [Sample limitation of pool of kafka](#sample-limitation-of-pool-of-kafka)

## Deps

On this layer propose for grouping depency and setting close function after we termiate the system. On this part to we
can setting for thread pool to so we can create the limitation of thread / go-routine onn every pool like Kafka, Redis,
etc.

```go
package deps

type (
	Deps struct {
		dig.In
		Config            *config.Configuration
		RedisWrapper      *RedisWrapper
		Kafka             *KafkaWrapper
		WorkerPoolWrapper *WorkerPoolWrapper
		Logger            logs.Logger
		Metric            metrics.MonitorStatsd
		Echo              *echo.Echo
		Mongo             *mongo.Database
		PredefinedCache   *PredefinedCache
	}
)

// this function is mandatory if you have connection session like for mongoDB, redis, etc so can be recicle the session
func (impl *Deps) Close() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	impl.RedisWrapper.Close()
	impl.Kafka.Close()
	impl.Mongo.Client().Disconnect(ctx)
	impl.Echo.Close()
	impl.PredefinedCache.Close()
	return nil
}


// register new method of limitation worker pool on here
func Register(container *dig.Container) error {
	if err := container.Provide(NewRedisWrapper); err != nil {
		return err
	}

	if err := container.Provide(NewWorkerPoolWrapper); err != nil {
		return err
	}

	// - register predefined
	if err := container.Provide(NewPredefined); err != nil {
		return err
	}

	// - register kafka
	if err := container.Provide(NewKafka); err != nil {
		return err
	}

	return nil
}

```

## Sample limitation of pool of kafka

This code is sample setting pool of kafka

```go
package deps

type (
	WorkerPoolWrapper struct {
		SearchRequestPool *pool.Pool
		KafkaGeneralPool  *pool.Pool
	}
)

func NewWorkerPoolWrapper(cfg *config.Configuration) (*WorkerPoolWrapper, error) {
	var (
		p1 = pool.New().WithMaxGoroutines(1) // value can be change via config if you need (can costomize)
		p2 = pool.New().WithMaxGoroutines(1)
	)

	return &WorkerPoolWrapper{
		SearchRequestPool: p1,
		KafkaGeneralPool:  p2,
	}, nil
}
```
