package user

import (
	"context"

	myHelper "github.com/livegoplayer/go_helper"
	user "go_user_rpc/user/grpc"
)

//专门的go rpc server代码，为了和基础的service代码解耦，主要是避免输入输出的对象化，那样写很难受

type UserServiceServer struct {
}

//以下函数是给endpoints调用的，最后会被安装到对应的成员变量上，所以一定要是grpc.Handler 类型的
func (u *UserServiceServer) Login(_ context.Context, request *user.LoginRequest) (loginResponse *user.LoginResponse, err error) {
	loginResponse = &user.LoginResponse{}
	uid, userSession, token, err := UserServiceInstance.Login(request.GetUserName(), request.GetPassword(), request.GetHost())
	if err != nil {
		loginResponse.Msg = err.Error()
		loginResponse.ErrorCode = 1
		return
	}
	loginData := &user.LoginData{}
	loginData.Uid = uid
	loginData.UserSession = utilUserSessionToResponse(userSession)
	loginData.Token = token

	loginResponse.Data = loginData

	return
}

func (u *UserServiceServer) CheckUserStatus(_ context.Context, checkUserStatusRequest *user.CheckUserStatusRequest) (checkUserStatusResponse *user.CheckUserStatusResponse, err error) {
	checkUserStatusResponse = &user.CheckUserStatusResponse{}
	isLogin, tokenStr, err := UserServiceInstance.CheckLoginStatus(checkUserStatusRequest.GetToken(), checkUserStatusRequest.GetHost())
	if err != nil {
		checkUserStatusResponse.Msg = err.Error()
		checkUserStatusResponse.ErrorCode = 1
		return
	}
	checkUserStatusData := &user.CheckUserStatusData{}
	checkUserStatusData.IsLogin = isLogin
	checkUserStatusData.Token = tokenStr

	checkUserStatusResponse.Data = checkUserStatusData

	return
}

func (u *UserServiceServer) Register(_ context.Context, request *user.RegisterRequest) (registerResponse *user.RegisterResponse, err error) {
	registerResponse = &user.RegisterResponse{}
	uid, err := UserServiceInstance.Register(request.GetToken(), request.GetHost())
	if err != nil {
		registerResponse.Msg = err.Error()
		registerResponse.ErrorCode = 1
		return
	}
	registerData := &user.RegisterData{}
	registerData.Uid = uid

	registerResponse.Data = registerData

	return
}

func (u *UserServiceServer) AddUser(_ context.Context, request *user.AddUserRequest) (addUserResponse *user.AddUserResponse, err error) {
	addUserResponse = &user.AddUserResponse{}
	uid, err := UserServiceInstance.AddUser(request.GetUserName(), request.GetPassword(), request.GetOperationUid())
	if err != nil {
		addUserResponse.Msg = err.Error()
		addUserResponse.ErrorCode = 1
		return
	}
	addUserData := &user.AddUserData{}
	addUserData.Uid = uid

	addUserResponse.Data = addUserData
	return
}

func (u *UserServiceServer) DelUser(_ context.Context, request *user.DelUserRequest) (delUserResponse *user.DelUserResponse, err error) {
	delUserResponse = &user.DelUserResponse{}
	success, err := UserServiceInstance.DelUser(request.GetUid(), request.GetOperationUid())
	if err != nil {
		delUserResponse.Msg = err.Error()
		delUserResponse.ErrorCode = 1
		return
	}
	delUserData := &user.DelUserData{}
	delUserData.Success = success

	delUserResponse.Data = delUserData
	return
}

func (u *UserServiceServer) CheckUserAuthority(_ context.Context, request *user.CheckUserAuthorityRequest) (checkUserAuthorityResponse *user.CheckUserAuthorityResponse, err error) {
	checkUserAuthorityResponse = &user.CheckUserAuthorityResponse{}
	exists := UserServiceInstance.CheckUserAuthority(request.GetUid(), request.GetAuthorityId())

	checkUserAuthorityData := &user.CheckUserAuthorityData{}
	checkUserAuthorityData.Exist = exists

	checkUserAuthorityResponse.Data = checkUserAuthorityData
	return
}
func (u *UserServiceServer) GetUserAuthorityList(_ context.Context, request *user.GetUserAuthorityListRequest) (getUserAuthorityListResponse *user.GetUserAuthorityListResponse, err error) {
	getUserAuthorityListResponse = &user.GetUserAuthorityListResponse{}
	userAuthorityList := UserServiceInstance.GetUserAuthorityList(request.GetUid())

	getUserAuthorityListData := &user.GetUserAuthorityListData{}
	getUserAuthorityListData.UserAuthorityList = userAuthorityList

	getUserAuthorityListResponse.Data = getUserAuthorityListData
	return
}

func (u *UserServiceServer) GetAuthorityList(_ context.Context, _ *user.GetAuthorityListRequest) (getAuthorityListResponse *user.GetAuthorityListResponse, err error) {
	getAuthorityListResponse = &user.GetAuthorityListResponse{}

	authorityList := UserServiceInstance.GetAuthorityList()

	getAuthorityListData := &user.GetAuthorityListData{}
	getAuthorityListData.UserAuthorityList = authorityList

	getAuthorityListResponse.Data = getAuthorityListData
	return
}

func (u *UserServiceServer) AddUserRole(_ context.Context, request *user.AddUserRoleRequest) (addUserRoleResponse *user.AddUserRoleResponse, err error) {
	addUserRoleResponse = &user.AddUserRoleResponse{}

	success := UserServiceInstance.AddUserRole(request.GetUid(), request.GetRoleId(), request.GetOperationUid())

	addUserRoleData := &user.AddUserRoleData{}
	addUserRoleData.Success = success

	addUserRoleResponse.Data = addUserRoleData
	return
}

func (u *UserServiceServer) DelUserRole(_ context.Context, request *user.DelUserRoleRequest) (delUserRoleResponse *user.DelUserRoleResponse, err error) {
	delUserRoleResponse = &user.DelUserRoleResponse{}
	success := UserServiceInstance.AddUserRole(request.GetUid(), request.GetRoleId(), request.GetOperationUid())

	delUserRoleData := &user.DelUserRoleData{}
	delUserRoleData.Success = success

	delUserRoleResponse.Data = delUserRoleData
	return
}

func (u *UserServiceServer) GetUserRoleList(_ context.Context, request *user.GetUserRoleListRequest) (getUserRoleListResponse *user.GetUserRoleListResponse, err error) {
	getUserRoleListResponse = &user.GetUserRoleListResponse{}
	userRoleList := UserServiceInstance.GetUserRoleList(request.GetUid())

	getUserRoleListData := &user.GetUserRoleListData{}
	getUserRoleListData.UserRoleList = userRoleList

	getUserRoleListResponse.Data = getUserRoleListData
	return
}

func (u *UserServiceServer) GetRoleList(_ context.Context, _ *user.GetRoleListRequest) (getRoleListResponse *user.GetRoleListResponse, err error) {
	getRoleListResponse = &user.GetRoleListResponse{}
	roleList := UserServiceInstance.GetRoleList()

	getRoleListData := &user.GetRoleListData{}
	getRoleListData.UserRoleList = roleList

	getRoleListResponse.Data = getRoleListData
	return
}

func utilUserSessionToResponse(session *myHelper.UserSession) *user.UserSessions {
	responseUserSession := &user.UserSessions{
		Uid:               int32(session.Uid),
		Username:          session.UserName,
		UserRoleList:      session.UserRoleList,
		UserAuthorityList: session.UserAuthorityList,
		AddDatetime:       session.AddDatetime,
		UpdateDatetime:    session.UpdateDatetime,
	}

	return responseUserSession
}

func responseToUtilUserSession(session *user.UserSessions) *myHelper.UserSession {
	responseUserSession := &myHelper.UserSession{
		Uid:               session.Uid,
		UserName:          session.Username,
		UserRoleList:      session.UserRoleList,
		UserAuthorityList: session.UserAuthorityList,
		AddDatetime:       session.AddDatetime,
		UpdateDatetime:    session.UpdateDatetime,
	}

	return responseUserSession
}

////以下是中间件部分
//type ServiceMiddleware func(UserServiceServer) UserServiceServer
//
//// loggingMiddleware Make a new type
//// that contains Service interface and logger instance
//type loggingMiddleware struct {
//	UserServiceServer
//	logger logrus.Logger
//}
//
//// LoggingMiddleware make logging middleware
//func LoggingMiddleware(logger logrus.Logger) ServiceMiddleware {
//	return func(next UserServiceServer) UserServiceServer {
//		return loggingMiddleware{next, logger}
//	}
//}
