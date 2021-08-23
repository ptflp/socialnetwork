package infoblog

type Repositories struct {
	Files        FileRepository
	Posts        PostRepository
	Users        UserRepository
	Subscribers  SubscriberRepository
	Likes        LikeRepository
	Chats        ChatRepository
	ChatMessages ChatMessagesRepository
	Events       EventRepository
}

type Entity interface {
	TableName() string
}
