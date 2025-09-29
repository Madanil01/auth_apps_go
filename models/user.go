package models


type User struct {
	Id           uint
	Name         string
	Email        string
	Password     []byte
	Alamat       string
	Username     string
	NoHp         string
	ConfirmEmail bool
	IsAdmin      bool `gorm:"default:false"`
	IsCashier    bool `gorm:"default:false"`

	// Relasi ke Role
	RoleID uint  `json:"role_id"`
	Role   Role  `json:"role" gorm:"foreignKey:RoleID"`
}

type Role struct {
	ID            uint           `json:"id" gorm:"primaryKey"`
	Name          string         `json:"name"`
	AccessRights  []AccessRight  `json:"access_rights" gorm:"foreignKey:RoleID"`
	Users         []User         `json:"users" gorm:"foreignKey:RoleID"`
}

type AccessRight struct {
	ID         uint      `json:"id" gorm:"primaryKey"`
	Name       string    `json:"name"`
	Access     string    `json:"access"`
	RoleID     uint      `json:"role_id"`
	PageAppsID uint      `json:"page_apps_id"`
	PageApps   PageApps  `json:"page_apps" gorm:"foreignKey:PageAppsID"`
}

type PageApps struct {
	ID           uint       `json:"id" gorm:"primaryKey"`
	AppsName     string     `json:"apps_name"`
	Name         string     `json:"name"`
	Path         string     `json:"path"`
	OptionalPath *string    `json:"optional_path"`
	Icon         string     `json:"icon"`
	Status       bool       `json:"status"`

	// Self referencing (parent-child)
	ParentID *uint      `json:"parent_id"`
	Parent   *PageApps  `json:"parent" gorm:"foreignKey:ParentID"`
	Children []PageApps `json:"children" gorm:"foreignKey:ParentID"`
}
