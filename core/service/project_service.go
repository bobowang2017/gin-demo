package service

import (
	"gin-demo/core/dao"
	"gin-demo/core/dto"
	m "gin-demo/core/model"
	"gin-demo/infra/utils/log"
	"strconv"
	"strings"
)

// ProjectTreeNode 定义Project树状结构体
type ProjectTreeNode struct {
	*m.Project
	Children []*ProjectTreeNode `json:"children,omitempty"`
}

type IProjectService interface {
	List(queryDto *dto.ProjectListQuery) ([]*ProjectTreeNode, int64, error)
}

type projectService struct {
	projectDao dao.IProjectDao
}

func (p *projectService) List(queryDto *dto.ProjectListQuery) ([]*ProjectTreeNode, int64, error) {
	var (
		total  int64
		rows   []*m.Project
		params = make(map[string]interface{})
		err    error
	)
	//没有搜索条件，直接按照root_parent_id进行归类查询即可
	if queryDto.Age == nil && queryDto.Name == "" && queryDto.Code == "" {
		if rows, err = p.listWithOutCondition(queryDto.Page, queryDto.Size); err != nil {
			log.Logger.Error(err)
			return nil, 0, err
		}
		total = p.totalWithOutCondition()
		return p.generateTree(rows), total, nil
	}

	if queryDto.Name != "" {
		params["name"] = queryDto.Name
	}
	if queryDto.Age != nil {
		params["age"] = *queryDto.Age
	}
	if queryDto.Code != "" {
		params["code"] = queryDto.Code
	}

	if rows, err = p.listWithCondition(queryDto.Page, queryDto.Size, params); err != nil {
		log.Logger.Error(err)
		return nil, 0, err
	}

	if total, err = p.projectDao.GetTotalDisRootParentId(params); err != nil {
		log.Logger.Error(err)
		return nil, 0, err
	}

	return p.generateTree(rows), total, nil
}

func (p *projectService) listWithOutCondition(page, size int) ([]*m.Project, error) {
	var (
		result  []*m.Project
		rootIds []int
		err     error
	)
	if rootIds, err = p.projectDao.GetRootParentIdByCondition(page, size, nil); err != nil {
		return result, err
	}
	if result, err = p.projectDao.GetByRootParentIds(rootIds); err != nil {
		return result, err
	}
	return result, err
}

func (p *projectService) totalWithOutCondition() int64 {
	var (
		total int64
		err   error
	)
	if total, err = p.projectDao.GetTotalDistinctRootParentId(); err != nil {
		log.Logger.Error(err)
	}
	return total
}

func (p *projectService) generateTree(rows []*m.Project) []*ProjectTreeNode {
	var (
		sortedKeys []int
		result     []*ProjectTreeNode
		rowMap     = make(map[int]*ProjectTreeNode)
	)

	for _, row := range rows {
		treeNode := &ProjectTreeNode{
			Project:  row,
			Children: make([]*ProjectTreeNode, 0, 0),
		}
		rowMap[row.ID] = treeNode
		sortedKeys = append(sortedKeys, row.ID)
	}

	for _, v := range sortedKeys {
		if rowMap[v].ParentId == -1 {
			result = append(result, rowMap[v])
		} else {
			rowMap[rowMap[v].ParentId].Children = append(rowMap[rowMap[v].ParentId].Children, rowMap[v])
		}
	}

	return result
}

func (p *projectService) listWithCondition(page, size int, params map[string]interface{}) ([]*m.Project, error) {
	var (
		ids           = make([]int64, 0, 0)
		idMap         = make(map[int64]bool)
		result        = make([]*m.Project, 0, 0)
		rootParentIds []int
		pathIds       []string
		err           error
	)
	// 首先根据搜索条件找到对应的节点，找到的结果集中可能是父节点，也可能是叶子节点，但是分页是按照根节点进行的，故先查询出分页后的根节点
	if rootParentIds, err = p.projectDao.GetDisRootParentIds(page, size, params); err != nil {
		return result, err
	}
	// 结合查询后的根节点和搜索条件可以查询到分页后的根节点及叶子节点信息，但是可能会漏掉相关的中间节点
	if pathIds, err = p.projectDao.GetPathIdsByRootParentIds(rootParentIds, params); err != nil {
		return result, err
	}
	//根据查询命中的节点的pathIds信息即可查找到相关结果集的所有节点列表
	for _, pathId := range pathIds {
		tmpIds := strings.Split(pathId, ",")
		for _, tmpId := range tmpIds {
			tmpIdInt, _ := strconv.ParseInt(tmpId, 10, 32)
			if _, ok := idMap[tmpIdInt]; ok {
				continue
			}
			idMap[tmpIdInt] = true
			ids = append(ids, tmpIdInt)
		}
	}
	//根据节点列表反查即可
	if result, err = p.projectDao.GetByIds(ids); err != nil {
		return result, err
	}
	return result, err
}

func NewProjectService() IProjectService {
	return &projectService{dao.GetProjectDao()}
}
