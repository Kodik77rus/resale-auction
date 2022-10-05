package models

import "net"

type Context struct {
	Ip        net.IP `json:"string"`
	UserAgent string `json:"user_agent"`
}
