package services

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"html/template"
	"math/rand"
	"net/url"
	"path"
	"time"

	sq "github.com/Masterminds/squirrel"

	"gitlab.com/InfoBlogFriends/server/email"
	"gitlab.com/InfoBlogFriends/server/types"

	"gitlab.com/InfoBlogFriends/server/utils"

	"gitlab.com/InfoBlogFriends/server/components"
	"go.uber.org/zap"

	"gitlab.com/InfoBlogFriends/server/validators"

	"gitlab.com/InfoBlogFriends/server/decoder"

	"gitlab.com/InfoBlogFriends/server/hasher"

	infoblog "gitlab.com/InfoBlogFriends/server"
	"gitlab.com/InfoBlogFriends/server/request"
)

const (
	PhoneRecoverKey      = "phone:recover:%s"
	PasswordRecoveryUUID = "R"
	RecoveryIDKey        = "recover:id:%s"
)

type User struct {
	*decoder.Decoder
	userRepository  infoblog.UserRepository
	subsRepository  infoblog.SubscriberRepository
	likesRepository infoblog.LikeRepository
	post            *Post
	file            *File
	components.Componenter
}

func NewUserService(rs infoblog.Repositories, post *Post, cmps components.Componenter, file *File) *User {
	return &User{userRepository: rs.Users, subsRepository: rs.Subscribers, Decoder: decoder.NewDecoder(), post: post, likesRepository: rs.Likes, Componenter: cmps, file: file}
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

	user.Password = types.NewNullString(passHash)
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
	cond := infoblog.Condition{
		Equal: &sq.Eq{"user_uuid": user.UUID},
	}
	subscribers, err := u.subsRepository.Listx(ctx, cond)
	if err != nil {
		return request.UserData{}, err
	}

	user.Subscribers.Uint64.Uint64 = uint64(len(subscribers))
	user.Subscribers.Uint64.Valid = true

	cond = infoblog.Condition{
		Equal: &sq.Eq{"subscriber_uuid": user.UUID},
	}

	subscribes, err := u.subsRepository.Listx(ctx, cond)
	if err != nil {
		return request.UserData{}, err
	}

	user.Subscribes.Uint64.Uint64 = uint64(len(subscribes))
	user.Subscribes.Uint64.Valid = true

	userData := request.UserData{}
	err = u.MapStructs(&userData, &user)
	if err != nil {
		return request.UserData{}, err
	}

	postCount, err := u.post.CountByUser(ctx, user)
	if err != nil {
		return request.UserData{}, err
	}
	subsCount, err := u.subsRepository.CountByUser(ctx, user)
	if err != nil {
		return request.UserData{}, err
	}
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

	user.Active = types.NewNullBool(true)
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
	user.Password = types.NewNullString(passHash)

	return u.userRepository.SetPassword(ctx, user)
}

func (u *User) prepareRecoveryTemplate(recoverUrl string) (bytes.Buffer, error) {
	tmpl, err := template.ParseFiles("./templates/password_recovery.html")
	if err != nil {
		u.Logger().Error("recover password template parse", zap.Error(err))
		return bytes.Buffer{}, err
	}
	type PasswordRecover struct {
		PasswordRecover string
	}

	b := bytes.Buffer{}

	err = tmpl.Execute(&b, PasswordRecover{PasswordRecover: recoverUrl})

	return b, err
}

func (u *User) Subscribe(ctx context.Context, subscribeRequest request.UserIDRequest) error {
	user, err := extractUser(ctx)
	if err != nil {
		return err
	}
	sub, err := u.userRepository.Find(ctx, infoblog.User{UUID: types.NewNullUUID(subscribeRequest.UUID)})
	if err != nil {
		return err
	}
	if !sub.UUID.Valid {
		return errors.New("user with specified id not found")
	}

	_, err = u.subsRepository.Create(ctx, infoblog.Subscriber{
		UserUUID:       sub.UUID,
		SubscriberUUID: user.UUID,
		Active:         types.NewNullBool(true),
	})
	if err != nil {
		return err
	}

	user, err = u.userRepository.Count(ctx, user, "subscribes", "incr")
	if err != nil {
		return err
	}
	sub, err = u.userRepository.Count(ctx, sub, "subscribers", "incr")
	if err != nil {
		return err
	}

	return err
}

func (u *User) Unsubscribe(ctx context.Context, user infoblog.User, subscribeRequest request.UserIDRequest) error {
	sub, err := u.userRepository.Find(ctx, infoblog.User{UUID: types.NewNullUUID(subscribeRequest.UUID)})
	if err != nil {
		return err
	}
	if !sub.UUID.Valid {
		return errors.New("user with specified id not found")
	}

	err = u.subsRepository.Delete(ctx, infoblog.Subscriber{
		UserUUID:       types.NewNullUUID(subscribeRequest.UUID),
		SubscriberUUID: user.UUID,
		Active:         types.NewNullBool(false),
	})
	if err != nil {
		return err
	}

	user, err = u.userRepository.Count(ctx, user, "subscribes", "decr")
	if err != nil {
		return err
	}

	sub, err = u.userRepository.Count(ctx, sub, "subscribers", "decr")
	if err != nil {
		return err
	}

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

func (u *User) Listx(ctx context.Context, condition infoblog.Condition) ([]infoblog.User, error) {
	return u.userRepository.Listx(ctx, condition)
}

func (u *User) TempList(ctx context.Context, req request.LimitOffsetReq) ([]request.UserData, error) {
	users, err := u.userRepository.FindLimitOffset(ctx, uint64(req.Limit), uint64(req.Offset))
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

func (u *User) Recommends(ctx context.Context, req request.LimitOffsetReq) ([]request.UserData, error) {
	user, err := extractUser(ctx)
	if err != nil {
		return nil, err
	}
	condition := infoblog.Condition{
		NotIn: &infoblog.In{
			Field: "uuid",
			Args:  []interface{}{user.UUID},
		},
		Order: &infoblog.Order{
			Field: "likes",
		},
		Other: &infoblog.Other{
			Condition: "nickname IS NOT null",
			Args:      nil,
		},
		LimitOffset: &infoblog.LimitOffset{
			Offset: req.Offset,
			Limit:  req.Limit,
		},
	}
	users, err := u.userRepository.Listx(ctx, condition)
	if err != nil {
		return nil, err
	}

	var usersData []request.UserData

	err = u.MapStructs(&usersData, &users)
	if err != nil {
		return nil, err
	}

	return usersData, nil
}

func (u *User) GetUsersByCondition(ctx context.Context, condition infoblog.Condition) ([]request.UserData, error) {
	users, err := u.userRepository.Listx(ctx, condition)
	if err != nil {
		return nil, err
	}

	var usersData []request.UserData

	err = u.MapStructs(&usersData, &users)
	if err != nil {
		return nil, err
	}

	return usersData, nil
}

func (u *User) GetUserSubscribesUUIDs(ctx context.Context, user infoblog.User, subscribesCondition infoblog.Condition) ([]interface{}, error) {
	mySubs, err := u.subsRepository.Listx(ctx, subscribesCondition)
	if err != nil {
		return nil, err
	}
	var userUUIDs []interface{}

	for i := range mySubs {
		userUUIDs = append(userUUIDs, mySubs[i].UserUUID)
	}

	return userUUIDs, nil
}

func (u *User) Subscribes(ctx context.Context, req request.LimitOffsetReq) ([]request.UserData, error) {
	user, err := extractUser(ctx)
	if err != nil {
		return nil, err
	}
	subscribesCondition := infoblog.Condition{
		Equal: &sq.Eq{"active": true, "subscriber_uuid": user.UUID},
		LimitOffset: &infoblog.LimitOffset{
			Limit:  req.Limit,
			Offset: req.Offset,
		},
	}
	userUUIDs, err := u.GetUserSubscribesUUIDs(ctx, user, subscribesCondition)

	if err != nil {
		return nil, err
	}

	if len(userUUIDs) < 1 {
		return []request.UserData{}, nil
	}

	condition := infoblog.Condition{
		In: &infoblog.In{
			Field: "uuid",
			Args:  userUUIDs,
		},
	}

	return u.GetUsersByCondition(ctx, condition)
}

func (u *User) Get(ctx context.Context, req request.UserIDNickRequest) (request.UserData, error) {
	user := infoblog.User{}
	var err error
	if req.UUID != nil {
		user.UUID = types.NewNullUUID(*req.UUID)
		user, err = u.userRepository.Find(ctx, user)
		if err != nil {
			return request.UserData{}, err
		}
	}
	if req.NickName != nil {
		user.NickName = types.NewNullString(*req.NickName)
		user, err = u.userRepository.FindByNickname(ctx, user)
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

func (u *User) Autocomplete(ctx context.Context, req request.UserNicknameRequest) ([]request.UserData, error) {
	users, err := u.userRepository.FindLikeNickname(ctx, req.Nickname)
	if err != nil {
		return nil, err
	}
	usersData := make([]request.UserData, 0, len(users))
	for i := range users {
		var userData request.UserData
		err = u.MapStructs(&userData, users[i])
		if err != nil {
			return nil, err
		}
		usersData = append(usersData, userData)
	}

	return usersData, nil
}

func (u *User) PasswordRecover(ctx context.Context, req request.PasswordRecoverRequest) error {
	user := infoblog.User{}

	err := u.MapStructs(&user, &req)
	if err != nil {
		return err
	}

	if user.Email.Valid {
		err = validators.CheckEmailFormat(user.Email.String)
		if err != nil {
			return err
		}
		user, err = u.userRepository.FindByEmail(ctx, user)
		if err != nil {
			return err
		}
		// send email
		var recoverUrl string
		recoverUrl, _, err = u.generateRecoverUrl(user)
		if err != nil {
			return err
		}

		var body bytes.Buffer
		body, err = u.prepareRecoveryTemplate(recoverUrl)
		if err != nil {
			return err
		}

		msg := email.NewMessage()
		msg.SetSubject("Восстановление пароля")
		msg.SetType(email.TypeHtml)
		msg.SetReceiver(user.Email.String)
		msg.SetBody(body)

		err = u.Email().Send(msg)
		if err != nil {
			return err
		}
	}

	if user.Phone.Valid {
		user.Phone.String, err = validators.CheckPhoneFormat(user.Phone.String)
		if err != nil {
			return err
		}
		user, err = u.userRepository.FindByPhone(ctx, user)
		if err != nil {
			return err
		}

		code := genCode()
		if u.Config().SMSC.Dev {
			code = 3455
		}
		u.Cache().Set(fmt.Sprintf(PhoneRecoverKey, user.Phone.String), &code, 15*time.Minute)
		if u.Config().SMSC.Dev {
			return nil
		}

		err = u.Componenter.SMS().Send(ctx, user.Phone.String, fmt.Sprintf("Ваш код: %d", code))
		if err != nil {
			u.Logger().Error("send sms err", zap.String("user.Phone.String", user.Phone.String), zap.Int("code", code))
		}

		return err
	}

	return errors.New("bad request params")
}

func (u *User) CheckPhoneCode(ctx context.Context, req request.CheckPhoneCodeRequest) (request.RecoverChekPhoneResponse, error) {
	var code int64
	var user infoblog.User
	err := u.Cache().Get(fmt.Sprintf(PhoneRecoverKey, req.Phone), &code)
	if err != nil {
		return request.RecoverChekPhoneResponse{}, err
	}
	if code != req.Code {
		return request.RecoverChekPhoneResponse{}, errors.New("user code error")
	}
	user.Phone = types.NewNullString(req.Phone)
	user, err = u.userRepository.FindByPhone(ctx, user)
	if err != nil {
		return request.RecoverChekPhoneResponse{}, err
	}

	recoverID, err := u.GenerateRecoverID(user)
	if err != nil {
		return request.RecoverChekPhoneResponse{}, err
	}

	return request.RecoverChekPhoneResponse{
		Success: true,
		Data: request.RecoverCheckPhoneData{
			RecoverID: recoverID,
		},
	}, nil
}

func (u *User) GenerateRecoverID(user infoblog.User) (string, error) {
	recoverID, err := utils.ProjectUUIDGen(PasswordRecoveryUUID)
	if err != nil {
		return "", err
	}
	u.Cache().Set(fmt.Sprintf(RecoveryIDKey, recoverID), &user.UUID, 15*time.Minute)

	return recoverID, err
}

func (u *User) PasswordReset(ctx context.Context, req request.PasswordResetRequest) error {
	var user infoblog.User
	err := u.Cache().Get(fmt.Sprintf(RecoveryIDKey, req.RecoverID), &user.UUID)
	if err != nil {
		return err
	}
	user, err = u.userRepository.Find(ctx, user)
	if err != nil {
		return err
	}
	passHash, err := hasher.HashPassword(req.Password)
	if err != nil {
		return err
	}
	user.Password = types.NewNullString(passHash)
	err = u.userRepository.SetPassword(ctx, user)
	if err != nil {
		return err
	}

	return err
}

func (u *User) EmailExist(ctx context.Context, req request.EmailRequest) error {
	var user infoblog.User
	user.Email = types.NewNullString(req.Email)
	_, err := u.userRepository.FindByEmail(ctx, user)

	return err
}

func (u *User) NicknameExist(ctx context.Context, req request.NicknameRequest) error {
	var user infoblog.User
	user.NickName = types.NewNullString(req.Nickname)
	_, err := u.userRepository.FindByNickname(ctx, user)

	return err
}

func (u *User) generateRecoverUrl(user infoblog.User) (string, string, error) {

	recoverID, err := u.GenerateRecoverID(user)
	if err != nil {
		return "", "", err
	}

	uri, err := url.Parse(u.Config().App.FrontEnd)
	if err != nil {
		return "", "", err
	}
	uri.Path = fmt.Sprintf("profile/password/%s", recoverID)

	return uri.String(), recoverID, err
}

func (u *User) SaveAvatar(ctx context.Context, formFile FormFile) (request.UserData, error) {
	// 1. save file to filesystem
	user, err := extractUser(ctx)
	if err != nil {
		return request.UserData{}, err
	}

	fileUUID := types.NewNullUUID()

	file, err := u.file.SaveFileSystem(formFile, user, fileUUID)
	if err != nil {
		return request.UserData{}, err
	}

	// 2. save post to db

	// 3. update file info, save to db
	file.Active = 1
	file.Type = types.TypeAvatar
	file.UserUUID = user.UUID

	err = u.file.SaveDB(ctx, &file)
	if err != nil {
		return request.UserData{}, err
	}

	link := "/" + path.Join(file.Dir, file.Name)

	user, err = u.userRepository.Find(ctx, user)
	if err != nil {
		return request.UserData{}, err
	}

	user.Avatar = types.NewNullString(link)

	err = u.userRepository.Update(ctx, user)
	if err != nil {
		return request.UserData{}, err
	}

	userData, err := u.GetUserData(user)
	if err != nil {
		return request.UserData{}, err
	}

	return userData, nil
}

func (u *User) GetUserData(user infoblog.User) (request.UserData, error) {

	var userData request.UserData
	err := u.MapStructs(&userData, &user)

	userData.AvatarSet = user.Avatar.Valid
	return userData, err
}

func (u *User) Count(ctx context.Context, user infoblog.User, field, ops string) (infoblog.User, error) {
	return u.userRepository.Count(ctx, user, field, ops)
}

func extractUser(ctx context.Context) (infoblog.User, error) {
	u, ok := ctx.Value(types.User{}).(*infoblog.User)
	if !ok {
		return infoblog.User{}, errors.New("type assertion to user err")
	}

	return *u, nil
}

func genCode() int {
	rand.Seed(time.Now().UnixNano())
	code := rand.Intn(8999) + 1000

	return code
}
