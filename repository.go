package infoblog

type Repositories struct {
	Files           FileRepository
	Posts           PostRepository
	Users           UserRepository
	Subscribers     SubscriberRepository
	Likes           LikeRepository
	Chats           ChatRepository
	ChatMessages    ChatMessagesRepository
	ChatParticipant ChatParticipantRepository
	Comments        CommentsRepository
	Events          EventRepository
	Friends         FriendRepository
}

type Tabler interface {
	TableName() string
	OnCreate() string
}
