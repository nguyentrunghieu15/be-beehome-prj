package datasource

import "github.com/nguyentrunghieu15/be-beehome-prj/internal/database"

type IPostalCodeRepo interface {
	FindPostalCodes(interface{}) ([]*PostalCode, error)
	FindOneById(int32) (*PostalCode, error)
	FindPostalCodesByCountryCode(string) ([]*PostalCode, error)
	FindPostalCodesByZipcode(string) ([]*PostalCode, error)
	UpdateOneById(int32, map[string]interface{}) (*PostalCode, error)
	CreatePostalCode(map[string]interface{}) (*PostalCode, error)
	DeleteOneById(int32) error
}

type PostalCodeRepo struct {
	db *database.PostgreDb
}

func (pr *PostalCodeRepo) FindOneById(id int32) (*PostalCode, error) {
	postalCode := &PostalCode{}
	result := pr.db.First(postalCode, "id = ?", id)
	if result.Error != nil {
		return nil, result.Error
	}
	return postalCode, nil
}

func (pr *PostalCodeRepo) FindPostalCodesByCountryCode(countryCode string) ([]*PostalCode, error) {
	var postalCodes []*PostalCode
	result := pr.db.Find(&postalCodes, "country_code = ?", countryCode)
	if result.Error != nil {
		return nil, result.Error
	}
	return postalCodes, nil
}

func (pr *PostalCodeRepo) FindPostalCodesByZipcode(zipcode string) ([]*PostalCode, error) {
	var postalCodes []*PostalCode
	result := pr.db.Find(&postalCodes, "zipcode = ?", zipcode)
	if result.Error != nil {
		return nil, result.Error
	}
	return postalCodes, nil
}

func (pr *PostalCodeRepo) UpdateOneById(id int32, updateParams map[string]interface{}) (*PostalCode, error) {
	_, err := pr.FindOneById(id)
	if err != nil {
		return nil, err
	}

	result := pr.db.Model(&PostalCode{}).Where("id = ?", id).Updates(updateParams)
	if result.Error != nil {
		return nil, result.Error
	}
	return pr.FindOneById(id)
}

func (pr *PostalCodeRepo) CreatePostalCode(data PostalCode) (*PostalCode, error) {
	// Create the record in the database
	result := pr.db.Create(&data)
	if result.Error != nil {
		return nil, result.Error
	}

	return &data, nil
}

func (pr *PostalCodeRepo) DeleteOneById(id int32) error {
	result := pr.db.Delete(&PostalCode{}, "id = ?", id)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (pr *PostalCodeRepo) FindPostalCodes(req interface{}) ([]*PostalCode, error) {
	return nil, nil
}

func NewPostalCodeRepo(db *database.PostgreDb) *PostalCodeRepo {
	return &PostalCodeRepo{
		db: db,
	}
}
