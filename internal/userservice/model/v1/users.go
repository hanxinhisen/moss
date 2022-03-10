package v1

import (
	"github.com/marmotedu/component-base/pkg/auth"
	metav1 "github.com/marmotedu/component-base/pkg/meta/v1"
	"github.com/marmotedu/component-base/pkg/util/idutil"
	"gorm.io/gorm"
	"time"
)

type User struct {
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Status            int    `json:"status" gorm:"column:status" validate:"omitempty"`
	Nickname          string `json:"nickname" gorm:"column:nickname" validate:"required,min=1,max=30"`

	// Required: true
	Password string `json:"password,omitempty" gorm:"column:password" validate:"required"`

	// Required: true
	Email string `json:"email" gorm:"column:email" validate:"required,email,min=1,max=100"`

	Phone string `json:"phone" gorm:"column:phone" validate:"omitempty"`

	IsAdmin int `json:"isAdmin,omitempty" gorm:"column:isAdmin" validate:"omitempty"`

	TotalPolicy int64 `json:"totalPolicy" gorm:"-" validate:"omitempty"`

	LoginedAt time.Time `json:"loginedAt,omitempty" gorm:"column:loginedAt"`
}

type UserList struct {
	metav1.ListMeta `json:",inline"`
	Items           []*User `json:"items"`
}

func (u *User) TableName() string {
	return "user"
}

func (u *User) Compare(pwd string) (err error) {
	err = auth.Compare(u.Password, pwd)

	return
}
func (u *User) AfterCreate(tx *gorm.DB) error {
	u.InstanceID = idutil.GetInstanceID(u.ID, "user-")

	return tx.Save(u).Error
}
