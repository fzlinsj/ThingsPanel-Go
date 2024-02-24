package services

import (
	"ThingsPanel-Go/initialize/psql"
	"ThingsPanel-Go/models"
	valid "ThingsPanel-Go/validate"
	"errors"

	"gorm.io/gorm"
)

type GoViewDashboardService struct {
	//可搜索字段
	SearchField []string
	//可作为条件的字段
	WhereField []string
	//可做为时间范围查询的字段
	TimeField []string
}

func (*GoViewDashboardService) GetTpDashboardDetail(tp_dashboard_id string) []models.TpDashboard {
	var tp_dashboard []models.TpDashboard
	psql.Mydb.First(&tp_dashboard, "id = ?", tp_dashboard_id)
	return tp_dashboard
}

// 获取列表
func (*GoViewDashboardService) GetGoViewDashboardList(PaginationValidate valid.GoViewDashboardPaginationValidate, tenantId string) (bool, []models.GoviewProjects, int64) {
	var GoViewDashboards []models.GoviewProjects
	offset := (PaginationValidate.Page - 1) * PaginationValidate.Limit
	db := psql.Mydb.Model(&models.GoviewProjects{}).Where("tenant_id = ? ", tenantId)
	// if PaginationValidate.RelationId != "" {
	// 	db.Where("relation_id = ?", PaginationValidate.RelationId)
	// }
	// if PaginationValidate.Id != "" {
	// 	db.Where("id = ?", PaginationValidate.Id)
	// }
	var count int64
	db.Count(&count)
	result := db.Limit(PaginationValidate.Limit).Offset(offset).Order("CreateTime desc").Find(&GoViewDashboards)
	if result.Error != nil {
		errors.Is(result.Error, gorm.ErrRecordNotFound)
		return false, GoViewDashboards, 0
	}
	return true, GoViewDashboards, count
}

// 新增数据
func (*GoViewDashboardService) AddTpDashboard(tp_dashboard models.TpDashboard) (bool, models.TpDashboard) {
	result := psql.Mydb.Create(&tp_dashboard)
	if result.Error != nil {
		errors.Is(result.Error, gorm.ErrRecordNotFound)
		return false, tp_dashboard
	}
	return true, tp_dashboard
}

// 修改数据
func (*GoViewDashboardService) EditTpDashboard(tp_dashboard valid.TpDashboardValidate, tenantId string) bool {
	result := psql.Mydb.Model(&models.TpDashboard{}).Where("id = ? and tenant_id = ?", tp_dashboard.Id, tenantId).Updates(&tp_dashboard)
	if result.Error != nil {
		errors.Is(result.Error, gorm.ErrRecordNotFound)
		return false
	}
	return true
}

// 删除数据
func (*GoViewDashboardService) DeleteTpDashboard(tp_dashboard models.TpDashboard) bool {
	result := psql.Mydb.Delete(&tp_dashboard)
	if result.Error != nil {
		errors.Is(result.Error, gorm.ErrRecordNotFound)
		return false
	}
	return true
}
