package services

import (
	"context"
	"errors"
	"fmt"

	"gitlab.com/InfoBlogFriends/server/decoder"

	"gitlab.com/InfoBlogFriends/server/hasher"

	infoblog "gitlab.com/InfoBlogFriends/server"
	"gitlab.com/InfoBlogFriends/server/request"
)

type User struct {
	*decoder.Decoder
	userRepository  infoblog.UserRepository
	subsRepository  infoblog.SubscriberRepository
	likesRepository infoblog.LikeRepository
	post            *Post
}

func NewUserService(rs infoblog.Repositories, subs infoblog.SubscriberRepository, post *Post) *User {
	return &User{userRepository: rs.Users, subsRepository: subs, Decoder: decoder.NewDecoder(), post: post, likesRepository: rs.Likes}
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

func (u *User) GetProfile(ctx context.Context) (request.UserData, error) {
	user, err := extractUser(ctx)
	if err != nil {
		return request.UserData{}, err
	}

	user, err = u.userRepository.Find(ctx, user)
	if err != nil {
		return request.UserData{}, err
	}
	userData := request.UserData{}
	err = u.MapStructs(&userData, &user)
	if err != nil {
		return request.UserData{}, err
	}

	postCount, err := u.post.CountByUser(ctx, user)
	subsCount, err := u.subsRepository.CountByUser(ctx, user)
	likesCount, err := u.likesRepository.CountByUser(ctx, user)

	if err != nil {
		return request.UserData{}, err
	}

	userData.Counts = &request.UserDataCounts{
		Posts:       postCount,
		Subscribers: subsCount,
		Friends:     377,
		Likes:       likesCount,
	}

	userData.PasswordSet = &user.Password.Valid

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

func (u *User) SetPassword(ctx context.Context, setPasswordReq request.SetPasswordReq) error {
	user, err := extractUser(ctx)
	if err != nil {
		return err
	}
	user, err = u.userRepository.Find(ctx, user)
	if err != nil {
		return err
	}
	if user.Password.Valid {
		if setPasswordReq.OldPassword == nil {
			return fmt.Errorf("old password is required")
		}
		if !hasher.CheckPasswordHash(*setPasswordReq.OldPassword, user.Password.String) {
			return fmt.Errorf("wrong old password")
		}
	}

	passHash, err := hasher.HashPassword(setPasswordReq.Password)
	if err != nil {
		return err
	}
	user.Password = infoblog.NewNullString(passHash)

	return u.userRepository.SetPassword(ctx, user)
}

func (u *User) Subscribe(ctx context.Context, user infoblog.User, subscribeRequest request.UserIDRequest) error {
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

func (u *User) Unsubscribe(ctx context.Context, user infoblog.User, subscribeRequest request.UserIDRequest) error {
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

func (u *User) List(ctx context.Context) ([]request.UserData, error) {
	users, err := u.userRepository.FindAll(ctx)
	if err != nil {
		return nil, err
	}

	usersData := []request.UserData{}
	for _, user := range users {
		userData := request.UserData{}
		err = u.MapStructs(&userData, &user)
		if err != nil {
			return nil, err
		}
		usersData = append(usersData, userData)
	}

	return usersData, nil
}

func (u *User) Get(ctx context.Context, req request.UserIDNickRequest) (request.UserData, error) {
	user := infoblog.User{}
	var err error
	if req.UUID != nil {
		user.UUID = *req.UUID
		user, err = u.userRepository.Find(ctx, user)
		if err != nil {
			return request.UserData{}, err
		}
	}
	if req.NickName != nil {
		user.NickName = infoblog.NewNullString(*req.NickName)
		user, err = u.userRepository.FindNickname(ctx, user)
		if err != nil {
			return request.UserData{}, err
		}
	}

	userData := request.UserData{}
	err = u.MapStructs(&userData, &user)
	if err != nil {
		return request.UserData{}, err
	}

	return userData, nil
}

func extractUser(ctx context.Context) (infoblog.User, error) {
	u, ok := ctx.Value("user").(*infoblog.User)
	if !ok {
		return infoblog.User{}, errors.New("type assertion to user err")
	}

	if u.ID == 0 {
		return infoblog.User{}, errors.New("user not exists")
	}

	return *u, nil
}
