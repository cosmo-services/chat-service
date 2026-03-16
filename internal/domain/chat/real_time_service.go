package chat

type RealTimeService struct {
	publisher Publisher
}

func (s *RealTimeService) NewMessage(payload NewMessagePayload) error {

	return nil
}
