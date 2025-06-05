package models

import (
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// UserRole represents the different types of users in the system
type UserRole string

const (
	RoleFamille       UserRole = "famille"
	RoleEnseignant    UserRole = "enseignant"
	RoleAdministrator UserRole = "administrator"
)

// Base User model - represents the base user class from the diagram
type User struct {
	ID          uint           `json:"id" gorm:"primaryKey"`
	Username    string         `json:"username" gorm:"unique;not null"`
	Password    string         `json:"-" gorm:"not null"`
	Email       string         `json:"email" gorm:"unique;not null"`
	PhoneNumber string         `json:"phone_number"`
	Role        UserRole       `json:"role" gorm:"not null"`
	IsActive    bool           `json:"is_active" gorm:"default:true"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`

	// Relationships
	Addresses []Address  `json:"addresses,omitempty" gorm:"foreignKey:UserID"`
	Payments  []Payment  `json:"payments,omitempty" gorm:"foreignKey:UserID"`
	Resources []Resource `json:"resources,omitempty" gorm:"many2many:user_resources;"`
}

// Famille model - inherits from User
type Famille struct {
	UserID     uint   `json:"user_id" gorm:"primaryKey"`
	FamilyName string `json:"family_name"`

	// Relationships
	User     User      `json:"user" gorm:"foreignKey:UserID"`
	Missions []Mission `json:"missions,omitempty" gorm:"foreignKey:FamilleID"`
	Courses  []Course  `json:"courses,omitempty" gorm:"foreignKey:FamilleID"`
	Options  []Option  `json:"options,omitempty" gorm:"foreignKey:FamilleID"`
}

// Enseignant model - inherits from User
type Enseignant struct {
	UserID         uint   `json:"user_id" gorm:"primaryKey"`
	Specialization string `json:"specialization"`
	Qualifications string `json:"qualifications"`

	// Relationships
	User     User      `json:"user" gorm:"foreignKey:UserID"`
	Missions []Mission `json:"missions,omitempty" gorm:"foreignKey:EnseignantID"`
	Courses  []Course  `json:"courses,omitempty" gorm:"foreignKey:EnseignantID"`
	Reports  []Report  `json:"reports,omitempty" gorm:"foreignKey:EnseignantID"`
	Options  []Option  `json:"options,omitempty" gorm:"foreignKey:EnseignantID"`
	Offers   []Offer   `json:"offers,omitempty" gorm:"many2many:enseignant_offers;"`
}

// Administrator model - inherits from User
type Administrator struct {
	UserID uint `json:"user_id" gorm:"primaryKey"`

	// Relationships
	User      User       `json:"user" gorm:"foreignKey:UserID"`
	Reports   []Report   `json:"validated_reports,omitempty" gorm:"foreignKey:ValidatedByID"`
	Offers    []Offer    `json:"created_offers,omitempty" gorm:"foreignKey:CreatedByID"`
	Resources []Resource `json:"managed_resources,omitempty" gorm:"foreignKey:ManagedByID"`
}

// User methods
func (u *User) Login() error {
	// Login logic will be implemented in controllers
	return nil
}

func (u *User) Logout() error {
	// Logout logic will be implemented in controllers
	return nil
}

func (u *User) UpdateProfile() error {
	// Update profile logic will be implemented in controllers
	return nil
}

func (u *User) ManageContacts() error {
	// Contact management logic will be implemented in controllers
	return nil
}

// Famille methods
func (f *Famille) ConsultTeachers() ([]Enseignant, error) {
	// Logic to consult available teachers
	return nil, nil
}

func (f *Famille) PlanNextSession() error {
	// Logic to plan next session
	return nil
}

func (f *Famille) ViewPayments() ([]Payment, error) {
	// Logic to view payments
	return nil, nil
}

func (f *Famille) StopMission(missionID uint) error {
	// Logic to stop a mission
	return nil
}

func (f *Famille) ViewCourses() ([]Course, error) {
	// Logic to view courses
	return nil, nil
}

func (f *Famille) SearchTeachers(criteria map[string]interface{}) ([]Enseignant, error) {
	// Logic to search teachers based on criteria
	return nil, nil
}

func (f *Famille) SelectTeacher(teacherID uint) error {
	// Logic to select a teacher
	return nil
}

func (f *Famille) ValidateSession(courseID uint) error {
	// Logic to validate a session
	return nil
}

func (f *Famille) DeclineSession(courseID uint) error {
	// Logic to decline a session
	return nil
}

// Enseignant methods
func (e *Enseignant) CompleteProfile() error {
	// Logic to complete profile
	return nil
}

func (e *Enseignant) ViewOffers() ([]Offer, error) {
	// Logic to view available offers
	return nil, nil
}

func (e *Enseignant) ReserveOption(optionID uint) error {
	// Logic to reserve an option
	return nil
}

func (e *Enseignant) PlanCourses() error {
	// Logic to plan courses
	return nil
}

func (e *Enseignant) DeclareSession(courseID uint) error {
	// Logic to declare a session
	return nil
}

func (e *Enseignant) ProvideReport(missionID uint, content string) error {
	// Logic to provide a report
	return nil
}

func (e *Enseignant) StopMission(missionID uint) error {
	// Logic to stop a mission
	return nil
}

func (e *Enseignant) DeclareHours(courseID uint, hours float64) error {
	// Logic to declare hours worked
	return nil
}

func (e *Enseignant) SelectStudent(familleID uint) error {
	// Logic to select a student/family
	return nil
}

// Administrator methods
func (a *Administrator) ManageAccount(userID uint) error {
	// Logic to manage user accounts
	return nil
}

func (a *Administrator) ManageStudents() error {
	// Logic to manage students
	return nil
}

func (a *Administrator) ManageTeachers() error {
	// Logic to manage teachers
	return nil
}

func (a *Administrator) ManageOffers() error {
	// Logic to manage offers
	return nil
}

func (a *Administrator) ConsultReports() ([]Report, error) {
	// Logic to consult reports
	return nil, nil
}

func (a *Administrator) ValidateReports(reportID uint) error {
	// Logic to validate reports
	return nil
}

func (a *Administrator) EditProfiles(userID uint) error {
	// Logic to edit user profiles
	return nil
}

func (a *Administrator) ProvideSupport() error {
	// Logic to provide support
	return nil
}

// Request/Response structures for API
type UserCreateRequest struct {
	Username       string   `json:"username" binding:"required,min=3,max=50"`
	Password       string   `json:"password" binding:"required,min=6"`
	Email          string   `json:"email" binding:"required,email"`
	PhoneNumber    string   `json:"phone_number"`
	Role           UserRole `json:"role" binding:"required"`
	FamilyName     string   `json:"family_name,omitempty"`    // For Famille
	Specialization string   `json:"specialization,omitempty"` // For Enseignant
	Qualifications string   `json:"qualifications,omitempty"` // For Enseignant
}

type UserLoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type UserUpdateRequest struct {
	Username       string `json:"username,omitempty"`
	Email          string `json:"email,omitempty"`
	PhoneNumber    string `json:"phone_number,omitempty"`
	FamilyName     string `json:"family_name,omitempty"`
	Specialization string `json:"specialization,omitempty"`
	Qualifications string `json:"qualifications,omitempty"`
}

// BeforeCreate hash le mot de passe avant de créer l'utilisateur
func (u *User) BeforeCreate(tx *gorm.DB) error {
	if u.Password != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
		if err != nil {
			return err
		}
		u.Password = string(hashedPassword)
	}
	return nil
}

// CheckPassword vérifie si le mot de passe fourni correspond au hash stocké
func (u *User) CheckPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	return err == nil
}

// GetSafeUser retourne une version sécurisée de l'utilisateur sans le mot de passe
func (u *User) GetSafeUser() User {
	safeUser := *u
	safeUser.Password = ""
	return safeUser
}
