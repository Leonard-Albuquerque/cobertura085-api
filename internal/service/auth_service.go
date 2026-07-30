package service

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/Leonard-Albuquerque/cobertura085-api/internal/config"
	"github.com/Leonard-Albuquerque/cobertura085-api/internal/model"
	"github.com/Leonard-Albuquerque/cobertura085-api/internal/repository"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

var (
	ErrInvalidCredentials     = errors.New("credenciais inválidas")
	ErrUserInactive           = errors.New("conta de usuário inativa")
	ErrEmailAlreadyInUse      = errors.New("e-mail já está em uso")
	ErrStoreNameAlreadyExists = errors.New("já existe uma loja cadastrada com este nome ou slug")
	ErrInvalidRefreshToken    = errors.New("refresh token inválido, expirado ou revogado")
)

type RegisterInput struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

type LoginInput struct {
	Email     string  `json:"email" binding:"required,email"`
	Password  string  `json:"password" binding:"required"`
	UserAgent *string `json:"userAgent,omitempty"`
	IP        *string `json:"ip,omitempty"`
}

type RefreshInput struct {
	RefreshToken string  `json:"refreshToken" binding:"required"`
	UserAgent    *string `json:"userAgent,omitempty"`
	IP           *string `json:"ip,omitempty"`
}

type AuthResponse struct {
	AccessToken  string      `json:"accessToken"`
	RefreshToken string      `json:"refreshToken"`
	User         *model.User `json:"user"`
}

type AuthService interface {
	Register(ctx context.Context, input RegisterInput) (*model.User, error)
	Login(ctx context.Context, input LoginInput) (*AuthResponse, error)
	RefreshToken(ctx context.Context, input RefreshInput) (*AuthResponse, error)
	Logout(ctx context.Context, refreshTokenStr string) error
	GetProfile(ctx context.Context, userID string) (*model.User, error)
}

type authService struct {
	userRepo         repository.UserRepository
	storeRepo        repository.StoreRepository
	refreshTokenRepo repository.RefreshTokenRepository
	jwtService       JWTService
	cfg              *config.Config
}

func NewAuthService(
	userRepo repository.UserRepository,
	storeRepo repository.StoreRepository,
	refreshTokenRepo repository.RefreshTokenRepository,
	jwtService JWTService,
	cfg *config.Config,
) AuthService {
	return &authService{
		userRepo:         userRepo,
		storeRepo:        storeRepo,
		refreshTokenRepo: refreshTokenRepo,
		jwtService:       jwtService,
		cfg:              cfg,
	}
}

func (s *authService) Register(ctx context.Context, input RegisterInput) (*model.User, error) {
	name := strings.TrimSpace(input.Name)
	if name == "" {
		return nil, errors.New("o nome não pode ser vazio")
	}

	// 1. Gerar slug do nome (ex: "Lojão do menino" -> "lojao-do-menino")
	slug := Slugify(name)
	if slug == "" {
		return nil, errors.New("nome da loja inválido para geração de slug")
	}

	// 2. Verificar unicidade de Name e Slug da Store
	existingStore, err := s.storeRepo.FindByNameOrSlug(ctx, name, slug)
	if err == nil && existingStore != nil {
		return nil, ErrStoreNameAlreadyExists
	}

	// 3. Verificar unicidade de E-mail
	existingUser, err := s.userRepo.FindByEmail(ctx, input.Email)
	if err == nil && existingUser != nil {
		return nil, ErrEmailAlreadyInUse
	}

	// 4. Hash da Senha
	passwordHash, err := s.jwtService.HashPassword(input.Password)
	if err != nil {
		return nil, fmt.Errorf("erro ao processar senha: %w", err)
	}

	now := time.Now().UTC()
	storeID := uuid.New().String()
	userID := uuid.New().String()

	store := &model.Store{
		ID:                  storeID,
		Slug:                slug,
		Name:                name,
		Whatsapp:            "5585999999999",
		Address:             "Endereço da Loja, Fortaleza - CE",
		OperatingHours:      "Segunda a Sexta: 08:00 às 18:00",
		PickupEnabled:       true,
		DeliveryTimeDefault: "2 horas",
		QRToken:             uuid.New().String(),
		CreatedAt:           now,
		UpdatedAt:           now,
	}

	user := &model.User{
		ID:           userID,
		StoreID:      storeID,
		Email:        input.Email,
		PasswordHash: passwordHash,
		Name:         &name,
		IsActive:     true,
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	// 5. Executar Transação no Banco: Criar Store -> Criar User
	db := s.userRepo.GetDB()
	if db != nil {
		err = db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
			if err := s.storeRepo.CreateTx(ctx, tx, store); err != nil {
				return err
			}
			if err := s.userRepo.CreateTx(ctx, tx, user); err != nil {
				return err
			}
			return nil
		})
		if err != nil {
			return nil, fmt.Errorf("falha ao salvar cadastro no banco: %w", err)
		}
	} else {
		// Fallback se DB for nulo (ex: ambiente sem DB)
		if err := s.storeRepo.Create(ctx, store); err != nil {
			return nil, err
		}
		if err := s.userRepo.Create(ctx, user); err != nil {
			return nil, err
		}
	}

	user.Store = store
	return user, nil
}

func (s *authService) Login(ctx context.Context, input LoginInput) (*AuthResponse, error) {
	// 1. Buscar Usuário pelo Email
	user, err := s.userRepo.FindByEmail(ctx, input.Email)
	if err != nil {
		return nil, ErrInvalidCredentials
	}

	// 2. Verificar se conta está ativa
	if !user.IsActive {
		return nil, ErrUserInactive
	}

	// 3. Validar Senha
	if !s.jwtService.CheckPassword(input.Password, user.PasswordHash) {
		return nil, ErrInvalidCredentials
	}

	// 4. Atualizar lastLoginAt
	now := time.Now().UTC()
	_ = s.userRepo.UpdateLastLogin(ctx, user.ID, now)
	user.LastLoginAt = &now

	// 5. Gerar Access Token e Refresh Token
	return s.generateTokenPair(ctx, user, input.UserAgent, input.IP)
}

func (s *authService) RefreshToken(ctx context.Context, input RefreshInput) (*AuthResponse, error) {
	// 1. Calcular hash do Refresh Token recebido
	tokenHash := s.jwtService.HashToken(input.RefreshToken)

	// 2. Buscar Refresh Token no banco de dados
	storedToken, err := s.refreshTokenRepo.FindByTokenHash(ctx, tokenHash)
	if err != nil {
		return nil, ErrInvalidRefreshToken
	}

	// 3. Validar expiração e revogação
	now := time.Now().UTC()
	if storedToken.RevokedAt != nil || now.After(storedToken.ExpiresAt) {
		return nil, ErrInvalidRefreshToken
	}

	// 4. Revogar o token atual (Rotação de Tokens)
	_ = s.refreshTokenRepo.Revoke(ctx, storedToken.ID)

	// 5. Buscar usuário e checar status
	user, err := s.userRepo.FindByID(ctx, storedToken.UserID)
	if err != nil || !user.IsActive {
		return nil, ErrUserInactive
	}

	// 6. Gerar novo par de tokens
	return s.generateTokenPair(ctx, user, input.UserAgent, input.IP)
}

func (s *authService) Logout(ctx context.Context, refreshTokenStr string) error {
	tokenHash := s.jwtService.HashToken(refreshTokenStr)
	storedToken, err := s.refreshTokenRepo.FindByTokenHash(ctx, tokenHash)
	if err != nil {
		return nil // Se não encontrou, já está inativo
	}

	return s.refreshTokenRepo.Revoke(ctx, storedToken.ID)
}

func (s *authService) GetProfile(ctx context.Context, userID string) (*model.User, error) {
	return s.userRepo.FindByID(ctx, userID)
}

// Auxiliar para gerar Access Token (JWT) e Refresh Token seguro
func (s *authService) generateTokenPair(ctx context.Context, user *model.User, userAgent, ip *string) (*AuthResponse, error) {
	accessDuration, err := time.ParseDuration(s.cfg.JWTAccessExpiration)
	if err != nil {
		accessDuration = 15 * time.Minute
	}

	refreshDuration, err := time.ParseDuration(s.cfg.JWTRefreshExpiration)
	if err != nil {
		refreshDuration = 7 * 24 * time.Hour
	}

	// 1. Access Token
	accessToken, err := s.jwtService.GenerateAccessToken(user.ID, user.StoreID, user.Email, accessDuration)
	if err != nil {
		return nil, err
	}

	// 2. Refresh Token Aleatório
	rawRefreshToken := generateRandomString(32)
	refreshTokenHash := s.jwtService.HashToken(rawRefreshToken)

	now := time.Now().UTC()
	var ipHash *string
	if ip != nil && *ip != "" {
		h := s.jwtService.HashToken(*ip)
		ipHash = &h
	}

	refreshTokenRecord := &model.RefreshToken{
		ID:        uuid.New().String(),
		UserID:    user.ID,
		TokenHash: refreshTokenHash,
		ExpiresAt: now.Add(refreshDuration),
		UserAgent: userAgent,
		IPHash:    ipHash,
		CreatedAt: now,
	}

	if err := s.refreshTokenRepo.Create(ctx, refreshTokenRecord); err != nil {
		return nil, fmt.Errorf("falha ao salvar refresh token: %w", err)
	}

	return &AuthResponse{
		AccessToken:  accessToken,
		RefreshToken: rawRefreshToken,
		User:         user,
	}, nil
}

func generateRandomString(n int) string {
	bytes := make([]byte, n)
	_, _ = rand.Read(bytes)
	return hex.EncodeToString(bytes)
}
