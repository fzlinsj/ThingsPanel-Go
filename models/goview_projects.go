package models

type GoviewProjects struct {
	ID           int32  `json:"ID,omitempty"`
	ProjectName  string `json:"ProjectName,omitempty"`
	State        int32  `json:"json_data,omitempty"`
	CreateTime   int64  `json:"CreateTime,omitempty"`
	UpdateTime   int64  `json:"UpdateTime,omitempty"`
	CreateUserId string `json:"CreateUserId,omitempty"`
	IsDelete     int8   `json:"IsDelete,omitempty"`
	IndexImage   string `json:"IndexImage,omitempty"`
	Remarks      string `json:"Remarks,omitempty"`
	TenantId     string `json:"tenant_id,omitempty"` // 租户id
}

func (GoviewProjects) TableName() string {
	return "GoviewProjects"
}
