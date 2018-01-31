package di

import (
    "errors"
    "sync"
    "time"

    "github.com/patrickmn/go-cache"
)

type (

    // Container is the manager behind the dependency injection
    // it holds the configuration to build a dependency and any
    // cached objects that should be shared.
    Container struct {
        Dependencies *sync.Map
        Cache        CacheInterface
    }

    // CacheInterface is the contract for the cache, again
    // you can roll out your own otherwise the standard cache
    // will be used.
    //
    // The standard cache in this case is the go-cache package
    // imported above.
    CacheInterface interface {
        Set(string, interface{}, time.Duration)
        Get(string) (interface{}, bool)
    }
)

// NewContainer creates a new container
func NewContainer() *Container {
    return &Container{
        Dependencies: new(sync.Map),
        Cache:        cache.New(5*time.Minute, 10*time.Minute),
    }
}

// Add is used to register a new dependency, calling add will simply
// add the given dependency to the map, it will not build it
func (c *Container) Add(d *Dependency) {
    c.Dependencies.Store(d.Name, d)
}

// MustGet will either return the underlying interface or
// panic
func (c *Container) MustGet(n string) interface{} {
    dep, err := c.Get(n)

    if err != nil {
        panic(err)
    }

    return dep
}

// Get returns the instance of the dependency.
//
// If its a singleton it will always build before return, if its
// a shared instance, we check the cache first and return that
// if available, otherwise, return a new instance and store in
// the cache.
func (c *Container) Get(n string) (interface{}, error) {
    // Check if the dependency exists first
    dep, err := c.Fetch(n)

    if err != nil {
        return nil, err
    }

    // If its a shared dependency, we should check the
    // cache before building it
    if dep.Shared {
        i, err := c.getFromCache(n)

        // If its not found in there, we need to build it
        if err == nil {
            return i, nil
        }
    }

    // Otherwise we can build and return
    // again, if its a shared, we need to
    // add it back into the cache for the next call
    i, err := dep.Build(c)

    if err == nil && dep.Shared {
        c.Cache.Set(dep.Name, i, cache.NoExpiration)
    }

    return i, err
}

// Fetch is used to return the dependency object by its name.
//
// Note: this will not be a built dependency, just the underlying
// entry struct.
func (c *Container) Fetch(n string) (*Dependency, error) {
    dep, ok := c.Dependencies.Load(n)

    if !ok {
        return &Dependency{}, errors.New("Unable to find dependency with name " + n)
    }

    return dep.(*Dependency), nil
}

// Attempts to fetch a dependency from the cache
func (c *Container) getFromCache(n string) (interface{}, error) {
    inst, ok := c.Cache.Get(n)

    if !ok {
        return nil, errors.New("Dependency not found in the cache")
    }

    return inst, nil
}
