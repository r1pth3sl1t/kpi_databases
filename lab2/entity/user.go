package entity

type User struct {
	UserId      int          `gorm:"column:user_id;primaryKey" requiredToInput:"false"`
	FirstName   string       `gorm:"column:first_name"`
	LastName    string       `gorm:"column:last_name"`
	Email       string       `gorm:"column:email"`
	Skills      []Skill      `gorm:"many2many:users_to_skills;joinForeignKey:user_id;joinReferences:skill_id;onDelete:SET NULL" requiredToInput:"false"`
	Connections []*User      `gorm:"many2many:connection;joinForeignKey:u1;joinReferences:u2;onDelete:SET NULL" requiredToInput:"false"`
	Education   []Education  `gorm:"foreignKey:UserId;onDelete:CASCADE" requiredToInput:"false"`
	Experience  []Experience `gorm:"foreignKey:UserId;onDelete:CASCADE" requiredToInput:"false"`
}

func (User) TableName() string {

	return "user"
}

func (User) Create(fields map[string]string) *User {

	return &User{
		FirstName: fields["FirstName"],
		LastName:  fields["LastName"],
		Email:     fields["Email"],
	}
}
