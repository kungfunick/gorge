package di

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type (

	// Used a base within the TestDependency
	// this is used to test inheritance.
	BaseDependency struct{}

	// Used to test basic functionality
	TestDependency struct {
		BaseDependency
		Count int
	}

	// Used to test dependency sharing
	ParamBag struct {
		T *TestDependency
	}
)

func NewParamBag(d *TestDependency) *ParamBag {
	return &ParamBag{T: d}
}

func (b *BaseDependency) foo() string {
	return "bar"
}

func TestAddAndGetDependency(t *testing.T) {
	c := NewContainer()

	c.Add(&Dependency{
		Name:   "TestDependency",
		Shared: false,
		Build: func(ctx *Container) (interface{}, error) {
			return "test", nil
		},
	})

	out, err := c.Get("TestDependency")

	assert.Nil(t, err)
	assert.Equal(t, "test", out.(string))
}

func TestGetWorksWithInheritance(t *testing.T) {
	c := NewContainer()

	c.Add(&Dependency{
		Name:   "TestDependency",
		Shared: false,
		Build: func(ctx *Container) (interface{}, error) {
			return TestDependency{}, nil
		},
	})

	out, err := c.Get("TestDependency")

	assert.Nil(t, err)

	dep := out.(TestDependency)

	assert.Equal(t, "bar", dep.foo())
}

func TestSharedInstance(t *testing.T) {
	c := NewContainer()

	c.Add(&Dependency{
		Name:   "TestDependency",
		Shared: true,
		Build: func(ctx *Container) (interface{}, error) {
			return &TestDependency{}, nil
		},
	})

	out, _ := c.Get("TestDependency")
	dep := out.(*TestDependency)
	dep.Count = 2

	tst := out.(*TestDependency)

	assert.Equal(t, 2, tst.Count)
}

func TestDependencyWithDependency(t *testing.T) {
	c := NewContainer()

	c.Add(&Dependency{
		Name:   "TestDependency",
		Shared: false,
		Build: func(ctx *Container) (interface{}, error) {
			return &TestDependency{}, nil
		},
	})

	c.Add(&Dependency{
		Name:   "ParamBag",
		Shared: true,
		Build: func(ctx *Container) (interface{}, error) {
			return NewParamBag(ctx.MustGet("TestDependency").(*TestDependency)), nil
		},
	})

	out, err := c.Get("ParamBag")

	assert.Nil(t, err)

	p := out.(*ParamBag)

	assert.Equal(t, "bar", p.T.foo())
}

func BenchmarkGetDependency(b *testing.B) {
	c := NewContainer()

	c.Add(&Dependency{
		Name:   "Test",
		Shared: false,
		Build: func(ctx *Container) (interface{}, error) {
			return &TestDependency{}, nil
		},
	})

	for i := 0; i < b.N; i++ {
		c.MustGet("Test")
	}
}

func BenchmarkGetDependencyCached(b *testing.B) {
	c := NewContainer()

	c.Add(&Dependency{
		Name:   "Test",
		Shared: true,
		Build: func(ctx *Container) (interface{}, error) {
			return &TestDependency{}, nil
		},
	})

	for i := 0; i < b.N; i++ {
		c.MustGet("Test")
	}
}

func BenchmarkWithMultipleDependencies(b *testing.B) {
	c := NewContainer()

	c.Add(&Dependency{
		Name:   "TestDependency",
		Shared: true,
		Build: func(ctx *Container) (interface{}, error) {
			return &TestDependency{}, nil
		},
	})

	c.Add(&Dependency{
		Name:   "Test",
		Shared: true,
		Build: func(ctx *Container) (interface{}, error) {
			return NewParamBag(ctx.MustGet("TestDependency").(*TestDependency)), nil
		},
	})

	for i := 0; i < b.N; i++ {
		c.MustGet("Test")
	}
}
