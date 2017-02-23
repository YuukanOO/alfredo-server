package domain

// Registry represents the domain registry needed for the application services.
type Registry interface {
	UserRepository() UserRepository
	CryptoService() CryptoService
}
