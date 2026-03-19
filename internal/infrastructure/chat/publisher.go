package chat_infrastructure

import (
	chat_domain "main/internal/domain/chat"
	"main/pkg"
)

type Publisher struct {
	logger pkg.Logger
}

func NewPublisher(logger pkg.Logger) chat_domain.Publisher {
	return &Publisher{logger: logger}
}
func (p *Publisher) PublishToUser(userId string, eventType string, collectionView *chat_domain.CollectionView) error {
	p.logger.Infof("collection type %s: %s", eventType, collectionView)

	return nil
}
func (p *Publisher) PublishToUsers(eventType string, collectionViewMap map[string]*chat_domain.CollectionView) error {
	p.logger.Infof("collection type %s: %s", eventType, collectionViewMap)

	return nil
}
