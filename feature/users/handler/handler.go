package handler

import (
	"fmt"
	middleware "hris-app-golang/app/middlewares"
	"hris-app-golang/feature/users"
	"hris-app-golang/helper"
	"net/http"
	"strconv"
	"strings"

	"github.com/labstack/echo/v4"
)

type UserHandler struct {
	userService users.UserServiceInterface
}

func New(service users.UserServiceInterface) *UserHandler {
	return &UserHandler{
		userService: service,
	}
}

func (handler *UserHandler) Add(c echo.Context) error {
	var input UserRequest
	roleId := middleware.ExtractTokenUserRoleId(c)
	fmt.Println("role id nya adalah:", roleId)
	if roleId == 4 {
		return c.JSON(http.StatusUnauthorized, helper.WebResponse(http.StatusUnauthorized, "operation failed, request resource not allowed", nil))
	}
	if roleId == 3 && (input.RoleID == 1 || input.RoleID == 2 || input.RoleID == 3) {
		return c.JSON(http.StatusUnauthorized, helper.WebResponse(http.StatusUnauthorized, "operation failed, request resource not allowed", nil))
	}
	if roleId == 2 && (input.RoleID == 1 || input.RoleID == 2) {
		return c.JSON(http.StatusUnauthorized, helper.WebResponse(http.StatusUnauthorized, "operation failed, request resource not allowed", nil))
	}
	errBind := c.Bind(&input)
	if errBind != nil {
		return c.JSON(http.StatusBadRequest, helper.WebResponse(http.StatusBadRequest, "operation failed, request resource not valid"+errBind.Error(), nil))
	}

	var fileName string
	file, header, errFile := c.Request().FormFile("profile_photo")

	if errFile != nil {
		if strings.Contains(errFile.Error(), "no such file") {
			fileName = "default.png"
		} else {
			return c.JSON(http.StatusBadRequest, helper.WebResponse(http.StatusBadRequest, "operation failed, request resource not valid "+errFile.Error(), nil))
		}
	}

	if fileName == "" {
		fileName = strings.ReplaceAll(header.Filename, " ", "_")
	}

	// errUp := helper.Uploader.UploadFile(file, header.Filename)

	// if errUp != nil {
	// 	return c.JSON(http.StatusInternalServerError, helper.WebResponse(http.StatusInternalServerError, "operation failed, internal server error", nil))
	// }
	input.Password = "qwerty"
	var userCore = UserRequestToCore(input)
	err := handler.userService.Add(userCore, file, fileName)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, helper.WebResponse(http.StatusInternalServerError, "operation failed, internal server error", nil))
	}
	return c.JSON(http.StatusOK, helper.WebResponse(http.StatusOK, "operation success, data added", nil))
}

func (handler *UserHandler) GetAll(c echo.Context) error {
	var pageConv, itemConv int
	var errPageConv, errItemConv error

	page := c.QueryParam("page")
	if page != "" {
		pageConv, errPageConv = strconv.Atoi(page)
		if errPageConv != nil {
			return c.JSON(http.StatusBadRequest, helper.WebResponse(http.StatusBadRequest, "operation failed, request resource not valid", nil))
		}
	}
	item := c.QueryParam("itemPerPage")
	if item != "" {
		itemConv, errItemConv = strconv.Atoi(item)
		if errItemConv != nil {
			return c.JSON(http.StatusBadRequest, helper.WebResponse(http.StatusBadRequest, "operation failed, request resource not valid", nil))
		}
	}

	search_name := c.QueryParam("searchName")
	role_id := middleware.ExtractTokenUserRoleId(c)
	fmt.Println("role id nya adalah:", role_id)
	division_id := middleware.ExtractTokenUserDivisionId(c)
	result, next, err := handler.userService.GetAll(role_id, division_id, uint(pageConv), uint(itemConv), search_name)
	if err != nil {
		if strings.Contains(err.Error(), "no row") {
			return c.JSON(http.StatusNotFound, helper.WebResponse(http.StatusNotFound, "operation failed, requested resource not found", nil))
		}
		return c.JSON(http.StatusInternalServerError, helper.WebResponse(http.StatusInternalServerError, "operation failed, internal server error", nil))
	}
	var userResp []UserResponseAll
	for _, value := range result {
		var user = UserCoreToResponseAll(value)
		userResp = append(userResp, user)
	}
	return c.JSON(http.StatusOK, helper.FindAllWebResponse(http.StatusOK, "success", userResp, next))
}

func (handler *UserHandler) Update(c echo.Context) error {
	id := c.Param("user_id")
	idConv, errConv := strconv.Atoi(id)
	if errConv != nil {
		return c.JSON(http.StatusBadRequest, helper.WebResponse(http.StatusBadRequest, "operation failed, request resource not valid", nil))
	}
	var input UserRequest
	errBind := c.Bind(&input)
	if errBind != nil {
		return c.JSON(http.StatusBadRequest, helper.WebResponse(http.StatusBadRequest, "operation failed, request resource not valid", nil))
	}
	var userCore = UserRequestToCore(input)

	var fileName string
	file, header, errFile := c.Request().FormFile("profil_photo")

	if errFile != nil {
		if strings.Contains(errFile.Error(), "no such file") {
			fileName = "default.jpg"
		}
		return c.JSON(http.StatusBadRequest, helper.WebResponse(http.StatusBadRequest, "operation failed, request resource not valid "+errFile.Error(), nil))
	}

	if fileName == "" {
		fileName = strings.ReplaceAll(header.Filename, " ", "_")
	}

	err := handler.userService.Update(uint(idConv), userCore, file, fileName)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, helper.WebResponse(http.StatusInternalServerError, "operation failed, internal server error", nil))
	}
	return c.JSON(http.StatusOK, helper.WebResponse(http.StatusOK, "success", nil))
}

func (handler *UserHandler) GetUserByID(c echo.Context) error {
	id := c.Param("user_id")
	role_id := middleware.ExtractTokenUserRoleId(c)
	user_id := middleware.ExtractTokenUserId(c)

	idConv, errConv := strconv.Atoi(id)
	if errConv != nil {
		return c.JSON(http.StatusBadRequest, helper.WebResponse(http.StatusBadRequest, "operation failed, request resource not valid", nil))
	}

	if role_id == 4 {
		if uint(idConv) != user_id {
			return c.JSON(http.StatusUnauthorized, helper.WebResponse(http.StatusUnauthorized, "operation failed, request resource not allowed", nil))
		}
	}

	result, err := handler.userService.GetById(uint(idConv))
	if err != nil {
		if strings.Contains(err.Error(), "no row affected") {
			return c.JSON(http.StatusNotFound, helper.WebResponse(http.StatusNotFound, "operation failed, requested resource not found", nil))
		}
		return c.JSON(http.StatusInternalServerError, helper.WebResponse(http.StatusInternalServerError, "operation failed, internal server error", nil))
	}

	if role_id == 2 {
		if result.RoleID == 1 {
			return c.JSON(http.StatusUnauthorized, helper.WebResponse(http.StatusUnauthorized, "operation failed, request resource not allowed", nil))
		}
	}

	resultResponse := UserCoreToResponse(result)

	return c.JSON(http.StatusOK, helper.WebResponse(http.StatusOK, "Success get user data", resultResponse))
}

func (handler *UserHandler) DeleteUser(c echo.Context) error {
	role_id := middleware.ExtractTokenUserRoleId(c)
	if role_id != 1 && role_id != 2 {
		return c.JSON(http.StatusUnauthorized, helper.WebResponse(http.StatusUnauthorized, "opeartion failed, request resource not allowed", nil))
	}
	id := c.Param("user_id")
	idConv, errConv := strconv.Atoi(id)
	if errConv != nil {
		return c.JSON(http.StatusBadRequest, helper.WebResponse(http.StatusBadRequest, "operation failed, request resource not valid", nil))
	}

	err := handler.userService.Delete(uint(idConv))
	if err != nil {
		if strings.Contains(err.Error(), "no row affected") {
			return c.JSON(http.StatusNotFound, helper.WebResponse(http.StatusNotFound, "operation failed, requested resource not found", nil))
		}
		return c.JSON(http.StatusInternalServerError, helper.WebResponse(http.StatusInternalServerError, "error delete data", nil))
	}
	return c.JSON(http.StatusOK, helper.WebResponse(http.StatusOK, "success delete data", nil))
}

func (handler *UserHandler) Login(c echo.Context) error {
	userInput := new(LoginRequest)
	errBind := c.Bind(&userInput)
	if errBind != nil {
		return c.JSON(http.StatusBadRequest, helper.WebResponse(http.StatusBadRequest, "error bind data. data not valid", nil))
	}
	dataLogin, token, err := handler.userService.Login(userInput.Email, userInput.Password)
	if err != nil {
		if strings.Contains(err.Error(), "validation") {
			return c.JSON(http.StatusBadRequest, helper.WebResponse(http.StatusBadRequest, err.Error(), nil))
		} else {
			return c.JSON(http.StatusInternalServerError, helper.WebResponse(http.StatusInternalServerError, "error login", nil))

		}
	}
	var response = LoginResponse{
		ID:       dataLogin.ID,
		Role:     dataLogin.Role.Name,
		Division: dataLogin.Division.Name,
		Token:    token,
	}
	return c.JSON(http.StatusOK, helper.WebResponse(http.StatusOK, "success login", response))
}

func (handler *UserHandler) GetAllManager(c echo.Context) error {
	result, err := handler.userService.GetAllManager()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, helper.WebResponse(http.StatusInternalServerError, "operation failed, internal server error", nil))
	}

	var userResp []ManagerResponse
	for _, value := range result {
		var resp = UserCoreToManagerResponse(value)
		userResp = append(userResp, resp)
	}
	return c.JSON(http.StatusOK, helper.WebResponse(http.StatusOK, "success", userResp))
}

func (handler *UserHandler) Dashboard(c echo.Context) error {
	var dashCore users.DashboardCore
	dashCore, err := handler.userService.GetDashboard()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, helper.WebResponse(http.StatusInternalServerError, "error getting employee count", nil))
	}

	var dashResp = DashboardCoreToResponse(dashCore)

	return c.JSON(http.StatusOK, helper.WebResponse(http.StatusOK, "success", dashResp))
}
