package user

import (
    "../../struct/DB"
    "../../utils/sha256"
    "../../utils/randomString"

    "os"
    "fmt"
    "net/http"

    "github.com/labstack/echo"
    _ "github.com/labstack/echo/middleware"

    "github.com/jinzhu/gorm"
    _ "github.com/mattn/go-sqlite3"
)

func transaction(db *gorm.DB, data interface{}, action string) error {
    tx := db.Begin()
    defer func() {
        if r := recover(); r != nil {
            tx.Rollback()
        }
    }()

    if err := tx.Error; err != nil {
        return err
    }

    if action == "save" {
        if err := tx.Save(data).Error; err != nil {
            tx.Rollback()
            return err
        }
    } else if action == "create" {
        if err := tx.Create(data).Error; err != nil {
            tx.Rollback()
            return err
        }
    }

    if err := tx.Commit().Error; err != nil {
        return err
    }

    return nil
}

func GenerateAuthCode(db *gorm.DB) echo.HandlerFunc {
    return func(c echo.Context) error {
        userid := c.Param("userid")

        code := new(DB.AuthCode)
        code.UserId = userid

        db.Where("user_id = ?", userid).First(&code)
        code.Code = randomString.RandString(8)

        transaction(db, code, "save")

        return c.HTML(http.StatusOK, code.Code)
    }
}

func Create(db *gorm.DB) echo.HandlerFunc {
    return func(c echo.Context) error {
        user := new(DB.Auth)
        if err := c.Bind(user); err != nil {
            fmt.Fprintln(os.Stderr, err)
            return err
        }

        if user.UserId == "" || user.PW == "" {
            return c.NoContent(http.StatusBadRequest)
        }

        isUsedUserId := new(DB.Auth)
        db.Where("user_id = ?", user.UserId).First(&isUsedUserId)
        if isUsedUserId.UserId == user.UserId {
            return c.NoContent(http.StatusBadRequest)
        }

        transaction(db, user, "create")

        return c.NoContent(http.StatusOK)
    }
}

func Login(db *gorm.DB, isLoginDB *gorm.DB) echo.HandlerFunc {
    return func(c echo.Context) error {
        formData := new(DB.Auth)
        if err := c.Bind(formData); err != nil {
            fmt.Fprintln(os.Stderr, err)
            return err
        }

        // TODO クエリの単一化
        user := new(DB.Auth)
        db.Where("user_id = ?", formData.UserId).First(&user)

        code := new(DB.AuthCode)
        db.Where("user_id = ?", formData.UserId).First(&code)

        auth := sha256.Sha256Sum([]byte(user.UserId + user.PW + code.Code))

        if formData.PW == auth {
            login := new(DB.LoginStatus)
            isLoginDB.Where("user_id = ?", formData.UserId).First(&login)
            login.UserId = formData.UserId
            login.IsLogin = true

            transaction(isLoginDB, login, "save")

            return c.NoContent(http.StatusOK)
        } else {
            return c.NoContent(http.StatusUnauthorized)
        }
    }
}
