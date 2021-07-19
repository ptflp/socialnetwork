package services

import (
	"context"
	"errors"

	"gitlab.com/InfoBlogFriends/server/hasher"

	"gitlab.com/InfoBlogFriends/server/request"
	"gitlab.com/InfoBlogFriends/server/validators"

	infoblog "gitlab.com/InfoBlogFriends/server"
)

type User struct {
	userRepository infoblog.UserRepository
	subsRepository infoblog.SubscriberRepository
}

func NewUserService(repository infoblog.UserRepository, subs infoblog.SubscriberRepository) *User {
	return &User{userRepository: repository, subsRepository: subs}
}

func (u *User) CheckEmailPass(ctx context.Context, user infoblog.User) bool {
	uDB, err := u.userRepository.FindByEmail(ctx, user)
	if err != nil {
		return false
	}

	return hasher.CheckPasswordHash(user.Password.String, uDB.Password.String)
}

func (u *User) CreateByEmailPassword(ctx context.Context, user infoblog.User) error {
	passHash, err := hasher.HashPassword(user.Password.String)
	if err != nil {
		return err
	}

	user.Password = infoblog.NewNullString(passHash)
	return u.userRepository.CreateUserByEmailPassword(ctx, user)
}

func (u *User) GetProfile(ctx context.Context, user infoblog.User) (infoblog.User, error) {
	user, err := u.userRepository.Find(ctx, user)
	if err != nil {
		return infoblog.User{}, err
	}

	return user, nil
}

func (u *User) UpdateProfile(ctx context.Context, profileUpdateReq request.ProfileUpdateReq, user infoblog.User) (infoblog.User, error) {
	user, err := u.userRepository.Find(ctx, user)
	if err != nil {
		return infoblog.User{}, err
	}

	if profileUpdateReq.Email != nil {
		if err = validators.CheckEmailFormat(*profileUpdateReq.Email); err != nil {
			return infoblog.User{}, err
		}
		user.Email = infoblog.NewNullString(*profileUpdateReq.Email)
	}
	if profileUpdateReq.Phone != nil {
		*profileUpdateReq.Phone, err = validators.CheckPhoneFormat(*profileUpdateReq.Phone)
		if err != nil {
			return infoblog.User{}, err
		}
		user.Phone = infoblog.NewNullString(*profileUpdateReq.Phone)
	}
	if profileUpdateReq.Name != nil {
		user.Name = infoblog.NewNullString(*profileUpdateReq.Name)
	}
	if profileUpdateReq.SecondName != nil {
		user.SecondName = infoblog.NewNullString(*profileUpdateReq.SecondName)
	}

	return user, u.userRepository.Update(ctx, user)
}

func (u *User) SetPassword(ctx context.Context, user infoblog.User) error {
	passHash, err := hasher.HashPassword(user.Password.String)
	if err != nil {
		return err
	}
	user.Password = infoblog.NewNullString(passHash)

	return u.userRepository.SetPassword(ctx, user)
}

func (u *User) Subscribe(ctx context.Context, user infoblog.User, subscribeRequest request.UserSubscriberRequest) error {
	sub, err := u.userRepository.Find(ctx, infoblog.User{UUID: subscribeRequest.UUID})
	if err != nil {
		return err
	}
	if sub.ID < 1 {
		return errors.New("user with specified id not found")
	}

	_, err = u.subsRepository.Create(ctx, infoblog.Subscriber{
		UserUUID:       user.UUID,
		SubscriberUUID: subscribeRequest.UUID,
		Active:         infoblog.NewNullBool(true),
	})

	return err
}

func (u *User) Unsubscribe(ctx context.Context, user infoblog.User, subscribeRequest request.UserSubscriberRequest) error {
	sub, err := u.userRepository.Find(ctx, infoblog.User{UUID: subscribeRequest.UUID})
	if err != nil {
		return err
	}
	if sub.ID < 1 {
		return errors.New("user with specified id not found")
	}

	err = u.subsRepository.Delete(ctx, infoblog.Subscriber{
		UserUUID:       user.UUID,
		SubscriberUUID: subscribeRequest.UUID,
		Active:         infoblog.NewNullBool(false),
	})

	return err
}
