package services

import (
	"context"
	"errors"

	"gitlab.com/InfoBlogFriends/server/decoder"

	"gitlab.com/InfoBlogFriends/server/hasher"

	infoblog "gitlab.com/InfoBlogFriends/server"
	"gitlab.com/InfoBlogFriends/server/request"
)

type User struct {
	*decoder.Decoder
	userRepository infoblog.UserRepository
	subsRepository infoblog.SubscriberRepository
}

func NewUserService(repository infoblog.UserRepository, subs infoblog.SubscriberRepository) *User {
	return &User{userRepository: repository, subsRepository: subs, Decoder: decoder.NewDecoder()}
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
	return u.userRepository.CreateUser(ctx, user)
}

func (u *User) GetProfile(ctx context.Context, user infoblog.User) (request.UserData, error) {
	user, err := u.userRepository.Find(ctx, user)
	if err != nil {
		return request.UserData{}, err
	}
	userData := request.UserData{}
	err = u.MapStructs(&userData, &user)
	if err != nil {
		return request.UserData{}, err
	}

	return userData, nil
}

func (u *User) UpdateProfile(ctx context.Context, profileUpdateReq request.ProfileUpdateReq, user infoblog.User) (request.UserData, error) {
	user, err := u.userRepository.Find(ctx, user)
	if err != nil {
		return request.UserData{}, err
	}

	err = u.MapStructs(&user, &profileUpdateReq)
	if err != nil {
		return request.UserData{}, err
	}
	user.Active = infoblog.NewNullBool(true)
	err = u.userRepository.Update(ctx, user)
	if err != nil {
		return request.UserData{}, err
	}

	userData := request.UserData{}
	err = u.MapStructs(&userData, &user)
	if err != nil {
		return request.UserData{}, err
	}

	return userData, nil
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

func (u *User) List(ctx context.Context, user infoblog.User, subscribeRequest request.UserSubscriberRequest) error {
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
