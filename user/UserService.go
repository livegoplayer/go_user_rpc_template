package user

import (
	"errors"
	"time"

	myHelper "github.com/livegoplayer/go_helper"
	"github.com/livegoplayer/go_user_rpc/model"
	sessionHelper "github.com/livegoplayer/go_user_rpc/redisHelper"
)

//定义一个接口，service包含的方法
type UserServiceInterface interface {
	Login(userName string, password string) (uid int32, userSession *myHelper.UserSession, tokenStr string, err error)
	CheckLoginStatus(token string, host string) (isLogin bool, tokenStr string, err error)
	Register(userName string, password string) (uid int32, err error)
	AddUser(userName string, password string, operationUid int32) (uid int32, err error)
	DelUser(uid int32, operationUid int32) (success bool, err error)

	CheckUserAuthority(uid int32, authorityId int32) bool
	GetUserAuthorityList(uid int32) map[int32]string
	GetAuthorityList() map[int32]string

	AddUserRole(uid int32, roleId int32, operationUid int32) bool
	DelUserRole(uid int32, roleId int32, operationUid int32) bool
	GetRoleList() map[int32]string
	GetUserRoleList() map[int32]string
}

type UserService struct{}

var UserServiceInstance *UserService

func init() {
	UserServiceInstance = &UserService{}
}

func getUserSession(uid int32) (session *myHelper.UserSession, exist bool) {
	key := getUserLoginStatusSessionKey(uid)
	res := sessionHelper.GetRedisKey(key)
	val := res.Val()
	if val != "" {
		exist = true
		myHelper.JsonDecode(myHelper.StringToBytes(val), session)
	}

	return
}

func setUserSession(uid int32, session *myHelper.UserSession) bool {
	key := getUserLoginStatusSessionKey(uid)
	_ = sessionHelper.SetRedisKey(key, myHelper.JsonEncode(session), time.Hour*24)
	return true
}

func getUserLoginStatusSessionKey(uid int32) string {
	key := "user_login_status_session_uid_" + string(uid)
	return key
}

//用户注册逻辑
func (userService *UserService) Register(userName string, password string) (uid int32, err error) {
	isRecordFound, err := model.CheckUserName(userName)
	if isRecordFound {
		err = errors.New("该用户名已经存在")
		return
	}

	return model.AddNewUser(userName, password, 0)
}

//管理员添加用户逻辑
func (userService *UserService) AddUser(userName string, password string, operationUid int32) (uid int32, err error) {
	return model.AddNewUser(userName, password, operationUid)
}

//管理员添加用户逻辑
func (userService *UserService) DelUser(uid int32, operationUid int32) (success bool, err error) {
	return model.DelUser(uid, operationUid)
}

func (userService *UserService) AddUserRole(uid int32, roleId int32, operationUid int32) bool {
	return model.AddUserRole(uid, roleId, operationUid)
}

//不给用
func (userService *UserService) DelUserRole(uid int32, roleId int32, operationUid int32) bool {
	return model.DelUserRole(uid, roleId, operationUid)
}

func (userService *UserService) GetRoleList() (roleList map[int32]string) {
	roles := model.GetRoleList()
	roleList = make(map[int32]string)
	for _, role := range roles {
		roleList[role.ID] = role.Name
	}

	return roleList
}

func (userService *UserService) GetUserRoleList(uid int32) (userRoleList map[int32]string) {
	roles := model.GetUserRoleList(uid)
	userRoleList = make(map[int32]string)
	for _, role := range roles {
		userRoleList[role.ID] = role.Name
	}
	return userRoleList
}

func (userService *UserService) Login(userName string, password string) (uid int32, userSession *myHelper.UserSession, tokenStr string, err error) {
	//第一步获取用户信息
	isRecordFound, user, err := model.CheckUserPassword(userName, password)
	if err != nil {
		return
	}

	if !isRecordFound {
		err = errors.New("用户名或者密码错误")
		return
	}

	//获取detail信息
	userSession, err = model.GetUserDetailInfo(user.ID)
	if err != nil {
		return
	}

	//设置session
	setUserSession(uid, userSession)

	//生成新的token
	tokenStr, err = myHelper.CreateToken(userSession)

	return
}

func (userService *UserService) CheckLoginStatus(token string, host string) (isLogin bool, tokenStr string, err error) {
	claims, err := myHelper.ParseToken(token, host)
	if err != nil {
		//如果token过期了
		if err.Error() == "jwt过期" || err.Error() == "host错误" {
			userSession := claims.UserSession
			//检查session是否过期
			userSessionNew, exsit := getUserSession(userSession.Uid)
			if exsit {
				// 重新根据当前的session生成token
				// 生成新的token
				isLogin = true
				tokenStr, err = myHelper.CreateToken(userSessionNew)
			}
		}

		return
	}

	isLogin = true

	return
}

func (userService *UserService) CheckUserAuthority(uid int32, authorityId int32) (exists bool) {
	exists, _ = model.CheckUserAuthority(uid, authorityId)
	return
}

func (userService *UserService) GetUserAuthorityList(uid int32) (userAuthorityList map[int32]string) {
	userAuthorityList, _ = model.GetUserAuthorityList(uid)
	return
}

func (userService *UserService) GetAuthorityList() (authorityList map[int32]string) {
	authorityList, _ = model.GetAuthorityList()
	return
}
