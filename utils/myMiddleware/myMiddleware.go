package myMiddleware

import (
    "../../struct/DB"

    "github.com/labstack/echo"
    _ "github.com/labstack/echo/middleware"

    "github.com/jinzhu/gorm"
    _ "github.com/mattn/go-sqlite3"
)

func CustomMiddleware(db *gorm.DB) echo.MiddlewareFunc {
    return func(next echo.HandlerFunc) echo.HandlerFunc {
        return func(c echo.Context) error {
            a := new(DB.Auth)
            a.UserId = "mstn_"
            a.PW = "admin"
            db.Create(&a)

            err := next(c)
            return err
        }
    }
}
