package models

type Context struct {
	Ip        string `valid:"required,ip" json:"ip"`
	UserAgent string `valid:"required,ascii" json:"user_agent"`
}
