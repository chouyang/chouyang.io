package models

import (
	"time"
)

// File model represents an actual file or an article.
type File struct {
	ID         uint      `json:"id"         gorm:"not null; primaryKey; autoIncrement; comment:file id"`
	Name       string    `json:"name"       gorm:"not null; size:64; index; comment:file name"`
	Path       string    `json:"path"       gorm:"not null; size:256; unique; comment:file path"`
	Size       int64     `json:"size"       gorm:"not null; default:0; comment:file size"`
	Mime       string    `json:"mime"       gorm:"not null; size:32; comment:file mime"`
	Hash       string    `json:"hash"       gorm:"not null; size:32; comment:file hash"`
	Permission string    `json:"permission" gorm:"not null; size:16; comment:file permission"`
	Content    string    `json:"content"    gorm:"not null; type:longtext; comment:file content"`
	CreatedBy  uint      `json:"created_by" gorm:"not null; default:0; foreignKey:User; comment:file created by"`
	CreatedAt  time.Time `json:"created_at" gorm:"not null; autoUpdateTime; comment:file created at"`
	UpdatedAt  time.Time `json:"updated_at" gorm:"not null; autoCreateTime; comment:file updated at"`
}

// Tree model represents the path of an actual file or an article.
type Tree struct {
	ID        uint      `json:"id"         gorm:"not null; primaryKey; autoIncrement; comment:file id"`
	Name      string    `json:"name"`
	Path      string    `json:"path"`
	Extra     string    `json:"extra"`
	CreatedAt time.Time `json:"created_at" gorm:"not null; autoUpdateTime; comment:file created at"`
	UpdatedAt time.Time `json:"updated_at" gorm:"not null; autoCreateTime; comment:file updated at"`
	Files     []*File   `json:"files"`
	Trees     []*Tree   `json:"trees"      gorm:"foreignKey:ParentID; comment:tree list"`
}
