package errx

const (
	pkgNamespace = "errx"
)

var DuplicateFallbackError = NewError("ERR_1", "Cannot create new Error that has same code with Fallback Error",
	WithNamespace(pkgNamespace))
