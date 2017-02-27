package application

import "github.com/YuukanOO/alfredo/identity/domain"

// AuthApplicationService is the auth application service.
type AuthApplicationService struct {
	userRepository domain.UserRepository
	cryptoService  domain.CryptoService
}

// NewAuthApplicationService instantiates a new auth application service.
func NewAuthApplicationService(registry domain.Registry) *AuthApplicationService {
	return &AuthApplicationService{
		userRepository: registry.UserRepository(),
		cryptoService:  registry.CryptoService(),
	}
}

// Register a new user in the system.
func (s *AuthApplicationService) Register(email, password string) (*domain.User, error) {
	usr, err := domain.Register(s.userRepository, s.cryptoService, email, password)

	if err != nil {
		return nil, err
	}

	if err = s.userRepository.Add(usr); err != nil {
		return nil, err
	}

	return usr, nil
}
