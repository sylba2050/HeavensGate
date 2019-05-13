package DB

import (
    "github.com/jinzhu/gorm"
    _ "github.com/mattn/go-sqlite3"
)

type Auth struct {
    gorm.Model
    UserId string `json:"userid" form:"userid" query:"userid"`
    PW string `json:"pw" form:"pw" query:"pw"`
}

type AuthCode struct {
    gorm.Model
    UserId string `json:"userid" form:"userid" query:"userid"`
    Code string `json:"code" form:"code" query:"code"`
}

type LoginStatus struct {
    gorm.Model
    UserId string `json:"userid" form:"userid" query:"userid"`
    IsLogin bool `json:"is_login" form:"is_login" query:"is_login"`
}

