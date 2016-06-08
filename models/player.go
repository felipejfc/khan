package models

import "github.com/jinzhu/gorm"

//Player identifies uniquely one player in a given game
type Player struct {
	gorm.Model
	PlayerID string `gorm:"column:player_id;size:255"`
	Name     string `gorm:"size:2000"`
	GameID   string `gorm:"column:game_id;size:10"`
	Metadata string `sql:"type:JSONB NOT NULL DEFAULT '{}'::JSONB"`
}
