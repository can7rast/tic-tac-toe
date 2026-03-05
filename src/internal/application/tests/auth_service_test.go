package tests_test

import (
	"context"
	"errors"
	"github.com/stretchr/testify/mock"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"school21/internal/application"
	"school21/internal/domain"
	"school21/internal/infrastructure/datasource/mocks" // ← путь к твоему моку
	"school21/internal/web/dto"
	"school21/pkg"
)

func TestAuthService_Login_Success(t *testing.T) {
	mockRepo := mocks.NewUserRepository(t)

	// 2. Создаём сервис, передаём ему мок
	service := application.NewAuthService(mockRepo)

	// 3. Подготавливаем тестовые данные
	ctx := context.Background()
	req := dto.SignUpRequest{
		Login:    "alice",
		Password: "secret123",
	}

	// Хэш пароля, который мы будем "имитировать"
	correctHash, _ := pkg.HashPassword("secret123")

	userID := uuid.New()

	// 4. Говорим моку: когда вызовут FindByLogin с этими аргументами — верни вот это
	mockRepo.EXPECT().
		FindByLogin(mock.Anything, "alice").
		Return(&domain.User{
			ID:           userID,
			Login:        "alice",
			PasswordHash: correctHash,
		}, nil).
		Once() // ровно один раз

	// 5. Выполняем то, что тестируем
	user, err := service.Login(ctx, req)

	// 6. Проверяем результат
	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, userID, user.ID)
	assert.Equal(t, "alice", user.Login)
}

func TestAuthService_Login_UserNotFound(t *testing.T) {
	mockRepo := mocks.NewUserRepository(t)
	service := application.NewAuthService(mockRepo)

	req := dto.SignUpRequest{
		Login:    "unknown",
		Password: "whatever",
	}

	// Пользователь не найден: репозиторий возвращает nil, nil
	mockRepo.EXPECT().
		FindByLogin(mock.Anything, "unknown").
		Return(nil, nil).
		Once()

	user, err := service.Login(context.Background(), req)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "user not found")
	assert.Nil(t, user)
}

func TestAuthService_Login_DatabaseError(t *testing.T) {
	mockRepo := mocks.NewUserRepository(t)
	service := application.NewAuthService(mockRepo)

	req := dto.SignUpRequest{
		Login:    "alice",
		Password: "secret123",
	}

	dbErr := errors.New("connection timeout")

	// Ошибка базы данных
	mockRepo.EXPECT().
		FindByLogin(mock.Anything, "alice").
		Return(nil, dbErr).
		Once()

	user, err := service.Login(context.Background(), req)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "data base error") // если у тебя есть такая обёртка
	assert.Nil(t, user)
}

func TestAuthService_SignUp_Success(t *testing.T) {
	mockRepo := mocks.NewUserRepository(t)
	service := application.NewAuthService(mockRepo)

	req := dto.SignUpRequest{
		Login:    "alice",
		Password: "secret123",
	}

	mockRepo.EXPECT().
		FindByLogin(mock.Anything, "alice").
		Return(nil, nil).
		Once()
	mockRepo.EXPECT().
		Create(mock.Anything, mock.MatchedBy(func(u *domain.User) bool {
			if u == nil {
				return false
			}
			// Проверяем, что login правильный
			if u.Login != "alice" {
				return false
			}
			// Проверяем, что ID сгенерирован (не пустой)
			if u.ID == uuid.Nil {
				return false
			}
			// Проверяем, что пароль захэширован (не равен открытому)
			if u.PasswordHash == "secret123" || u.PasswordHash == "" {
				return false
			}
			return true
		})).
		Return(nil).
		Once()

	err := service.SignUp(context.Background(), req)

	assert.NoError(t, err)
}

func TestAuthService_SignUp_DatabaseError(t *testing.T) {
	mockRepo := mocks.NewUserRepository(t)
	service := application.NewAuthService(mockRepo)
	req := dto.SignUpRequest{
		Login:    "alice",
		Password: "secret123",
	}

	mockRepo.EXPECT().
		FindByLogin(mock.Anything, "alice").
		Return(nil, errors.New("data base error")).
		Once()

	mockRepo.EXPECT().
		Create(mock.Anything, mock.AnythingOfType("*domain.User")).Times(0)

	err := service.SignUp(context.Background(), req)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "data base error")

}

func TestAuthService_SignUp_UserAlreadyExists(t *testing.T) {
	mockRepo := mocks.NewUserRepository(t)
	service := application.NewAuthService(mockRepo)
	req := dto.SignUpRequest{
		Login:    "alice",
		Password: "secret123",
	}
	mockRepo.EXPECT().
		FindByLogin(mock.Anything, "alice").
		Return(&domain.User{
			ID:           uuid.New(),
			Login:        "alice",
			PasswordHash: "secret123",
		}, nil).
		Once()
	mockRepo.EXPECT().
		Create(mock.Anything, mock.AnythingOfType("*domain.User")).Times(0)

	err := service.SignUp(context.Background(), req)
	assert.Contains(t, err.Error(), "user with login alice already exists")
}
