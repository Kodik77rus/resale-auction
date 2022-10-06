package models

import "net"

type Context struct {
	Ip        net.IP `valid:"required,ip" json:"ip"`
	UserAgent string `valid:"required,alpha" json:"user_agent"`
}
