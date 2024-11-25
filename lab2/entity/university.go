package entity

type University struct {
	UniversityId int    `gorm:"column:university_id;primaryKey" requiredToInput:"false"`
	Name         string `gorm:"column:name"`
	Country      string `gorm:"country"`
}

func (University) TableName() string {
	return "university"
}

func (University) Create(fields map[string]string) *University {
	return &University{
		Name:    fields["Name"],
		Country: fields["Country"],
	}
}
