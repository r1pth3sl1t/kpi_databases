package model

import (
	"encoding/json"
	"errors"
	"fmt"
	_ "github.com/lib/pq"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/logger"
	"os"
	"reflect"
	"rgr/entity"
	"rgr/model/repository"
	"rgr/queries"
	"rgr/queries/search"
	"strconv"
	"strings"
)

type DatabaseCredentials struct {
	Host     string `json:"host"`
	Port     string `json:"port"`
	User     string `json:"user"`
	Password string `json:"password"`
	Dbname   string `json:"dbname"`
}

type DatabaseEntities struct {
	User       entity.User
	Skill      entity.Skill
	University entity.University
	Company    entity.Company
	Education  entity.Education
	Experience entity.Experience
}

type Model struct {
	db                   *gorm.DB
	userRepository       *repository.UserRepository
	skillRepository      *repository.Repository[entity.Skill, int]
	universityRepository *repository.Repository[entity.University, int]
	companyRepository    *repository.Repository[entity.Company, int]
	educationRepository  *repository.Repository[entity.Education, int]
	experienceRepository *repository.Repository[entity.Experience, int]
}

func New() (*Model, error) {
	config, err := os.Open("config.json")
	if err != nil {
		return nil, err
	}
	var dbc DatabaseCredentials
	decoder := json.NewDecoder(config)
	err = decoder.Decode(&dbc)

	if err != nil {
		return nil, err
	}

	m := new(Model)
	connStr := fmt.Sprintf("user=%s host=%s port=%s dbname=%s password=%s sslmode=disable",
		dbc.User, dbc.Host, dbc.Port, dbc.Dbname, dbc.Password)

	db, err := gorm.Open(postgres.Open(connStr), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})

	if err != nil {
		return nil, err
	}
	m.db = db

	m.userRepository = repository.NewUserRepository(m.db)
	m.skillRepository = repository.NewRepository[entity.Skill, int](m.db)
	m.universityRepository = repository.NewRepository[entity.University, int](m.db)
	m.companyRepository = repository.NewRepository[entity.Company, int](m.db)
	m.educationRepository = repository.NewRepository[entity.Education, int](m.db)
	m.experienceRepository = repository.NewRepository[entity.Experience, int](m.db)
	return m, err
}

func (m *Model) Close() {
	fmt.Println("Closing connection")
	db, err := m.db.DB()
	if err != nil {
		return
	}
	err = db.Close()
	if err != nil {
		return
	}
}

func (m *Model) FetchTableData() map[string][]string {

	tables := make(map[string][]string)
	for i := 0; i < reflect.ValueOf(DatabaseEntities{}).NumField(); i++ {
		entityName := reflect.TypeOf(DatabaseEntities{}).Field(i).Name
		for j := 0; j < reflect.ValueOf(DatabaseEntities{}).Field(i).NumField(); j++ {
			if (reflect.ValueOf(DatabaseEntities{}).Field(i).Type().Field(j).Tag.Get("requiredToInput") == "false") {
				continue
			}
			tables[entityName] =
				append(tables[entityName], reflect.ValueOf(DatabaseEntities{}).Field(i).Type().Field(j).Name)
		}
	}

	tables["Connection"] = []string{"UserId1", "UserId2"}
	tables["UserSkill"] = []string{"UserId", "SkillId"}
	return tables
}

func (m *Model) Insert(table string, data map[string]string) error {
	switch table {
	case "User":
		return m.userRepository.Create(entity.User{}.Create(data))
	case "Skill":
		return m.skillRepository.Create(entity.Skill{}.Create(data))
	case "University":
		return m.universityRepository.Create(entity.University{}.Create(data))
	case "Company":
		return m.companyRepository.Create(entity.Company{}.Create(data))
	case "Connection":
		user1Id, err := strconv.Atoi(data["UserId1"])
		user2Id, err := strconv.Atoi(data["UserId2"])

		if err != nil {
			return err
		}

		user1 := m.userRepository.FindById(user1Id)
		user2 := m.userRepository.FindById(user2Id)

		return m.userRepository.CreateConnection(user1, user2)
	case "UserSkill":
		userId, err := strconv.Atoi(data["UserId"])
		skillId, err := strconv.Atoi(data["SkillId"])

		if err != nil {
			return err
		}

		user := m.userRepository.FindById(userId)
		skill := m.skillRepository.FindById(skillId)

		return m.userRepository.AddSkillToUser(user, skill)
	case "Education":
		universityId, err := strconv.Atoi(data["UniversityId"])
		userId, err := strconv.Atoi(data["UserId"])

		user := m.userRepository.FindById(userId)
		university := m.universityRepository.FindById(universityId)

		education, err := entity.Education{}.Create(data)
		if err != nil {
			return err
		}

		return m.userRepository.AddEducationToUser(user, university, education)
	case "Experience":
		companyId, err := strconv.Atoi(data["CompanyId"])
		userId, err := strconv.Atoi(data["UserId"])

		if err != nil {
			return err
		}

		user := m.userRepository.FindById(userId)
		company := m.companyRepository.FindById(companyId)

		experience, err := entity.Experience{}.Create(data)
		if err != nil {
			return err
		}
		return m.userRepository.AddExperienceToUser(user, company, experience)
	}

	return nil
}

func (m *Model) FetchTablePrimaryKeys() map[string][]string {
	tables := make(map[string][]string)
	for i := 0; i < reflect.ValueOf(DatabaseEntities{}).NumField(); i++ {
		entityName := reflect.TypeOf(DatabaseEntities{}).Field(i).Name
		for j := 0; j < reflect.ValueOf(DatabaseEntities{}).Field(i).NumField(); j++ {

			if (!strings.Contains(reflect.ValueOf(DatabaseEntities{}).Field(i).Type().Field(j).Tag.Get("gorm"), "primaryKey")) {
				continue
			}
			tables[entityName] =
				append(tables[entityName], reflect.ValueOf(DatabaseEntities{}).Field(i).Type().Field(j).Name)
		}
	}

	tables["Connection"] = []string{"UserId1", "UserId2"}
	tables["UserSkill"] = []string{"UserId", "SkillId"}
	return tables
}

func (m *Model) Update(table string, data map[string]string, pkey map[string]string) error {
	switch table {
	case "User":
		userId, err := strconv.Atoi(pkey["UserId"])
		if err != nil {
			return err
		}

		user := m.userRepository.FindById(userId)
		return m.userRepository.Update(user, data)
	case "Skill":
		skillId, err := strconv.Atoi(pkey["SkillId"])
		if err != nil {
			return err
		}

		skill := m.skillRepository.FindById(skillId)
		return m.skillRepository.Update(skill, data)
	case "University":
		universityId, err := strconv.Atoi(pkey["UniversityId"])
		if err != nil {
			return err
		}

		university := m.universityRepository.FindById(universityId)
		return m.universityRepository.Update(university, data)
	case "Company":
		companyId, err := strconv.Atoi(pkey["CompanyId"])
		if err != nil {
			return err
		}

		company := m.companyRepository.FindById(companyId)
		return m.companyRepository.Update(company, data)
	case "Connection":
		return errors.New("update is unsupported")
	case "UserSkill":
		return errors.New("update is unsupported")
	case "Education":
		educationId, err := strconv.Atoi(pkey["EducationId"])
		if err != nil {
			return err
		}

		education := m.educationRepository.FindById(educationId)

		return m.educationRepository.Update(education, data)
	case "Experience":
		experienceId, err := strconv.Atoi(pkey["ExperienceId"])

		if err != nil {
			return err
		}

		experience := m.experienceRepository.FindById(experienceId)
		return m.experienceRepository.Update(experience, data)
	}
	return nil
}

func (m *Model) Delete(table string, pkey map[string]string) error {
	switch table {
	case "User":
		userId, err := strconv.Atoi(pkey["UserId"])
		if err != nil {
			return err
		}
		var user entity.User
		m.db.Preload(clause.Associations).Find(&user, userId)
		return m.userRepository.Delete(&user)
	case "Skill":
		skillId, err := strconv.Atoi(pkey["SkillId"])
		if err != nil {
			return err
		}

		skill := m.skillRepository.FindById(skillId)
		return m.skillRepository.Delete(skill)
	case "University":
		universityId, err := strconv.Atoi(pkey["UniversityId"])
		if err != nil {
			return err
		}

		university := m.universityRepository.FindById(universityId)
		return m.universityRepository.Delete(university)
	case "Company":
		companyId, err := strconv.Atoi(pkey["CompanyId"])
		if err != nil {
			return err
		}

		company := m.companyRepository.FindById(companyId)
		return m.companyRepository.Delete(company)
	case "Connection":
		user1Id, err := strconv.Atoi(pkey["UserId1"])
		if err != nil {
			return err
		}

		user2Id, err := strconv.Atoi(pkey["UserId2"])
		if err != nil {
			return err
		}

		user1 := m.userRepository.FindById(user1Id)
		user2 := m.userRepository.FindById(user2Id)

		return m.db.Model(&user1).Association("Connections").Delete(&user2)
	case "UserSkill":
		userId, err := strconv.Atoi(pkey["UserId"])
		if err != nil {
			return err
		}

		skillId, err := strconv.Atoi(pkey["SkillId"])
		if err != nil {
			return err
		}
		user := m.userRepository.FindById(userId)
		skill := m.skillRepository.FindById(skillId)

		return m.db.Model(&user).Association("Skills").Delete(&skill)
	case "Education":
		educationId, err := strconv.Atoi(pkey["EducationId"])
		if err != nil {
			return err
		}

		education := m.educationRepository.FindById(educationId)

		return m.educationRepository.Delete(education)
	case "Experience":
		experienceId, err := strconv.Atoi(pkey["ExperienceId"])

		if err != nil {
			return err
		}

		experience := m.experienceRepository.FindById(experienceId)
		return m.experienceRepository.Delete(experience)
	}
	return nil
}

func (m *Model) GenerateDataSet(size int) error {
	err := m.db.Exec(queries.GetUserGeneratingQuery(), size).Error
	if err != nil {
		return err
	}
	err = m.db.Exec(queries.GenerateSkillsQuery(), size).Error
	if err != nil {
		return err
	}
	err = m.db.Exec(queries.GenerateConnectionQuery(), size).Error
	if err != nil {
		return err
	}

	return nil
}

func (m *Model) Search(searcher search.SearchQueryHandler, data map[string]string) (int64, []string, error) {

	return 0, searcher.Search(m.db, repository.UpcastMap(data)), nil
}
