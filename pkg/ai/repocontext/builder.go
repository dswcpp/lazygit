package repocontext

// Builder is implemented by the GUI layer to provide current repository state.
// Injecting this interface keeps AI packages free of GUI dependencies.
type Builder interface {
	Build() RepoContext
}
