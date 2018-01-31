package di

type (
    // Dependency represents instructions for the container
    // on how to build and maintain a dependency, a dependency
    // can come in several forms.
    //
    // It could be a singleton, in which case will be built on the fly
    // every time it is requested.
    //
    // It could also be a shared dependency which will only build on
    // the first request, any subsequent get's will be sourced from
    // the initial instance that the first request initiates.
    Dependency struct {
        // Name needs to be unique and will be used in calls to fetch
        // the dependency
        Name string
        // Shared bool denotes whether this is a shared instance
        // or a singleton
        Shared bool
        // Build see BuildFunc
        Build BuildFunc
    }

    // BuildFunc function is the instruction the container will follow
    // when it comes to build the given dependency.
    BuildFunc func(*Container) (interface{}, error)
)
