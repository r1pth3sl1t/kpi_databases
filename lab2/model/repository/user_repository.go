package repository

import (
	"gorm.io/gorm"
	"rgr/entity"
)

type UserRepository struct {
	Repository[entity.User, int]
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{
		Repository[entity.User, int]{db},
	}
}

func (u *UserRepository) CreateConnection(user1 *entity.User, user2 *entity.User) error {
	//avoid manipulating join table directly, use ORM entities instead
	return u.db.Model(user1).Association("Connections").Append(user2)
}

func (u *UserRepository) AddSkillToUser(user *entity.User, skill *entity.Skill) error {
	return u.db.Model(user).Association("Skills").Append(skill)
}

func (u *UserRepository) AddEducationToUser(user *entity.User,
	university *entity.University,
	education *entity.Education) error {

	education.University = university

	return u.db.Model(user).Association("Education").Append(education)
}

func (u *UserRepository) AddExperienceToUser(user *entity.User,
	company *entity.Company,
	experience *entity.Experience) error {

	experience.Company = company

	return u.db.Model(user).Association("Experience").Append(experience)
}

func (u *UserRepository) Delete(user *entity.User) error {

	err := u.db.Model(user).Association("Skills").Clear()
	if err != nil {
		return err
	}
	err = u.db.Model(user).Association("Connections").Clear()
	if err != nil {
		return err
	}
	return u.db.Select("Experience", "Education").Delete(user).Error
}
