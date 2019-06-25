package service

import (
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/satori/go.uuid"
	"gochat/db"
	"gochat/variable"
	"net/http"
	"strings"
)

func getUserId(c *gin.Context) (userId string) {
	userName, exited := c.Get(variable.AuthKey)
	if !exited {
		sendHTTPError(c, http.StatusUnauthorized, "用户未登录")
		return
	}
	userId = userName.(string)
	return
}

// 验证用户
func verifyAndGetUser(c *gin.Context, userName string) (user db.User, ok bool) {
	u := db.User{}
	errDB := db.DB.Where(&db.User{UserName:userName}).First(&u).Error
	if errDB != nil && errDB != gorm.ErrRecordNotFound {
		sendServerInternelError(errDB.Error(), c)
		return user ,false
	}
	if errDB != nil && errDB == gorm.ErrRecordNotFound {
		sendFail("用户不存在", c)
		return user, false
	}
	return u, true
}

// 判断是否是一个房间的管理员
func verifyRoomManager(c *gin.Context, userName string, roomId string) (room db.Room, isManager bool) {
	room = db.Room{}
	isManager = false
	dbErr := db.DB.Where(&db.Room{RoomId:roomId}).First(&room).Error
	if dbErr != nil {
		sendServerInternelError(dbErr.Error(), c)
		goto END
	}

	dbErr = db.DB.Model(&room).Association("Managers").Find(&room.Managers).Error
	if dbErr != nil {
		sendServerInternelError(dbErr.Error(), c)
		goto END
	}

	for _,m := range room.Managers {
		if m.UserName == userName {
			isManager = true
			goto END
		}
	}
	END:
		return
}


// 创建房间
func CreateRoom(c *gin.Context) {
	userName := getUserId(c)
	user, ok := verifyAndGetUser(c, userName)
	if !ok {
		return
	}
	_uuid := strings.ReplaceAll(uuid.NewV4().String(),"-","")
	room := db.Room{Owned:userName, RoomId:_uuid, Users:[]db.User{user}, Managers:[]db.User{user}}

	d := db.DB.Create(&room)
	if d.Error != nil {
		sendServerInternelError(d.Error.Error(), c)
		return
	}
	sendSuccess("创建成功", _uuid, c)
}

type GetRoomInfoParams struct {
	RoomId string `json:"room_id"`
}
// 获取房间信息
func GetRoomInfo(c *gin.Context) {
	var getRoomInfoParams GetRoomInfoParams
	if err := c.ShouldBindJSON(&getRoomInfoParams); err != nil {
		sendParamError(err.Error(), c)
		return
	}
	room := db.Room{}
	dbErr := db.DB.Where(&db.Room{RoomId:getRoomInfoParams.RoomId}).First(&room).Error
	if dbErr != nil && dbErr != gorm.ErrRecordNotFound {
		sendServerInternelError(dbErr.Error(), c)
		return
	}
	if dbErr != nil && dbErr == gorm.ErrRecordNotFound {
		sendFail(dbErr.Error(), c)
		return
	}
	dbErr = db.DB.Model(&room).Association("Users").Find(&room.Users).Error
	if dbErr != nil {
		sendServerInternelError(dbErr.Error(), c)
		return
	}
	dbErr = db.DB.Model(&room).Association("Managers").Find(&room.Managers).Error
	if dbErr != nil {
		sendServerInternelError(dbErr.Error(), c)
		return
	}
	sendSuccess("获取成功", room, c)
}

type AddMemeberParams struct {
	Users []string `json:"users"`
	RoomId string	`json:"room_id"`
} 

// 添加群成员
func AddMemberToToom(c *gin.Context) {
	var addMemeParams AddMemeberParams
	if err := c.ShouldBindJSON(&addMemeParams); err != nil {
		sendParamError(err.Error(), c)
		return
	}

	// 判断参数中是否有好友
	if len(addMemeParams.Users) <= 0 {
		sendFail("没有任何的好友", c)
		return
	}

	// 过滤重复的id
	fids := map[string]bool{}
	for _, f := range addMemeParams.Users {
		fids[f] = true
	}

	// 判断用户是否登录
	userName := getUserId(c)
	user, ok := verifyAndGetUser(c, userName)
	if !ok {
		return
	}

	ok = true
	memes := []db.User{}
	// 判断所有用户id是否有效
	for id,_ := range fids {
		u, ok := verifyAndGetUser(c, id)
		if !ok {
			return
		}
		memes = append(memes, u)
	}
	if !ok {
		return
	}

	//判断用户是否是房间的管理员
	room, isManager := verifyRoomManager(c, user.UserName, addMemeParams.RoomId)

	if !isManager {
		sendFail("您不是管理员，没有权限添加成员", c)
		return
	}

	// 关联查找房间用户
	dbErr := db.DB.Model(&room).Association("Users").Find(&room.Users).Error
	if dbErr != nil {
		sendServerInternelError(dbErr.Error(), c)
		return
	}

	// 过滤掉重复的
	for _, u := range room.Users {
		var isOk = true
		for _, _u := range memes {
			if _u.UserName == u.UserName {
				isOk = false
				break
			}
		}
		if u.UserName == user.UserName {
			isOk = false
		}
		if isOk {
			memes = append(memes, u)
		}
	}

	room.Users = memes
	// 添加好友
	dbErr = db.DB.Save(&room).Error
	if dbErr != nil {
		sendServerInternelError(dbErr.Error(), c)
		return
	}
	sendSuccess("添加成功", nil, c)
}

type DeleMemeParams struct {
	Users []string 	`json:"users"`
	RoomId string	`json:"room_id"`
} 

// 删除群成员
func DeleteMemeberFromRoom(c *gin.Context) {
	var deleMemeParams DeleMemeParams
	if err := c.ShouldBindJSON(&deleMemeParams); err != nil {
		sendParamError(err.Error(), c)
		return
	}

	// 判断参数中是否有好友
	if len(deleMemeParams.Users) <= 0 {
		sendFail("没有任何的好友", c)
		return
	}

	// 过滤重复的id
	fids := map[string]bool{}
	for _, f := range deleMemeParams.Users {
		fids[f] = true
	}

	// 判断用户是否登录
	userName := getUserId(c)
	user, ok := verifyAndGetUser(c, userName)
	if !ok {
		return
	}

	// 判断所有用户id是否有效
	ok = true
	memes := []db.User{}
	for id,_ := range fids {
		u, ok := verifyAndGetUser(c, id)
		if !ok {
			return
		}
		memes = append(memes, u)
	}
	if !ok {
		return
	}

	//判断用户是否是房间的管理员
	room, isManager := verifyRoomManager(c, user.UserName, deleMemeParams.RoomId)

	if !isManager {
		sendFail("您不是管理员，没有权限添加成员", c)
		return
	}

	dbErr := db.DB.Model(&room).Association("Users").Delete(memes).Error
	if dbErr != nil {
		sendServerInternelError(dbErr.Error(), c)
		return
	}
	sendSuccess("已经移出群组", nil, c)
}

type EmpowerManagerParams struct {
	Managers []string 	`json:"managers"`
	RoomId string		`json:"room_id"`
}

// 授权管理员
func EmpowerManager(c *gin.Context) {
	var empowerManagerParams EmpowerManagerParams
	if err := c.ShouldBindJSON(&empowerManagerParams); err != nil {
		sendParamError(err.Error(), c)
		return
	}

	// 判断参数中是否有好友
	if len(empowerManagerParams.Managers) <= 0 {
		sendFail("没有任何的管理员", c)
		return
	}

	// 过滤重复的id
	mids := map[string]bool{}
	for _, f := range empowerManagerParams.Managers {
		mids[f] = true
	}

	// 判断用户是否登录
	userName := getUserId(c)
	user, ok := verifyAndGetUser(c, userName)
	if !ok {
		return
	}

	//判断用户是否是房间的管理员
	room, isManager := verifyRoomManager(c, user.UserName, empowerManagerParams.RoomId)
	if !(isManager && room.Owned == userName) {
		sendFail("你没有权限设置管理员", c)
		return
	}

	// 判断所有用户id是否有效
	ok = true
	memes := []db.User{}
	for id,_ := range mids {
		u, ok := verifyAndGetUser(c, id)
		if !ok {
			return
		}
		memes = append(memes, u)
	}
	if !ok {
		return
	}

	// 验证是否都是该房间的成员
	dbErr := db.DB.Model(&room).Association("Users").Find(&room.Users).Error
	if dbErr != nil {
		sendServerInternelError(dbErr.Error(), c)
		return
	}
	var isInRoom = true
	for _, _u := range memes {
		var isIn = false
		for _, u := range room.Users {
			if u.UserName == _u.UserName {
				isIn = true
				break
			}
		}
		if !isIn {
			isInRoom = false
			break
		}
	}

	if !isInRoom {
		sendFail("管理员必须是该房间成员", c)
		return
	}

	dbErr = db.DB.Model(&room).Association("Managers").Append(memes).Error
	if dbErr != nil {
		sendServerInternelError(dbErr.Error(), c)
		return
	}
	sendSuccess("设置成功", nil, c)
}

type CallOffManagerParams struct {
	Managers []string 	`json:"managers"`
	RoomId   string 	`json:"room_id"`
}

// 取消管理员
func CallOffManager(c *gin.Context) {
	var callOffManagerParams CallOffManagerParams
	if err := c.ShouldBindJSON(&callOffManagerParams); err != nil {
		sendParamError(err.Error(), c)
		return
	}

	// 判断参数中是否是空数组
	if len(callOffManagerParams.Managers) <= 0 {
		sendFail("managers 为空数组", c)
		return
	}

	// 过滤重复的id
	mids := map[string]bool{}
	for _, f := range callOffManagerParams.Managers {
		mids[f] = true
	}

	// 判断用户是否登录
	userName := getUserId(c)
	user, ok := verifyAndGetUser(c, userName)
	if !ok {
		return
	}

	//判断用户是否是房间的管理员
	room, isManager := verifyRoomManager(c, user.UserName, callOffManagerParams.RoomId)
	if !(isManager && room.Owned == userName) {
		sendFail("你没有权限设置管理员", c)
		return
	}

	// 判断是否是是否是管理员
	var callOfMangers = []db.User{}
	var allIsManager = true
	for mid, _ := range mids {
		var isIn = false
		for _, u := range room.Managers {
			if u.UserName == mid {
				callOfMangers = append(callOfMangers, u)
				isIn = true
				break
			}
		}
		if !isIn {
			allIsManager = false
			break
		}
	}

	if !allIsManager {
		sendFail("所有人必须是管理员才能取消", c)
		return
	}

	dbErr := db.DB.Model(&room).Association("Managers").Delete(callOfMangers).Error
	if dbErr != nil {
		sendServerInternelError(dbErr.Error(), c)
		return
	}
	sendSuccess("取消成功", nil, c)
}

type TransferOwnerParams struct {
	User string 	`json:"user"`
	RoomId string 	`json:"room_id"`
}

// 转让群
func TransferOwner(c *gin.Context) {
	var transferOwnerParams TransferOwnerParams
	if err := c.ShouldBindJSON(&transferOwnerParams); err != nil {
		sendParamError(err.Error(), c)
		return
	}

	// 判断用户是否登录
	userName := getUserId(c)
	user, ok := verifyAndGetUser(c, userName)
	if !ok {
		return
	}

	room := db.Room{}
	dbErr := db.DB.Where(&db.Room{RoomId:transferOwnerParams.RoomId}).First(&room).Error
	if dbErr != nil {
		sendServerInternelError(dbErr.Error(), c)
		return
	}
	if room.Owned != user.UserName {
		sendFail("你没有权限操作", c)
		return
	}

	// 检查用户是否存在
	_, ok = verifyAndGetUser(c, transferOwnerParams.User)
	if !ok {
		sendFail("不存在的用户", c)
		return
	}

	// 修改拥有者
	room.Owned = transferOwnerParams.User
	dbErr = db.DB.Save(&room).Error
	if dbErr != nil {
		sendServerInternelError(dbErr.Error(), nil)
		return
	}
	sendSuccess("转让成功", nil , c)
}