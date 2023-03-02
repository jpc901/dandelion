package logic

import (
	"bluebell/dao/mysql"
	"bluebell/models"
)

func GetCommunityList() ([]*models.Community, error) {

	// 查数据库， 查找到所有的community并且返回
	return mysql.GetCommunityList()
}
