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
	"github.com/semelyanov86/vtiger-portal/pkg/vtiger"
	"strconv"
	"sync"
	"time"
)

const CacheUsersTTL = 50000

var ErrUserNotFound = errors.New("user not found")

var ErrUserIsNotActive = errors.New("user is not active")

var ErrOldPasswordIsWrong = errors.New("old password is wrong")

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

type UserUpdateInput struct {
	AccountName     string `json:"account_name"`
	Department      string `json:"department" binding:"omitempty,min=2,max=50"`
	Description     string `json:"description" binding:"omitempty,min=5,max=200"`
	Email           string `json:"email" binding:"omitempty,email"`
	FirstName       string `json:"firstname" binding:"omitempty,min=2,max=50"`
	LastName        string `json:"lastname" binding:"omitempty,min=2,max=50"`
	Mailingcity     string `json:"mailingcity" binding:"omitempty,min=2,max=50"`
	Mailingcountry  string `json:"mailingcountry" binding:"omitempty,min=2,max=50"`
	Mailingpobox    string `json:"mailingpobox"`
	Mailingstate    string `json:"mailingstate"`
	Mailingstreet   string `json:"mailingstreet" binding:"omitempty,min=2,max=100"`
	Mailingzip      string `json:"mailingzip" binding:"omitempty,min=2,max=10"`
	Othercity       string `json:"othercity"`
	Othercountry    string `json:"othercountry"`
	Otherpobox      string `json:"otherpobox"`
	Otherstate      string `json:"otherstate"`
	Otherstreet     string `json:"otherstreet"`
	Otherzip        string `json:"otherzip"`
	Password        string `json:"password"`
	OldPassword     string `json:"old_password"`
	ConfirmPassword string `json:"confirmPassword"`
	Phone           string `json:"phone" binding:"omitempty,min=5,max=15"`
	Title           string `json:"title" binding:"omitempty,min=2,max=50"`
	Imagename       string `json:"imagename"`
	Imagetype       string `json:"imagetype"`
	Imagecontent    string `json:"imagecontent"`
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
	account         AccountService
	config          config.Config
}

func NewUsersService(repo repository.Users, crm repository.UsersCrm, wg *sync.WaitGroup, email EmailServiceInterface, company Company, tokenRepository repository.Tokens, document repository.Document, cache cache.Cache, account AccountService, config config.Config) UsersService {
	return UsersService{repo: repo, crm: crm, wg: wg, email: email, company: company, tokenRepository: tokenRepository, document: document, cache: cache, account: account, config: config}
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
	if user.AccountId != "" {
		account, err := s.account.GetAccountById(ctx, user.AccountId)
		if err != nil {
			return user, e.Wrap("can not get account info", err)
		}
		user.AccountName = account.AccountName
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

func (s UsersService) Update(ctx context.Context, id int64, updateData UserUpdateInput) (*domain.User, error) {
	user, err := s.GetUserById(ctx, id)
	if err != nil {
		return nil, e.Wrap("can not get user by id "+strconv.Itoa(int(id)), err)
	}
	if updateData.FirstName != "" {
		user.FirstName = updateData.FirstName
	}
	if updateData.LastName != "" {
		user.LastName = updateData.LastName
	}
	if updateData.Email != "" {
		user.Email = updateData.Email
	}
	if updateData.AccountName != "" {
		user.AccountName = updateData.AccountName
	}
	if updateData.Department != "" {
		user.Department = updateData.Department
	}
	if updateData.Description != "" {
		user.Description = updateData.Description
	}
	if updateData.Email != "" {
		user.Email = updateData.Email
	}
	if updateData.Mailingcity != "" {
		user.MailingCity = updateData.Mailingcity
	}
	if updateData.Mailingpobox != "" {
		user.MailingPoBox = updateData.Mailingpobox
	}
	if updateData.Mailingstate != "" {
		user.MailingState = updateData.Mailingstate
	}
	if updateData.Mailingstreet != "" {
		user.MailingStreet = updateData.Mailingstreet
	}
	if updateData.Mailingzip != "" {
		user.MailingZip = updateData.Mailingzip
	}
	if updateData.Othercity != "" {
		user.OtherCity = updateData.Othercity
	}
	if updateData.Othercountry != "" {
		user.OtherCountry = updateData.Othercountry
	}
	if updateData.Otherpobox != "" {
		user.OtherPoBox = updateData.Otherpobox
	}
	if updateData.Otherstate != "" {
		user.OtherState = updateData.Otherstate
	}
	if updateData.Otherstreet != "" {
		user.OtherStreet = updateData.Otherstreet
	}
	if updateData.Otherzip != "" {
		user.OtherZip = updateData.Otherzip
	}
	if updateData.Phone != "" {
		user.Phone = updateData.Phone
	}
	if updateData.Title != "" {
		user.Title = updateData.Title
	}
	if updateData.Imagename != "" {
		user.Imagename = updateData.Imagename
	}
	if updateData.Imagetype != "" {
		user.Imagetype = updateData.Imagetype
	}
	if updateData.Imagecontent != "" {
		user.Imagecontent = updateData.Imagecontent
	}
	if updateData.Password != "" {
		if updateData.Password != updateData.ConfirmPassword {
			return user, ErrPasswordDoesNotMatch
		}
		if !user.Password.Matches(updateData.OldPassword) {
			return user, ErrOldPasswordIsWrong
		}
		user.Password.Set(updateData.Password)
	}
	err = s.repo.Update(ctx, user)
	if err != nil {
		return user, e.Wrap("can not update user", err)
	}
	_, err = s.crm.Update(ctx, user.Crmid, *user)
	return user, err
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
		updatedUser.Otp_verified = user.Otp_verified
		updatedUser.Otp_enabled = user.Otp_enabled
		updatedUser.Otp_auth_url = user.Otp_auth_url
		updatedUser.Otp_secret = user.Otp_secret
		if err != nil {
			logger.Error(logger.GenerateErrorMessageFromString(err.Error()))
			return
		}
		if updatedUser.AccountId != "" {
			account, err := s.account.GetAccountById(ctx2, updatedUser.AccountId)
			if err != nil {
				logger.Error(logger.GenerateErrorMessageFromString(err.Error()))
				return
			}
			updatedUser.AccountName = account.AccountName
		}
		err = s.repo.Update(ctx2, &updatedUser)
		if err != nil {
			logger.Error(logger.GenerateErrorMessageFromString(err.Error()))
			return
		}
	}()
	return &user, nil
}

func (s UsersService) FindContactsFromAccount(ctx context.Context, filter vtiger.PaginationQueryFilter) ([]domain.User, int, error) {
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

func (s UsersService) FindByCrmid(ctx context.Context, id string) (domain.User, error) {
	var user domain.User
	err := GetFromCache[*domain.User](id, &user, s.cache)
	if err == nil {
		return user, nil
	}
	if errors.Is(cache.ErrItemNotFound, err) {
		user, err := s.crm.RetrieveById(ctx, id)
		if err != nil {
			return user, err
		}
		if user.Imageattachmentids != "" {
			file, err := s.document.RetrieveFile(ctx, user.Imageattachmentids)
			if err == nil && file.Filecontents != "" {
				user.Imagecontent = file.Filecontents
			}
		}
		if user.AccountId != "" {
			account, err := s.account.GetAccountById(ctx, user.AccountId)
			if err != nil {
				return user, e.Wrap("can not retrieve account", err)
			}
			user.AccountName = account.AccountName
		}

		err = StoreInCache[*domain.User](id, &user, CacheUsersTTL, s.cache)
		if err != nil {
			return user, err
		}
		return user, err
	} else {
		return user, e.Wrap("can not convert caches data to user", err)
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

func (s UsersService) GetUserSettings(ctx context.Context, id string) (map[string]bool, error) {
	result := make(map[string]bool)
	contact, err := s.crm.RetrieveContactMap(ctx, id)
	if err != nil {
		return result, e.Wrap("can not get user info", err)
	}
	settingsFields := s.config.Vtiger.Business.UserSettingsFields

	for _, field := range settingsFields {
		if contact[field] == "1" || contact[field] == 1 {
			result[field] = true
		} else {
			result[field] = false
		}
	}
	return result, nil
}

func (s UsersService) ChangeUserSetting(ctx context.Context, id string, field string, value bool) error {
	return s.crm.ChangeSettingField(ctx, id, field, value)
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
