package address

import (
	"strings"

	"github.com/nguyentrunghieu15/be-beehome-prj/address-service/internal/datasource"
	"github.com/nguyentrunghieu15/be-beehome-prj/internal/database"
)

type IAddressRepo interface {
	FindAllAddressByQuery(string) ([]datasource.Address, error)
	GetAllWardsByDistrictCode(string) ([]datasource.Ward, error)
	GetAllDistrictsByProvinceCode(string) ([]datasource.District, error)
	GetAllProvince() ([]datasource.Province, error)
	CheckAddressExist(datasource.Address) (bool, error)
}

type AddressRepo struct {
	db *database.PostgreDb
}

func NewAddressRepo(db *database.PostgreDb) *AddressRepo {
	return &AddressRepo{
		db: db,
	}
}

var stopWord = []string{"Tỉnh", "Huyện", "Xã", "Phường", "Thị trấn", "Thành Phố"}

func fixToTextSearchQuery(query string) []string {

	// Split theo dau phẩy
	splited := strings.Split(query, ",")

	// Loại bỏ dấu cách dư thừa
	var normalizes []string
	for _, v := range splited {
		normalize := strings.Join(strings.Fields(v), " ")
		for _, stopW := range stopWord {
			normalize = strings.Replace(normalize, stopW, " ", -1)
		}
		normalize = strings.Replace(normalize, " ", " <-> ", -1)
		normalizes = append(normalizes, normalize)
	}

	switch len(normalizes) {
	case 1:
		return []string{normalizes[0], normalizes[0], normalizes[0]}
	case 2:
		return []string{normalizes[0], normalizes[0] + " | " + normalizes[1], normalizes[1]}
	case 3:
		return normalizes
	}
	return []string{strings.Join(normalizes, " | "), strings.Join(normalizes, " | "), strings.Join(normalizes, " | ")}
}

func (a *AddressRepo) FindAllAddressByQuery(query string) ([]datasource.Address, error) {
	sql := `
SELECT
  provinces.full_name AS province_full_name,
  districts.full_name AS district_full_name,
  wards.full_name AS ward_full_name,
  RANK() OVER (ORDER BY 
      COALESCE(ts_rank(to_tsvector('simple', unaccent(provinces.full_name)), to_tsquery('simple', unaccent(?))), 0) +
      COALESCE(ts_rank(to_tsvector('simple', unaccent(districts.full_name)), to_tsquery('simple', unaccent(?))), 0) +
      COALESCE(ts_rank(to_tsvector('simple', unaccent(wards.full_name)), to_tsquery('simple', unaccent(?))), 0) DESC
  ) AS rank
FROM provinces
LEFT JOIN districts ON provinces.code = districts.province_code
LEFT JOIN wards ON districts.code = wards.district_code
WHERE (
  to_tsvector('simple', unaccent(provinces.full_name)) @@ to_tsquery('simple', unaccent(?)) AND
  provinces.full_name IS NOT NULL
) OR (
  to_tsvector('simple', unaccent(districts.full_name)) @@ to_tsquery('simple', unaccent(?)) AND
  districts.full_name IS NOT NULL
) OR (
  to_tsvector('simple', unaccent(wards.full_name)) @@ to_tsquery('simple', unaccent(?)) AND
  wards.full_name IS NOT NULL
)
ORDER BY rank DESC
LIMIT 10;`

	var result []datasource.Address
	tsquery := fixToTextSearchQuery(query)

	err := a.db.Raw(sql, tsquery[2], tsquery[1], tsquery[0], tsquery[2], tsquery[1], tsquery[0]).Scan(&result)
	if err.Error != nil {
		return nil, err.Error
	}
	return result, nil
}

func (a *AddressRepo) GetAllWardsByDistrictCode(districtCode string) ([]datasource.Ward, error) {
	var wards []datasource.Ward
	err := a.db.Where("district_code = ?", districtCode).Find(&wards).Error
	if err != nil {
		return nil, err
	}
	return wards, nil
}

// GetAllDistrictsByProvinceCode gets all districts by province code
func (a *AddressRepo) GetAllDistrictsByProvinceCode(provinceCode string) ([]datasource.District, error) {
	var districts []datasource.District
	err := a.db.Where("province_code = ?", provinceCode).Find(&districts).Error
	if err != nil {
		return nil, err
	}
	return districts, nil
}

// GetAllDistrictsByProvinceCode gets all districts by province code
func (a *AddressRepo) GetAllProvince() ([]datasource.Province, error) {
	var province []datasource.Province
	err := a.db.Find(&province).Error
	if err != nil {
		return nil, err
	}
	return province, nil
}

func (a *AddressRepo) CheckAddressExist(address datasource.Address) (bool, error) {
	sql := `
SELECT
  provinces.full_name AS province_full_name,
  districts.full_name AS district_full_name,
  wards.full_name AS ward_full_name
FROM provinces
LEFT JOIN districts ON provinces.code = districts.province_code
LEFT JOIN wards ON districts.code = wards.district_code
WHERE provinces.full_name = ? AND districts.full_name = ? AND wards.full_name = ?;`

	var temp []datasource.Address
	err := a.db.Raw(sql, address.ProvinceFullName, address.DistrictFullName, address.WardFullName).Scan(&temp)
	if err.Error != nil {
		return false, err.Error
	}

	return err.RowsAffected > 0, nil
}
