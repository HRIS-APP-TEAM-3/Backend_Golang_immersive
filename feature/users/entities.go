package users

import (
	divisionCore "hris-app-golang/feature/divisions"
	roleCore "hris-app-golang/feature/roles"
	"time"
)

type UserCore struct {
	ID           uint
	FirstName    string
	LastName     string
	Email        string
	PhoneNumber  string
	Password     string
	Address      string
	ProfilePhoto string
	UserLeadID   uint
	RoleID       uint
	DivisionID   uint
	Role         roleCore.RoleCore
	Division     divisionCore.DivisionCore
	UserImport   UserImportantData
	UserEdu      []UserEducationData
	CreatedAt    time.Time
	UpdatedAt    time.Time
	// UserLead     *UserCore
}

type UserImportantData struct {
	ID              uint
	UserID          uint
	BirthPlace      string
	BirthDate       string
	EmergencyName   string
	EmergencyStatus string
	EmergencyPhone  string
	Npwp            string
	Bpjs            string
	Religion        string
	Gender          string
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

type UserEducationData struct {
	ID           uint
	UserID       uint
	Name         string
	StartYear    string
	GraduateYear string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

type UserDataInterface interface {
	Insert(input UserCore) error
	SelectAll(role_id, division_id uint) ([]UserCore, error)
	SelectById(id uint) (UserCore, error)
	Update(id uint, input UserCore) error
	Delete(id uint) error
	Login(email, password string) (UserCore, error)
}

type UserServiceInterface interface {
	Add(input UserCore) error
	GetAll(role_id, division_id uint) ([]UserCore, error)
	GetById(id uint) (UserCore, error)
	Update(id uint, input UserCore) error
	Delete(id uint) error
	Login(email, password string) (UserCore, string, error)
}
