package models

import "net"

type Context struct {
	Ip        net.IP `json:"ip"`
	UserAgent string `json:"user_agent"`
}
