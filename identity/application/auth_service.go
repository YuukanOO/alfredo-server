package application

import "github.com/YuukanOO/alfredo/identity/domain"

// AuthService is the auth application service.
type AuthService struct {
	userRepository domain.UserRepository
	cryptoService  domain.CryptoService
}

// NewAuthService instantiates a new auth application service.
func NewAuthService(registry domain.Registry) *AuthService {
	return &AuthService{
		userRepository: registry.UserRepository(),
		cryptoService:  registry.CryptoService(),
	}
}

// Register a new user in the system.
func (s *AuthService) Register(email, password string) (*domain.User, error) {
	usr, err := domain.Register(s.userRepository, s.cryptoService, email, password)

	if err != nil {
		return nil, err
	}

	if err = s.userRepository.Add(usr); err != nil {
		return nil, err
	}

	return usr, nil
}
