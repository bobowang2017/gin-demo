package dao

import (
	m "gin-demo/infra/model"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

type BaseDao struct {
}

func (b *BaseDao) Create(model interface{}) error {
	err := m.DB.Create(model).Error
	return err
}

// CreateInBatches 批量插入
// model:数组或者slice
// batchSize:一次要插入的数量,建议设置为100,如果要插入1000条,则会自动拆成10次插入
func (b *BaseDao) CreateInBatches(model interface{}, batchSize int) error {
	err := m.DB.CreateInBatches(model, batchSize).Error
	return err
}

func (b *BaseDao) GetObjById(model interface{}, id int) error {
	return m.DB.Model(model).Where(map[string]interface{}{"id": id}).First(model).Error
}

func (b *BaseDao) GetAllObj(model interface{}, result interface{}) error {
	return m.DB.Model(model).Find(result).Error
}

func (b *BaseDao) GetObjByCondition(model interface{}, params map[string]interface{}, result interface{}) error {
	if params == nil {
		return errors.New("params nil")
	}
	temp := m.DB.Model(model)
	page, okPage := params["page"].(int)
	size, okSize := params["size"].(int)
	if okPage && okSize {
		delete(params, "page")
		delete(params, "size")
		temp = temp.Offset((page - 1) * size).Limit(size).Where(params).Find(result)
	} else {
		temp = temp.Where(params).Find(result)
	}
	return temp.Error
}

func (b *BaseDao) GetByLikeCondition(model interface{}, params, likeParams map[string]interface{}, result interface{}) error {
	if params == nil {
		return errors.New("params nil")
	}
	temp := m.DB.Model(model)
	page, okPage := params["page"].(int)
	size, okSize := params["size"].(int)
	if okPage && okSize {
		delete(params, "page")
		delete(params, "size")
		temp = temp.Offset((page - 1) * size).Limit(size).Where(params)
	} else {
		temp = temp.Where(params)
	}
	for k, v := range likeParams {
		temp = temp.Where(k+" like ?", "%"+v.(string)+"%")
	}
	temp = temp.Find(result)
	return temp.Error
}

func (b *BaseDao) GetOneByCondition(model interface{}, params map[string]interface{}) (err error) {
	if params == nil {
		return errors.New("params nil")
	}
	return m.DB.Model(model).Offset(0).Limit(1).Where(params).First(model).Error
}

func (b *BaseDao) DeleteObjById(model interface{}, id int) error {
	if err := m.DB.Where(map[string]interface{}{"id": id}).Delete(model).Error; err != nil {
		return err
	}
	return nil
}

func (b *BaseDao) DeleteObjByIds(model interface{}, ids *[]int) error {
	if err := m.DB.Where("id in ?", *ids).Delete(model).Error; err != nil {
		return err
	}
	return nil
}

func (b *BaseDao) DeleteObjByCondition(model interface{}, params map[string]interface{}) error {
	err := m.DB.Where(params).Delete(model).Error
	if err != nil {
		return err
	}
	return nil
}

func (b *BaseDao) UpdateObjById(model interface{}, id int, params map[string]interface{}) error {
	err := m.DB.Model(model).Where(map[string]interface{}{"id": id}).Updates(params).Error
	if err != nil {
		return err
	}
	return nil
}

func (b *BaseDao) UpdateObjByCondition(model interface{}, params map[string]interface{}, field map[string]interface{}) error {
	err := m.DB.Model(model).Where(params).Updates(field).Error
	if err != nil {
		return err
	}
	return nil
}

func (b *BaseDao) TotalObj(model interface{}, params map[string]interface{}) (int64, error) {
	var total int64
	err := m.DB.Model(model).Where(params).Count(&total).Error
	return total, err
}

func (b *BaseDao) TotalPlus(model interface{}, params map[string]interface{}) (int64, error) {
	var total int64
	err := m.DB.Model(model).Select("id").Where(params).Count(&total).Error
	return total, err
}

func (b *BaseDao) TotalByCondition(model interface{}, params, likeParams map[string]interface{}) (int64, error) {
	var (
		temp  = m.DB.Model(model).Where(params)
		total int64
	)
	for k, v := range likeParams {
		temp = temp.Where(k+" like ?", "%"+v.(string)+"%")
	}
	temp = temp.Count(&total)
	return total, temp.Error
}

func (b *BaseDao) ExecSql(sql string, values ...interface{}) (*gorm.DB, error) {
	db := m.DB.Raw(sql, values...)
	if err := db.Error; err != nil {
		return db, err
	}
	return db, nil
}
