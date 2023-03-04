package logic

import (
	"bluebell/dao/mysql"
	"bluebell/models"
)

func GetCommunityList() ([]*models.Community, error) {
	// 查数据库， 查找到所有的community并且返回
	return mysql.GetCommunityList()
}

// GetCommunityDetail 分类详情
func GetCommunityDetail(id int64) (*models.CommunityDetail, error) {
	return mysql.GetCommunityDetailByID(id)
}
