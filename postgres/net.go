package postgres

import "github.com/aodin/sol/types"

func Cidr() types.BaseType {
	return types.Base("CIDR")
}

func Inet() types.BaseType {
	return types.Base("INET")
}

func Macaddr() types.BaseType {
	return types.Base("MACADDR")
}
