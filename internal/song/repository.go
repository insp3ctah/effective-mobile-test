package song

import (
	"effective-mobile-test/pkg"
)

type Repository struct{}

func NewRepository() *Repository {
	return &Repository{}
}

func (r *Repository) Create(song *Song) error {
	return pkg.DB.Create(song).Error
}

func (r *Repository) GetAll() ([]Song, error) {
	var songs []Song
	err := pkg.DB.Find(&songs).Error
	return songs, err
}

func (r *Repository) GetByID(id int) (*Song, error) {
	var song Song
	err := pkg.DB.First(&song, id).Error
	return &song, err
}

func (r *Repository) Delete(id uint) error {
	return pkg.DB.Delete(&Song{}, id).Error
}

func (r *Repository) Update(song *Song) error {
	return pkg.DB.Save(song).Error
}

func (r *Repository) GetAllWithFilter(filter Song) ([]Song, error) {
	var songs []Song
	query := pkg.DB.Model(&Song{})

	if filter.Group != "" {
		query = query.Where("\"group\" LIKE ?", "%"+filter.Group+"%")
	}
	if filter.Title != "" {
		query = query.Where("title LIKE ?", "%"+filter.Title+"%")
	}
	if filter.ReleaseDate != "" {
		query = query.Where("release_date = ?", filter.ReleaseDate)
	}

	err := query.Find(&songs).Error
	return songs, err
}
