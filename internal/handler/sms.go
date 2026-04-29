package handler

import (
	"net/http"
	"net/url"
	"sms-dashboard/internal/database"
	"sms-dashboard/internal/model"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type SMSHandler struct{}

func NewSMSHandler() *SMSHandler {
	return &SMSHandler{}
}

func (h *SMSHandler) Receive(c *gin.Context) {
	var req struct {
		Content     string `json:"content" form:"content" binding:"required"`
		SendTime    string `json:"sendTime" form:"sendTime"`
		OrgContent  string `json:"org_content" form:"org_content"`
		Sign        string `json:"sign" form:"sign"`
		Timestamp   string `json:"timestamp" form:"timestamp"`
		ReceiveTime string `json:"receive_time" form:"receive_time"`
		Sender      string `json:"sender" form:"sender"`
		Phone       string `json:"phone" form:"phone"`
		Device      string `json:"device" form:"device"`
	}

	// For JSON requests, we can also use ShouldBindJSON
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request format",
			"details": err.Error(),
		})
		return
	}
	//校验Sign - 已禁用
	// cfg := config.LoadConfig()
	// timestamp := req.Timestamp
	// secret := cfg.Secret
	// stringToSign := timestamp + "\n" + secret
	// hash := hmac.New(sha256.New, []byte(secret))
	// hash.Write([]byte(stringToSign))
	// signData := hash.Sum(nil)
	// base64Sign := base64.StdEncoding.EncodeToString(signData)
	// sign := url.QueryEscape(base64Sign)
	// if sign != req.Sign {
	// 	c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid Sign"})
	// 	return
	// }
	sms := model.SMS{
		Content:  req.Content,
		SendTime: req.SendTime,
		Sender:   req.Sender,
		Phone:    req.Phone,
		Device:   req.Device,
	}

	if sms.SendTime == "" {
		sms.SendTime = time.Now().Format("2006-01-02 15:04:05")
	}

	// 对 content 进行 URLDecode 处理
	if decodedContent, err := url.QueryUnescape(sms.Content); err == nil {
		sms.Content = decodedContent
	}

	if err := database.DB.Create(&sms).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "id": sms.ID})
}

func (h *SMSHandler) List(c *gin.Context) {
	var smsList []model.SMS
	var total int64
	pageStr := c.DefaultQuery("page", "1")
	pageSizeStr := c.DefaultQuery("pageSize", "10")

	page, _ := strconv.Atoi(pageStr)
	pageSize, _ := strconv.Atoi(pageSizeStr)
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}

	database.DB.Model(&model.SMS{}).Count(&total)

	offset := (page - 1) * pageSize
	if err := database.DB.Order("created_at desc").Offset(offset).Limit(pageSize).Find(&smsList).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch SMS"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": smsList,
		"pagination": gin.H{
			"page":     page,
			"pageSize": pageSize,
			"total":    total,
		},
	})
}

func (h *SMSHandler) Delete(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID is required"})
		return
	}

	if err := database.DB.Delete(&model.SMS{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete SMS"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success"})
}

// GroupedSMS 分组短信结构
type GroupedSMS struct {
	Device  string       `json:"device"`
	Total   int64        `json:"total"`
	HasMore bool         `json:"hasMore"`
	SMSList []model.SMS  `json:"smsList"`
}

func (h *SMSHandler) GroupedList(c *gin.Context) {
	limit := 10

	// 获取最新10条记录
	var latest []model.SMS
	database.DB.Order("created_at desc").Limit(limit).Find(&latest)

	// 获取所有不同的机型(包括空值)
	var devices []string
	database.DB.Model(&model.SMS{}).Distinct("device").Pluck("device", &devices)

	// 为每个机型获取最新的10条记录
	var groups []GroupedSMS
	for _, device := range devices {
		var total int64
		var query = database.DB.Model(&model.SMS{})
		
		if device == "" || device == "null" {
			query = query.Where("device IS NULL OR device = ''")
		} else {
			query = query.Where("device = ?", device)
		}
		query.Count(&total)

		var smsList []model.SMS
		listQuery := database.DB.Order("created_at desc").Limit(limit)
		if device == "" || device == "null" {
			listQuery = listQuery.Where("device IS NULL OR device = ''")
		} else {
			listQuery = listQuery.Where("device = ?", device)
		}
		listQuery.Find(&smsList)

		// 如果设备名为空,显示为"未知机型"
		displayDevice := device
		if device == "" || device == "null" {
			displayDevice = "未知机型"
		}

		groups = append(groups, GroupedSMS{
			Device:  displayDevice,
			Total:   total,
			HasMore: total > int64(limit),
			SMSList: smsList,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"latest": latest,
		"groups": groups,
	})
}

func (h *SMSHandler) LoadMore(c *gin.Context) {
	device := c.Query("device")
	offsetStr := c.DefaultQuery("offset", "10")
	limitStr := c.DefaultQuery("limit", "10")

	if device == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "device parameter is required"})
		return
	}

	offset, _ := strconv.Atoi(offsetStr)
	limit, _ := strconv.Atoi(limitStr)

	var total int64
	var query = database.DB.Model(&model.SMS{})
	
	// 处理"未知机型"
	if device == "未知机型" {
		query = query.Where("device IS NULL OR device = ''")
	} else {
		query = query.Where("device = ?", device)
	}
	query.Count(&total)

	var smsList []model.SMS
	listQuery := database.DB.Order("created_at desc").Offset(offset).Limit(limit)
	if device == "未知机型" {
		listQuery = listQuery.Where("device IS NULL OR device = ''")
	} else {
		listQuery = listQuery.Where("device = ?", device)
	}
	listQuery.Find(&smsList)

	hasMore := int64(offset+limit) < total

	c.JSON(http.StatusOK, gin.H{
		"device":  device,
		"smsList": smsList,
		"hasMore": hasMore,
	})
}
