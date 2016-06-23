// khan
// https://github.com/topfreegames/khan
//
// Licensed under the MIT license:
// http://www.opensource.org/licenses/mit-license
// Copyright © 2016 Top Free Games <backend@tfgco.com>

package models

import (
	"time"

	"gopkg.in/gorp.v1"
)

// Player identifies uniquely one player in a given game
type Player struct {
	ID        int    `db:"id"`
	GameID    string `db:"game_id"`
	PublicID  string `db:"public_id"`
	Name      string `db:"name"`
	Metadata  string `db:"metadata"`
	CreatedAt int64  `db:"created_at"`
	UpdatedAt int64  `db:"updated_at"`
}

// PreInsert populates fields before inserting a new player
func (p *Player) PreInsert(s gorp.SqlExecutor) error {
	p.CreatedAt = time.Now().UnixNano() / 1000000
	p.UpdatedAt = p.CreatedAt
	return nil
}

// PreUpdate populates fields before updating a player
func (p *Player) PreUpdate(s gorp.SqlExecutor) error {
	p.UpdatedAt = time.Now().UnixNano() / 1000000
	return nil
}

// GetPlayerByID returns a player by id
func GetPlayerByID(db DB, id int) (*Player, error) {
	obj, err := db.Get(Player{}, id)
	if err != nil || obj == nil {
		return nil, &ModelNotFoundError{"Player", id}
	}

	player := obj.(*Player)
	return player, nil
}

// GetPlayerByPublicID returns a player by their public id
func GetPlayerByPublicID(db DB, gameID string, publicID string) (*Player, error) {
	var player Player
	err := db.SelectOne(&player, "SELECT * FROM players WHERE game_id=$1 AND public_id=$2", gameID, publicID)
	if err != nil || &player == nil {
		return nil, &ModelNotFoundError{"Player", publicID}
	}
	return &player, nil
}

// CreatePlayer creates a new player
func CreatePlayer(db DB, gameID string, publicID string, name string, metadata string) (*Player, error) {
	player := &Player{
		GameID:   gameID,
		PublicID: publicID,
		Name:     name,
		Metadata: metadata,
	}
	err := db.Insert(player)
	if err != nil {
		return nil, err
	}
	return player, nil
}

// UpdatePlayer updates an existing player
func UpdatePlayer(db DB, gameID, publicID, name, metadata string) (*Player, error) {
	player, err := GetPlayerByPublicID(db, gameID, publicID)

	if err != nil {
		return nil, err
	}

	player.Name = name
	player.Metadata = metadata

	_, err = db.Update(player)

	if err != nil {
		return nil, err
	}

	return player, nil
}

// GetPlayerDetails returns detailed information about a player and their memberships
func GetPlayerDetails(db DB, gameID, publicID string) (map[string]interface{}, error) {
	query := `
	SELECT
		p.name PlayerName, p.metadata PlayerMetadata, p.public_id PlayerPublicID,
		p.created_at PlayerCreatedAt, p.updated_at PlayerUpdatedAt,
		m.membership_level MembershipLevel,
		m.approved MembershipApproved, m.denied MembershipDenied, m.banned MembershipBanned,
		c.public_id ClanPublicID, c.name ClanName, c.metadata ClanMetadata,
		r.name RequestorName, r.public_id RequestorPublicID, r.metadata RequestorMetadata,
		m.created_at MembershipCreatedAt,
		m.updated_at MembershipUpdatedAt,
		m.deleted_at MembershipDeletedAt,
		d.name DeletedByName, d.public_id DeletedByPublicID
	FROM players p
	LEFT OUTER JOIN memberships m on m.player_id = p.id
	INNER JOIN clans c on c.id=m.clan_id
	LEFT OUTER JOIN players d on d.id=m.deleted_by
	LEFT OUTER JOIN players r on r.id=m.requestor_id
	WHERE p.game_id=$1 and p.public_id=$2`

	var details []playerDetailsDAO
	_, err := db.Select(&details, query, gameID, publicID)
	if err != nil {
		return nil, err
	}

	if len(details) == 0 {
		return nil, &ModelNotFoundError{"Player", publicID}
	}

	result := make(map[string]interface{})

	result["name"] = details[0].PlayerName
	result["metadata"] = details[0].PlayerMetadata
	result["publicID"] = details[0].PlayerPublicID
	result["createdAt"] = details[0].PlayerCreatedAt
	result["updatedAt"] = details[0].PlayerUpdatedAt

	if details[0].MembershipLevel.Valid {
		// Player has memberships
		result["memberships"] = make([]JSON, len(details))

		approved := []JSON{}
		denied := []JSON{}
		banned := []JSON{}
		pending := []JSON{}

		clanFromDetail := func(clanDetail playerDetailsDAO) JSON {
			return JSON{
				"publicID": nullOrString(clanDetail.ClanPublicID),
				"name":     nullOrString(clanDetail.ClanName),
			}
		}

		for index, detail := range details {
			result["memberships"].([]JSON)[index] = detail.Serialize()

			ma := nullOrBool(detail.MembershipApproved)
			md := nullOrBool(detail.MembershipDenied)
			mb := nullOrBool(detail.MembershipBanned)

			clanDetail := clanFromDetail(detail)
			switch {
			case !ma && !md && !mb:
				pending = append(pending, clanDetail)
			case ma:
				approved = append(approved, clanDetail)
			case md:
				denied = append(denied, clanDetail)
			case mb:
				banned = append(banned, clanDetail)
			}
		}

		result["clans"] = JSON{
			"approved": approved,
			"denied":   denied,
			"banned":   banned,
			"pending":  pending,
		}

	} else {
		result["memberships"] = []JSON{}
	}

	return result, nil
}
