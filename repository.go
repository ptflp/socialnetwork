package infoblog

type Repositories struct {
	Files       FileRepository
	Posts       PostRepository
	Users       UserRepository
	Subscribers SubscriberRepository
	Likes       LikeRepository
}
