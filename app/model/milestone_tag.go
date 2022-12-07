package model

import "encoding/json"

type MilestoneTag struct {
	ID          uint   `json:"id"`
	MilestoneID uint   `json:"milestoneId" form:"milestoneId" gorm:"index:idx_milestone_id_name,unique"`
	Name        string `json:"name" form:"name" gorm:"index:idx_milestone_id_name,unique"`
}

func (m MilestoneTag) MarshalJSON() ([]byte, error) {
	return json.Marshal(m.Name)
}
