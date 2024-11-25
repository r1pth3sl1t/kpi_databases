package entity

import (
	"strconv"
	"time"
)

type Education struct {
	EducationId  int `gorm:"column:education_id;primaryKey " requiredToInput:"false"`
	UniversityId int
	UserId       int         `gorm:"column:user_id"`
	University   *University `gorm:"foreignKey:UniversityId" requiredToInput:"false"`
	DateStart    time.Time   `gorm:"column:date_start"`
	DateEnd      time.Time   `gorm:"column:date_end"`
	Degree       string      `gorm:"column:degree"`
	Specialty    string      `gorm:"column:specialty"`
}

func (Education) TableName() string {
	return "education"
}

func (Education) Create(fields map[string]string) (*Education, error) {
	education := Education{}
	universityId, err := strconv.Atoi(fields["UniversityId"])
	education.UniversityId = universityId
	education.DateStart, err = time.Parse("2006-01-02", fields["DateStart"])
	education.DateEnd, err = time.Parse("2006-01-02", fields["DateEnd"])
	if err != nil {
		return nil, err
	}
	education.Degree = fields["Degree"]
	education.Specialty = fields["Specialty"]
	return &education, err
}
