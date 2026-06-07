package user

type MemoryRepository struct{}

func (MemoryRepository) Save(User) error { return nil }
