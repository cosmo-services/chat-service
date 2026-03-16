package chat

type Publisher interface {
	Publish(userId string, collection *Collection) error
}
