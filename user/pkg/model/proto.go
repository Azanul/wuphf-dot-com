package model

import (
	"wuphf.com/user/gen"
)

// MetadataToProto converts a User struct into a generated proto counterpart.
func UserToProto(m *User) *gen.User {
	return &gen.User{
		Id:    m.ID,
		Email: m.Email,
	}
}

// MetadataFromProto converts a generated proto counterpart into a User struct.
func UserFromProto(m *gen.User) *User {
	return &User{
		ID:    m.Id,
		Email: m.Email,
	}
}
