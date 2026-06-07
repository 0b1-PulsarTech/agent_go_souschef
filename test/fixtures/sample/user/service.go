package user

type Repository interface {
	Save(User) error
}

type Publisher interface {
	Publish(User) error
}

type User struct{ Name string }

func CreateUser(repo Repository, pub Publisher, user User) error {
	if err := repo.Save(user); err != nil {
		return err
	}
	return pub.Publish(user)
}
