package entity

type Skill struct {
	SkillId   int    `gorm:"column:skill_id; primaryKey" requiredToInput:"false"`
	SkillName string `gorm:"column:name"`
	SkillType string `gorm:"column:skill_type"`
}

func (Skill) Create(fields map[string]string) *Skill {
	return &Skill{
		SkillName: fields["SkillName"],
		SkillType: fields["SkillType"],
	}
}

func (Skill) TableName() string {
	return "skill"
}
