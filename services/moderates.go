package services

import (
	"context"

	"gitlab.com/InfoBlogFriends/server/utils"

	sq "github.com/Masterminds/squirrel"

	"gitlab.com/InfoBlogFriends/server/decoder"
	"gitlab.com/InfoBlogFriends/server/request"
	"gitlab.com/InfoBlogFriends/server/types"

	infoblog "gitlab.com/InfoBlogFriends/server"
)

type Moderates struct {
	moderateRep infoblog.ModerateRepository
	Services    *Services
	*decoder.Decoder
}

func NewModeratesService(reps infoblog.Repositories, services *Services) *Moderates {
	return &Moderates{moderateRep: reps.Moderate, Decoder: decoder.NewDecoder(), Services: services}
}

func (m *Moderates) SaveFile(ctx context.Context, formFile FormFile) (request.FileData, error) {
	// 1. save file to filesystem
	u, err := extractUser(ctx)
	if err != nil {
		return request.FileData{}, err
	}

	fileUUID := types.NewNullUUID()

	file, err := m.Services.File.SaveFileSystem(formFile, u, fileUUID)
	if err != nil {
		return request.FileData{}, err
	}
	// 3. update file info, save to db
	file.Active = 1
	file.Type = 1
	file.UserUUID = u.UUID

	err = m.Services.File.SaveDB(ctx, &file)
	if err != nil {
		return request.FileData{}, err
	}
	var postFileData request.FileData

	err = m.MapStructs(&postFileData, &file)
	if err != nil {
		return request.FileData{}, err
	}
	postFileData.Link = utils.PublicLink(postFileData)

	return postFileData, nil
}

func (m *Moderates) CreateModerate(ctx context.Context, moderateReq request.ModerateCreateReq) error {
	user, err := extractUser(ctx)
	if err != nil {
		return err
	}

	moderate := infoblog.Moderate{
		UUID:     types.NewNullUUID(),
		Type:     types.TypeUserModerate,
		Active:   types.NewNullBool(true),
		UserUUID: user.UUID,
	}

	err = m.moderateRep.Create(ctx, moderate)
	if err != nil {
		return err
	}

	file := infoblog.File{
		Type:        types.TypeUserModerate,
		ForeignUUID: moderate.UUID,
	}

	err = m.Services.File.UpdateFileType(ctx, moderateReq.FilesID, file)

	return err
}

func (m *Moderates) GetModerates(ctx context.Context, limitOffsetReq request.LimitOffsetReq) ([]request.ModerateData, error) {
	condition := infoblog.Condition{
		Equal: &sq.Eq{"type": types.TypeUserModerate, "active": true},
		LimitOffset: &infoblog.LimitOffset{
			Offset: limitOffsetReq.Offset,
			Limit:  limitOffsetReq.Limit,
		},
	}

	moderates, err := m.moderateRep.Listx(ctx, condition)
	if err != nil {
		return nil, err
	}

	moderatesData := make([]request.ModerateData, 0, len(moderates))
	if len(moderates) < 1 {
		return moderatesData, nil
	}

	err = m.MapStructs(&moderatesData, &moderates)
	if err != nil {
		return nil, err
	}
	uuids := make([]interface{}, 0, len(moderatesData))
	moderatesDataMap := make(map[string]*request.ModerateData, len(moderatesData))
	for i := range moderatesData {
		uuids = append(uuids, moderatesData[i].UUID)
		moderatesDataMap[moderatesData[i].UUID.String] = &moderatesData[i]
	}
	_, err = m.FeedGetFiles(ctx, moderatesDataMap, uuids...)
	if err != nil {
		return nil, err
	}

	return moderatesData, nil
}

func (m *Moderates) Get(ctx context.Context, req request.UUIDReq) (request.ModerateData, error) {
	condition := infoblog.Condition{
		Equal: &sq.Eq{"type": types.TypeUserModerate, "active": true},
		In: &infoblog.In{
			Field: "uuid",
			Args:  []interface{}{types.NewNullUUID(req.UUID)},
		},
	}

	moderates, err := m.moderateRep.Listx(ctx, condition)
	if err != nil {
		return request.ModerateData{}, err
	}

	moderatesData := make([]request.ModerateData, 0, len(moderates))

	err = m.MapStructs(&moderatesData, &moderates)
	if err != nil {
		return request.ModerateData{}, err
	}
	if len(moderatesData) < 1 {
		return request.ModerateData{}, nil
	}

	uuids := make([]interface{}, 0, len(moderatesData))
	moderatesDataMap := make(map[string]*request.ModerateData, len(moderatesData))
	for i := range moderatesData {
		uuids = append(uuids, moderatesData[i].UUID)
		moderatesDataMap[moderatesData[i].UUID.String] = &moderatesData[i]
	}
	_, err = m.FeedGetFiles(ctx, moderatesDataMap, uuids...)
	if err != nil {
		return request.ModerateData{}, err
	}

	return moderatesData[0], nil
}

func (m *Moderates) FeedGetFiles(ctx context.Context, moderateDataMap map[string]*request.ModerateData, moderatesUUID ...interface{}) ([]request.FileData, error) {
	filesCondition := infoblog.Condition{
		Equal: &sq.Eq{"type": types.TypeUserModerate},
		In: &infoblog.In{
			Field: "foreign_uuid",
			Args:  moderatesUUID,
		},
	}

	files, err := m.Services.File.Listx(ctx, filesCondition)
	if err != nil {
		return nil, err
	}

	filesData := make([]request.FileData, 0, len(files))

	err = m.MapStructs(&filesData, &files)
	if err != nil {
		return nil, err
	}

	for i := range filesData {
		if filesData[i].Private.Valid && filesData[i].Private.Bool {
			filesData[i].Link = utils.PrivateLink(filesData[i])
		} else {
			filesData[i].Link = utils.PublicLink(filesData[i])
		}

		if moderateDataMap != nil {
			moderateDataMap[filesData[i].ForeignUUID.String].Files = append(moderateDataMap[filesData[i].ForeignUUID.String].Files, filesData[i])
		}
	}

	return filesData, err
}
