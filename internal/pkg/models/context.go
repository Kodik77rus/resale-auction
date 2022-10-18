package models

type Context struct {
	Ip        string `valid:"required,ip" json:"ip"`
	UserAgent string `valid:"required,halfwidth" json:"user_agent"`
}
