package entity

import (
	"strconv"
	"time"
)

type Experience struct {
	ExperienceId int       `gorm:"column:exp_id;primaryKey" requiredToInput:"false"`
	UserId       int       `gorm:"column:user_id"`
	Company      *Company  `gorm:"foreignKey:CompanyId" requiredToInput:"false"`
	CompanyId    int       `gorm:"column:company_id"`
	DateStart    time.Time `gorm:"column:date_start"`
	DateEnd      time.Time `gorm:"column:date_end"`
	Position     string    `gorm:"column:position"`
	JobOverview  string    `gorm:"column:job_overview"`
}

func (Experience) TableName() string {
	return "experience"
}

func (Experience) Create(fields map[string]string) (*Experience, error) {
	experience := Experience{}
	companyId, err := strconv.Atoi(fields["CompanyId"])

	experience.CompanyId = companyId
	experience.Position = fields["Position"]
	experience.JobOverview = fields["JobOverview"]
	experience.DateStart, err = time.Parse("2006-01-02", fields["DateStart"])
	experience.DateEnd, err = time.Parse("2006-01-02", fields["DateEnd"])
	if err != nil {
		return nil, err
	}
	return &experience, err

}
