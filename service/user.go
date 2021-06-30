package service

import (
	"context"

	"gitlab.com/InfoBlogFriends/server/request"
	"gitlab.com/InfoBlogFriends/server/validators"

	"golang.org/x/crypto/bcrypt"

	infoblog "gitlab.com/InfoBlogFriends/server"
)

type User struct {
	repository infoblog.UserRepository
}

func NewUserService(repository infoblog.UserRepository) *User {
	return &User{repository: repository}
}

func (u *User) CheckEmailPass(ctx context.Context, email, password string) bool {
	user, err := u.repository.FindByEmail(ctx, email)
	if err != nil {
		return false
	}

	return checkPasswordHash(password, user.Password)
}

func (u *User) CreateByEmailPassword(ctx context.Context, email, password string) error {
	passHash, err := hashPassword(password)
	if err != nil {
		return err
	}

	return u.repository.CreateUserByEmailPassword(ctx, email, passHash)
}

func (u *User) GetProfile(ctx context.Context, uid int64) (infoblog.User, error) {
	user, err := u.repository.Find(ctx, uid)
	if err != nil {
		return infoblog.User{}, err
	}

	return user, nil
}

func (u *User) UpdateProfile(ctx context.Context, profileUpdateReq request.ProfileUpdateReq, uid int64) (infoblog.User, error) {
	user, err := u.repository.Find(ctx, uid)
	if err != nil {
		return infoblog.User{}, err
	}

	if profileUpdateReq.Email != nil {
		if err = validators.CheckEmailFormat(*profileUpdateReq.Email); err != nil {
			return infoblog.User{}, err
		}
		user.Email = *profileUpdateReq.Email
	}
	if profileUpdateReq.Phone != nil {
		*profileUpdateReq.Phone, err = validators.CheckPhoneFormat(*profileUpdateReq.Phone)
		if err != nil {
			return infoblog.User{}, err
		}
		user.Phone = *profileUpdateReq.Phone
	}
	if profileUpdateReq.Name != nil {
		user.Name = *profileUpdateReq.Name
	}
	if profileUpdateReq.SecondName != nil {
		user.SecondName = *profileUpdateReq.SecondName
	}

	return user, u.repository.Update(ctx, user)
}

func (u *User) SetPassword(ctx context.Context, user infoblog.User) error {
	passHash, err := hashPassword(user.Password)
	if err != nil {
		return err
	}
	user.Password = passHash

	return u.repository.SetPassword(ctx, user)
}

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)

	return string(bytes), err
}

func checkPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))

	return err == nil
}
