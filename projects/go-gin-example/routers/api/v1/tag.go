package v1

import (
	"net/http"

	"github.com/astaxie/beego/validation"
	"github.com/fcf/go-gin-example/models"
	"github.com/fcf/go-gin-example/pkg/e"
	"github.com/fcf/go-gin-example/pkg/logging"
	"github.com/fcf/go-gin-example/pkg/setting"
	"github.com/fcf/go-gin-example/pkg/util"
	"github.com/gin-gonic/gin"
	"github.com/unknwon/com"
)

// @summary 获取多个文章标签
// @Produce json
// @Param token query string true "Token"
// @Success 200 {string} json "{"code":200,"data":{"lists": [{"id": 4,"created_on": 1685854990,"modified_on": 0,"deleted_on": 0,"name": "4","created_by": "fcf","modified_by": "","state": 1}]},"msg":"ok"}"
// @Router /api/v1/tags [GET]
func GetTags(c *gin.Context) {
	// c.Quuery 用户获取?name=test&state=1这类URL参数
	// c.DefaultQuery则支持设置一个默认值
	name := c.Query("name")
	logging.Debug("xxxxxxxx fcf 开始获取tags")

	maps := make(map[string]interface{})
	data := make(map[string]interface{})

	if name != "" {
		maps["name"] = name
	}
	var state int = -1
	if arg := c.Query("state"); arg != "" {
		state = com.StrTo(arg).MustInt()
		maps["state"] = state
	}
	code := e.SUCCESS
	logging.Debug("xxxxxxxx fcf : 开始获取tags")
	// util.GetPage保证了各接口的page处理是一致的
	lists, err := models.GetTags(util.GetPage(c), setting.PageSize, maps)
	if err != nil {
		code = e.ERROR
		logging.Debug("xxxxxxxx fcf : 开始获取tags list 失败:", err)
	}
	data["lists"] = lists
	data["total"] = models.GetTagTotal(maps)

	c.JSON(http.StatusOK, gin.H{
		"code": code,
		"msg":  e.GetMsg(code),
		"data": data,
	})
}

// @summary 新增文章标签
// @Produce json
// @Param token query string true "Token"
// @Param name query string true "Name"
// @Param state query int false "State"
// @Success 200 {string} json "{"code":200,"data":{},"msg":"ok"}"
// @Router /api/v1/tags [POST]
func AddTag(c *gin.Context) {
	name := c.Query("name")
	logging.Debug("xxxxxxxx fcf 开始添加tag name: %s", name)
	state := com.StrTo(c.DefaultQuery("state", "0")).MustInt()
	createdBy := c.Query("created_by")

	valid := validation.Validation{}
	valid.Required(name, "name").Message("名称不能为空")
	valid.MaxSize(name, 100, "name").Message("名称最长100字节")
	valid.Required(createdBy, "created_by").Message("创始人不能为空")
	valid.MaxSize(createdBy, 100, "created_by").Message("创始人最长100字节")
	valid.Range(state, 0, 1, "state").Message("状态只允许0或1")

	code := e.INVALID_PARAMS
	if !valid.HasErrors() {
		exist, err := models.ExistTagByName(name)
		if !exist {
			logging.Debug("xxxxxxxx fcf model 数据不存在")
			err := models.AddTag(name, state, createdBy)
			if err != nil {
				logging.Debug("add tag failed: %x", err)
				code = e.ERROR
			} else {
				logging.Debug("add tag success")
				code = e.SUCCESS
			}
		} else if err != nil {
			code = e.ERROR
		} else {
			code = e.ERROR_EXIST_TAG
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"code": code,
		"msg":  e.GetMsg(code),
		"data": make(map[string]string),
	})
}

// @summary 修改文章标签
// @Produce json
// @Param token query string true "Token"
// @Param name query string true "Name"
// @Param state query int false "State"
// @Param modified_by query string false "State"
// @Success 200 {string} json "{"code":200,"data":{},"msg":"ok"}"
// @Router /api/v1/tags/:id [PUT]
func EditTag(c *gin.Context) {
	logging.Debug("xxxxxxxx fcf 修改tag")
	id := com.StrTo(c.Param("id")).MustInt()
	name := c.Query("name")
	modifiedBy := c.Query("modified_by")

	valid := validation.Validation{}

	var state int = -1
	if arg := c.Query("state"); arg != "" {
		state = com.StrTo(arg).MustInt()
		valid.Range(state, 0, 1, "state").Message("状态只允许0或1")
	}
	valid.Required(id, "id").Message("ID不能为空")
	valid.Required(modifiedBy, "modified_by").Message("修改人不能为空")
	valid.MaxSize(modifiedBy, 100, "modified_by").Message("修改人最长为100字符")
	valid.MaxSize(name, 100, "name").Message("名称最长为100字符")

	code := e.INVALID_PARAMS
	if !valid.HasErrors() {
		code = e.SUCCESS
		if models.ExistTagByID(id) {
			data := make(map[string]interface{})
			data["modified_by"] = modifiedBy
			if name != "" {
				data["name"] = name
			}
			if state != -1 {
				data["state"] = state
			}
			models.EditTag(id, data)
		} else {
			code = e.ERROR_NOT_EXIST_TAG
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"code": code,
		"msg":  e.GetMsg(code),
		"data": make(map[string]string),
	})
}

// @summary 删除文章标签
// @Produce json
// @Param token query string true "Token"
// @Success 200 {string} json "{"code":200,"data":{},"msg":"ok"}"
// @Router /api/v1/tags/:id [DELETE]
func DeleteTag(c *gin.Context) {
	logging.Debug("xxxxxxxx fcf 删除tag")
	id := com.StrTo(c.Param("id")).MustInt()

	valid := validation.Validation{}
	valid.Min(id, 1, "id").Message("ID必须大于0")

	code := e.INVALID_PARAMS
	if !valid.HasErrors() {
		code = e.SUCCESS
		if models.ExistTagByID(id) {
			models.DeleteTag(id)
		} else {
			code = e.ERROR_NOT_EXIST_TAG
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"code": code,
		"msg":  e.GetMsg(code),
		"data": make(map[string]string),
	})
}