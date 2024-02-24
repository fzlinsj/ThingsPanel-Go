package valid

import "ThingsPanel-Go/models"

type GoViewDashboardPaginationValidate struct {
	Page     int    `json:"page"  alias:"当前页" valid:"Required;Min(1)"`
	Limit    int    `json:"limit"  alias:"每页页数" valid:"Required;Max(10000)"`
	TenantId string `json:"tenant_id"  alias:"租户ID" valid:"Required"`
}

type RspGoViewDashboardPaginationValidate struct {
	CurrentPage int                     `json:"current_page"  alias:"当前页" valid:"Required;Min(1)"`
	PerPage     int                     `json:"per_page"  alias:"每页页数" valid:"Required;Max(10000)"`
	Data        []models.GoviewProjects `json:"data" alias:"返回数据"`
	Total       int64                   `json:"total" alias:"总数" valid:"Max(10000)"`
}
