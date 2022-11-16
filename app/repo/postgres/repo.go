package postgres

type Repository struct {
}

func NewRepo() *Repository {
	return &Repository{}
}

func (r *Repository) Close() error {
	return nil
}
