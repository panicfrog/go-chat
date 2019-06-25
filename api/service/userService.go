package service

import (
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"gochat/db"
	"gochat/toolUtils"
	"gochat/variable"
	"log"
)

type UserParams struct {
	 Username string `json:"username" binding:"required"`
	 Password string `json:"password" binding:"required"`
}

// 注册
func Register(c *gin.Context) {
	var userParams UserParams
	if err := c.ShouldBindJSON(&userParams); err != nil {
		sendParamError(err.Error(), c)
		return
	}

	var user = db.User{}
	errDB := db.DB.Where(&db.User{UserName:userParams.Username}).First(&user).Error
	if errDB != nil && errDB != gorm.ErrRecordNotFound  {
		log.Println(errDB)
		sendServerInternelError(errDB.Error(), c)
	}

	if errDB == nil {
		sendFail("账号("+userParams.Username+")已被注册", c)
		return
	}

	user = db.User{UserName: userParams.Username, Passwd:userParams.Password}

	if err := db.DB.Create(&user).Error; err != nil {
		log.Println(err)
		return
	}
	sendSuccess("注册成功", nil, c)
}

// 登录
func Login(c *gin.Context) {
	var userParams UserParams
	if err := c.ShouldBindJSON(&userParams); err != nil {
		sendParamError(err.Error(), c)
		return
	}
	var user = db.User{}
	err := db.DB.Where(&db.User{UserName:userParams.Username, Passwd:userParams.Password}).First(&user).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		log.Println(err)
		sendServerInternelError(err.Error(), c)
		return
	}

	if err != nil && err == gorm.ErrRecordNotFound {
		sendFail("用户不存在或密码错误", c)
		return
	}
	platform := c.Request.Header.Get("platform")
	p := variable.StoP(platform)
	authUser := toolUtils.AuthUser{user, p}
	b, e := authUser.CreateToken()
	if e != nil {
		log.Println(e)
		sendServerInternelError(e.Error(), nil)
		return
	}
	token, e := toolUtils.EncodeToken(string(b))
	if e != nil {
		sendServerInternelError("加密出错", c)
		return
	}
	sendSuccess("登录成功", token, c)
}

func Friends(c *gin.Context) {
	userName := getUserId(c)
	user, ok := verifyAndGetUser(c, userName)
	if !ok {
		return
	}
	fridends := []db.User{}
	dbErr := db.DB.Model(&user).Association("Friends").Find(&fridends).Error
	if dbErr != nil {
		sendServerInternelError(dbErr.Error(), c)
		return
	}
	sendSuccess("成功", fridends, c)
}

type AddFriendsParams struct {
	Friends []string		`json:"friends"`
}

func AddFriends(c *gin.Context) {
	var params AddFriendsParams
	if err := c.ShouldBindJSON(&params); err != nil {
		sendFail(err.Error(), c)
		return
	}

	if len(params.Friends) <= 0 {
		sendFail("friends 字段为空数组", c)
		return
	}

	userName := getUserId(c)
	user, ok := verifyAndGetUser(c, userName)
	if !ok {
		return
	}

	// 去重
	fids := map[string]bool{}
	var isSelf = false
	for _,f := range params.Friends {
		if f == user.UserName {
			isSelf = true
			break
		}
		fids[f] = true
	}

	// 限制不能添加自己为好友
	if isSelf {
		sendFail("不能添加自己为好友", c)
		return
	}

	ok = true
	memes := []db.User{}
	// 判断所有用户id是否有效
	for id,_ := range fids {
		u, ok := verifyAndGetUser(c, id)
		if !ok {
			break
		}
		memes = append(memes, u)
	}
	if !ok {
		return
	}

	// 相互添加好友
	dbErr := db.DB.Model(&user).Association("Friends").Append(memes).Error
	if dbErr != nil {
		sendServerInternelError(dbErr.Error(), c)
		return
	}

	for _,u := range memes {
		if err := db.DB.Model(&u).Association("Friends").Append([]db.User{user}).Error; err != nil {
			sendServerInternelError(err.Error(), c)
			goto END
		}
	}

	sendSuccess("添加成功", nil, c)
	END:
}

type RemoveFriendsParams struct {
	Friends []string		`json:"friends"`
}
func RemoveFriends(c *gin.Context) {
	var params RemoveFriendsParams
	if err := c.ShouldBindJSON(&params); err != nil {
		sendFail(err.Error(), c)
		return
	}

	if len(params.Friends) <= 0 {
		sendFail("friends 字段为空数组", c)
		return
	}

	userName := getUserId(c)
	user, ok := verifyAndGetUser(c, userName)
	if !ok {
		return
	}

	// 去重
	fids := map[string]bool{}
	for _,f := range params.Friends {
		if f == user.UserName {
			break
		}
		fids[f] = true
	}

	ok = true
	memes := []db.User{}
	// 判断所有用户id是否有效
	for id := range fids {
		u, ok := verifyAndGetUser(c, id)
		if !ok {
			return
		}
		memes = append(memes, u)
	}
	if !ok {
		return
	}

	var (
		friends []db.User
	)
	dbErr := db.DB.Model(&user).Association("Friends").Find(&friends).Error
	if dbErr != nil {
		sendServerInternelError(dbErr.Error(), c)
		return
	}

	// 判断是否在好友列表中
	var isFriend = true
	for _,_u := range memes {
		isIn := false
		for _,u := range friends {
			if u.UserName == _u.UserName {
				isIn = true
				break
			}
		}
		if !isIn {
			isFriend = false
			break
		}
	}
	if !isFriend  {
		sendFail("只有好友才能删除好友", c)
		return
	}
	dbErr = db.DB.Model(&user).Association("Friends").Delete(memes).Error
	if dbErr != nil {
		sendServerInternelError(dbErr.Error(), c)
		return
	}
	for _, u := range memes {
		dbErr = db.DB.Model(&u).Association("Friends").Delete(user).Error
		if dbErr != nil {
			sendServerInternelError(dbErr.Error(), c)
			goto END
		}
	}
	sendSuccess("删除成功", nil , c)
	END:
}


type ResultRoom struct {
	Owned string 			`json:"owned"`
	RoomId string 			`json:"room_id"`
}

func newRoom(room *db.Room) (rr *ResultRoom) {
	rr = &ResultRoom{Owned:room.Owned, RoomId:room.RoomId}
	return
}

func Rooms(c *gin.Context) {
	userName := getUserId(c)
	user, ok := verifyAndGetUser(c, userName)
	if !ok {
		return
	}

	var (
		rooms []db.Room
		resultRooms []ResultRoom
	)
	dbErr := db.DB.Model(&user).Association("Rooms").Find(&rooms).Error
	if dbErr != nil {
		sendServerInternelError(dbErr.Error(), c)
		return
	}

	for _, r := range rooms {
		resultRooms = append(resultRooms, *newRoom(&r))
	}

	sendSuccess("请求成功", resultRooms, c)
}

