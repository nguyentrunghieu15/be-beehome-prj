package datasource

type AdministrativeRegion struct {
	ID         int    `json:"id"`
	Name       string `json:"name"`
	NameEn     string `json:"name_en"`
	CodeName   string `json:"code_name"`
	CodeNameEn string `json:"code_name_en"`
}

type AdministrativeUnit struct {
	ID          int    `json:"id"`
	FullName    string `json:"full_name"`
	FullNameEn  string `json:"full_name_en"`
	ShortName   string `json:"short_name"`
	ShortNameEn string `json:"short_name_en"`
	CodeName    string `json:"code_name"`
	CodeNameEn  string `json:"code_name_en"`
}

type Province struct {
	Code                   string `json:"code"`
	Name                   string `json:"name"`
	NameEn                 string `json:"name_en"`
	FullName               string `json:"full_name"`
	FullNameEn             string `json:"full_name_en"`
	CodeName               string `json:"code_name"`
	AdministrativeUnitID   int    `json:"administrative_unit_id"`
	AdministrativeRegionID int    `json:"administrative_region_id"`
}

type District struct {
	Code                 string `json:"code"`
	Name                 string `json:"name"`
	NameEn               string `json:"name_en"`
	FullName             string `json:"full_name"`
	FullNameEn           string `json:"full_name_en"`
	CodeName             string `json:"code_name"`
	ProvinceCode         string `json:"province_code"`
	AdministrativeUnitID int    `json:"administrative_unit_id"`
}

type Ward struct {
	Code                 string `json:"code"`
	Name                 string `json:"name"`
	NameEn               string `json:"name_en"`
	FullName             string `json:"full_name"`
	FullNameEn           string `json:"full_name_en"`
	CodeName             string `json:"code_name"`
	DistrictCode         string `json:"district_code"`
	AdministrativeUnitID int    `json:"administrative_unit_id"`
}

type Address struct {
	WardFullName     string `json:"wards_full_name"`
	DistrictFullName string `json:"districts_full_name"`
	ProvinceFullName string `json:"provinces_full_name"`
}
