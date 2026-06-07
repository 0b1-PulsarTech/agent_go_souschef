package user

type EventPublisher struct{}

func (EventPublisher) Publish(User) error { return nil }
