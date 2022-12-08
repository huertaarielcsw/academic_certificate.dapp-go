package service

import (
	"dapp/lib"
	"dapp/repo"
	"dapp/schema"
	"dapp/schema/dto"
	"dapp/schema/mapper"
	"dapp/schema/models"

	"github.com/kataras/iris/v12"
)

// region ======== SETUP =================================================================

// ISvcUser User request service interface
type ISvcUser interface {
	// user functions

	GetRolesSvc() (*dto.Pagination, *dto.Problem)
	GetUserSvc(userID int) (dto.UserResponse, *dto.Problem)
	GetUserByUsernameSvc(username string) (dto.UserResponse, *dto.Problem)
	GetUsersSvc(pagination *dto.Pagination) (*dto.Pagination, *dto.Problem)
	PutUserSvc(userID int, user dto.EditUserData) (dto.UserResponse, *dto.Problem)
	PostUserSvc(user dto.UserData) (dto.UserResponse, *dto.Problem)
	DeleteUserSvc(userID int) (dto.UserResponse, *dto.Problem)
	InvalidateUserSvc(userID int) (dto.UserResponse, *dto.Problem)
}

type svcUser struct {
	repoUser *repo.RepoUser
}

// endregion =============================================================================

// NewSvcUserReqs instantiate the User request services
func NewSvcUserReqs(repoUser *repo.RepoUser) ISvcUser {
	return &svcUser{repoUser}
}

// region ======== METHODS ======================================================

func (s *svcUser) GetRolesSvc() (*dto.Pagination, *dto.Problem) {
	res, err := (*s.repoUser).GetRoles()
	if err != nil {
		return nil, lib.NewProblem(iris.StatusExpectationFailed, schema.ErrBuntdb, err.Error())
	}
	return &dto.Pagination{Rows: res}, nil
}

func (s *svcUser) GetUserSvc(userID int) (dto.UserResponse, *dto.Problem) {
	res, err := (*s.repoUser).GetUser(userID)
	if err != nil {
		return dto.UserResponse{}, lib.NewProblem(iris.StatusExpectationFailed, schema.ErrBuntdb, err.Error())
	}
	return mapper.MapModelUser2DtoUserResponse(res), nil
}

func (s *svcUser) GetUserByUsernameSvc(username string) (dto.UserResponse, *dto.Problem) {
	res, err := (*s.repoUser).GetUserByUsername(username)
	if err != nil {
		return dto.UserResponse{}, lib.NewProblem(iris.StatusExpectationFailed, schema.ErrBuntdb, err.Error())
	}
	return mapper.MapModelUser2DtoUserResponse(res), nil
}

func (s *svcUser) GetUsersSvc(pagination *dto.Pagination) (*dto.Pagination, *dto.Problem) {
	res, err := (*s.repoUser).GetUsers(pagination)
	if err != nil {
		return nil, lib.NewProblem(iris.StatusExpectationFailed, schema.ErrBuntdb, err.Error())
	}
	var usersResponse []dto.UserResponse
	items := res.Rows.([]models.User)
	for i := 0; i < len(items); i++ {
		usersResponse = append(usersResponse, mapper.MapModelUser2DtoUserResponse(items[i]))
	}
	res.Rows = usersResponse
	return res, nil
}

func (s *svcUser) PutUserSvc(userID int, user dto.EditUserData) (dto.UserResponse, *dto.Problem) {
	userInDB, err := s.repoUser.GetUser(userID)
	if err != nil {
		return dto.UserResponse{}, lib.NewProblem(iris.StatusExpectationFailed, schema.ErrBuntdb, err.Error())
	}
	if user.Username != "" {
		userInDB.Username = user.Username
	}
	if user.Passphrase != "" {
		passphraseEncoded, _ := lib.Checksum("SHA256", []byte(user.Passphrase))
		userInDB.Passphrase = passphraseEncoded
	}
	if user.FirstName != "" {
		userInDB.FirstName = user.FirstName
	}
	if user.LastName != "" {
		userInDB.LastName = user.LastName
	}
	if user.Email != "" {
		userInDB.Email = user.Email
	}
	if user.Role != "" {
		userInDB.Role = user.Role
	}
	resUser, err := s.repoUser.UpdateUser(userID, userInDB)
	if err != nil {
		return dto.UserResponse{}, lib.NewProblem(iris.StatusExpectationFailed, schema.ErrBuntdb, err.Error())
	}
	return mapper.MapModelUser2DtoUserResponse(resUser), nil
}

func (s *svcUser) PostUserSvc(user dto.UserData) (dto.UserResponse, *dto.Problem) {
	passphraseEncoded, _ := lib.Checksum("SHA256", []byte(user.Passphrase))
	user.Passphrase = passphraseEncoded
	modelUser := mapper.MapUserData2ModelUser(0, user)
	resUser, err := s.repoUser.AddUser(modelUser)
	if err != nil {
		return dto.UserResponse{}, lib.NewProblem(iris.StatusExpectationFailed, schema.ErrBuntdb, err.Error())
	}
	return mapper.MapModelUser2DtoUserResponse(resUser), nil
}

func (s *svcUser) DeleteUserSvc(userID int) (dto.UserResponse, *dto.Problem) {
	user, err := s.repoUser.RemoveUser(userID)
	if err != nil {
		return dto.UserResponse{}, lib.NewProblem(iris.StatusExpectationFailed, schema.ErrBuntdb, err.Error())
	}
	return mapper.MapModelUser2DtoUserResponse(user), nil
}

func (s *svcUser) InvalidateUserSvc(userID int) (dto.UserResponse, *dto.Problem) {
	user, err := s.repoUser.InvalidateUser(userID)
	if err != nil {
		return dto.UserResponse{}, lib.NewProblem(iris.StatusExpectationFailed, schema.ErrBuntdb, err.Error())
	}
	return mapper.MapModelUser2DtoUserResponse(user), nil
}
