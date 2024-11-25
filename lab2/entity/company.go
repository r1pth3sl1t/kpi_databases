package entity

type Company struct {
	CompanyId   int    `gorm:"column:company_id; primaryKey" requiredToInput:"false"`
	Name        string `gorm:"column:name"`
	Website     string `gorm:"column:website_link"`
	Description string `gorm:"column:description"`
}

func (Company) TableName() string {
	return "company"
}

func (Company) Create(fields map[string]string) *Company {
	return &Company{
		Name:        fields["Name"],
		Website:     fields["Website"],
		Description: fields["Description"],
	}
}
