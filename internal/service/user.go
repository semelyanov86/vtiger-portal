package service

import (
	"context"
	"errors"
	"github.com/semelyanov86/vtiger-portal/internal/config"
	"github.com/semelyanov86/vtiger-portal/internal/domain"
	"github.com/semelyanov86/vtiger-portal/internal/repository"
	"github.com/semelyanov86/vtiger-portal/pkg/cache"
	"github.com/semelyanov86/vtiger-portal/pkg/e"
	"github.com/semelyanov86/vtiger-portal/pkg/logger"
	"strconv"
	"sync"
	"time"
)

const CacheUsersTTL = 50000

var ErrUserNotFound = errors.New("user not found")

var ErrUserIsNotActive = errors.New("user is not active")

type UserSignUpInput struct {
	Email    string `json:"email" binding:"required,email,max=64"`
	Code     string `json:"code" binding:"required,min=3,max=10"`
	Password string `json:"password" binding:"required,min=5,max=20"`
}

type UserSignInInput struct {
	Email    string `json:"email" binding:"required,email,max=64"`
	Password string `json:"password" binding:"required,min=5,max=20"`
}

type PasswordResetInput struct {
	Token    string `json:"token" binding:"required,min=5,max=40"`
	Password string `json:"password" binding:"required,min=5,max=20"`
}

type UsersService struct {
	repo            repository.Users
	crm             repository.UsersCrm
	wg              *sync.WaitGroup
	email           EmailServiceInterface
	company         Company
	tokenRepository repository.Tokens
	document        repository.Document
	cache           cache.Cache
}

func NewUsersService(repo repository.Users, crm repository.UsersCrm, wg *sync.WaitGroup, email EmailServiceInterface, company Company, tokenRepository repository.Tokens, document repository.Document, cache cache.Cache) UsersService {
	return UsersService{repo: repo, crm: crm, wg: wg, email: email, company: company, tokenRepository: tokenRepository, document: document, cache: cache}
}

func (s UsersService) SignUp(ctx context.Context, input UserSignUpInput, cfg *config.Config) (*domain.User, error) {
	_, err := s.repo.GetByEmail(ctx, input.Email)
	if !errors.Is(repository.ErrRecordNotFound, err) {
		return nil, repository.ErrDuplicateEmail
	}
	users, err := s.crm.FindByEmail(ctx, input.Email)
	var user *domain.User
	if err != nil {
		return user, e.Wrap("can not find current email in crm", err)
	}

	for _, u := range users {
		if u.Code == input.Code {
			user = &u
		}
	}
	err = FillVtigerContactWithAdditionalValues(user, input.Password)
	if err != nil {
		return nil, e.Wrap("can not fill data with additional values", err)
	}
	if user == nil || user.Crmid == "" {
		return user, e.Wrap("can not find user in vtiger", ErrUserNotFound)
	}

	if err := s.repo.Insert(ctx, user); err != nil {
		return user, err
	}

	if cfg.Vtiger.Business.ClearCode {
		s.wg.Add(1)
		go func() {
			defer s.wg.Done()
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			_, err := s.crm.ClearUserCodeField(ctx, user.Crmid)
			if err != nil {
				logger.Error(logger.GenerateErrorMessageFromString(err.Error()))
			}
		}()
	}

	if cfg.Email.SendWelcomeEmail {
		companyData, err := s.company.GetCompany(ctx)
		if err != nil {
			return nil, e.Wrap("can not send email because company data not received", err)
		}
		s.wg.Add(1)
		go func() {
			defer s.wg.Done()
			err = s.email.SendGreetingsToUser(VerificationEmailInput{
				Name:         user.FirstName + " " + user.LastName,
				CompanyName:  companyData.OrganizationName,
				SupportEmail: cfg.Vtiger.Business.SupportEmail,
				Email:        user.Email,
			})
			if err != nil {
				logger.Error(logger.GenerateErrorMessageFromString(err.Error()))
			}
		}()
	}

	return user, nil
}

func (s UsersService) GetUserByToken(ctx context.Context, token string) (*domain.User, error) {
	return s.repo.GetForToken(ctx, domain.ScopeAuthentication, token)
}

func (s UsersService) GetUserById(ctx context.Context, id int64) (*domain.User, error) {
	user, err := s.repo.GetById(ctx, id)
	if err != nil {
		return nil, e.Wrap("can not get user by id "+strconv.Itoa(int(id)), err)
	}
	if user.Imageattachmentids != "" {
		file, err := s.document.RetrieveFile(ctx, user.Imageattachmentids)
		if err == nil && file.Filecontents != "" {
			user.Imagecontent = file.Filecontents
		}
	}
	s.wg.Add(1)
	go func() {
		defer s.wg.Done()
		ctx2, cancel := context.WithTimeout(context.Background(), time.Second*5)
		defer cancel()
		updatedUser, err := s.crm.RetrieveById(ctx2, user.Crmid)
		updatedUser.Id = id
		updatedUser.IsActive = user.IsActive
		updatedUser.Password = user.Password
		if err != nil {
			logger.Error(logger.GenerateErrorMessageFromString(err.Error()))
			return
		}
		err = s.repo.Update(ctx2, &updatedUser)
		if err != nil {
			logger.Error(logger.GenerateErrorMessageFromString(err.Error()))
			return
		}
	}()
	return &user, nil
}

func (s UsersService) FindContactsFromAccount(ctx context.Context, filter repository.PaginationQueryFilter) ([]domain.User, int, error) {
	users := make([]domain.User, 0)
	err := GetFromCache[*[]domain.User]("account-"+filter.Client, &users, s.cache)
	if err == nil {
		return users, len(users), nil
	}
	if errors.Is(cache.ErrItemNotFound, err) {
		ids, err := s.crm.FindContactsInAccount(ctx, filter)
		if err != nil {
			return users, 0, err
		}
		for _, id := range ids {
			user, err := s.crm.RetrieveById(ctx, id)
			if err != nil {
				return users, len(users), e.Wrap("can not retrieve user with id "+id, err)
			}
			if user.Imageattachmentids != "" {
				file, err := s.document.RetrieveFile(ctx, user.Imageattachmentids)
				if err == nil && file.Filecontents != "" {
					user.Imagecontent = file.Filecontents
				}
			}
			users = append(users, user)
		}
		err = StoreInCache[*[]domain.User]("account-"+filter.Client, &users, CacheUsersTTL, s.cache)
		if err != nil {
			return users, len(users), err
		}
		return users, len(users), err
	} else {
		return users, len(users), e.Wrap("can not convert caches data to user", err)
	}
}

func (s UsersService) ResetUserPassword(ctx context.Context, input PasswordResetInput) (domain.User, error) {
	user, err := s.repo.GetForToken(ctx, domain.ScopePasswordReset, input.Token)
	if err != nil {
		return domain.User{}, err
	}
	err = user.Password.Set(input.Password)
	if err != nil {
		return *user, err
	}
	err = s.repo.Update(ctx, user)
	if err != nil {
		return *user, err
	}
	err = s.tokenRepository.DeleteAllForUser(ctx, domain.ScopePasswordReset, user.Id)
	if err != nil {
		return *user, err
	}
	return *user, nil
}

func FillVtigerContactWithAdditionalValues(user *domain.User, password string) error {
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()
	user.IsActive = true
	err := user.Password.Set(password)
	if err != nil {
		return e.Wrap("can not hash password", err)
	}
	return nil
}
