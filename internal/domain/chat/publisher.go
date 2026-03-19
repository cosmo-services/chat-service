package chat

type Publisher interface {
	PublishToUser(userId string, eventType string, collectionView *CollectionView) error
	PublishToUsers(eventType string, collectionViewMap map[string]*CollectionView) error
}
