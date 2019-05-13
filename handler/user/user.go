package user

import (
    "../../struct/DB"
    "../../utils/sha256"

    "os"
    "fmt"
    "net/http"

    "github.com/labstack/echo"
    _ "github.com/labstack/echo/middleware"

    "github.com/jinzhu/gorm"
    _ "github.com/mattn/go-sqlite3"
)

func GenerateAuthCode(db *gorm.DB) echo.HandlerFunc {
    return func(c echo.Context) error {
        userid := c.Param("userid")

        code := new(DB.AuthCode)
        code.UserId = userid
        //TODO ランダム生成
        code.Code = "code"

        isUsedUserId := new(DB.Auth)
        db.Where("user_id = ?", userid).First(&isUsedUserId)

        if isUsedUserId.UserId == userid {
            db.Save(&code)
        } else {
            tx := db.Begin()
            defer func() {
                if r := recover(); r != nil {
                    tx.Rollback()
                }
            }()

            if err := tx.Error; err != nil {
                return err
            }

            if err := tx.Create(&code).Error; err != nil {
                tx.Rollback()
                return err
            }

            if err := tx.Commit().Error; err != nil {
                return err
            }
        }

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

        tx := db.Begin()
        defer func() {
            if r := recover(); r != nil {
                tx.Rollback()
            }
        }()

        if err := tx.Error; err != nil {
            return err
        }

        if err := tx.Create(&user).Error; err != nil {
            tx.Rollback()
            return err
        }

        if err := tx.Commit().Error; err != nil {
            return err
        } else {
            return c.NoContent(http.StatusOK)
        }
    }
}

func Login(db *gorm.DB) echo.HandlerFunc {
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
            return c.NoContent(http.StatusOK)
        } else {
            return c.NoContent(http.StatusUnauthorized)
        }
    }
}
