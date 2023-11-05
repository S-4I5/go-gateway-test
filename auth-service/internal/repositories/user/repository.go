package user

type Holder interface {
	SaveUser(username, password string) (int, error)
	GetUser(username string) (*User, error)
}

type Repository struct {
	Holder
}

func (rep *Repository) SaveUser(username, password string) (int, error) {
	return rep.Holder.SaveUser(username, password)
}

func (rep *Repository) GetPassword(username string) (*User, error) {
	return rep.Holder.GetUser(username)
}
