// khan
// https://github.com/topfreegames/khan
//
// Licensed under the MIT license:
// http://www.opensource.org/licenses/mit-license
// Copyright © 2016 Top Free Games <backend@tfgco.com>

package models

import (
	"fmt"

	"github.com/Pallinder/go-randomdata"
	"github.com/bluele/factory-go/factory"
	"github.com/satori/go.uuid"
)

//PlayerFactory is responsible for constructing test player instances
var PlayerFactory = factory.NewFactory(
	&Player{},
).Attr("GameID", func(args factory.Args) (interface{}, error) {
	return uuid.NewV4().String(), nil
}).Attr("PublicID", func(args factory.Args) (interface{}, error) {
	return uuid.NewV4().String(), nil
}).Attr("Name", func(args factory.Args) (interface{}, error) {
	return randomdata.FullName(randomdata.RandomGender), nil
}).Attr("Metadata", func(args factory.Args) (interface{}, error) {
	return "{}", nil
})

//ClanFactory is responsible for constructing test clan instances
var ClanFactory = factory.NewFactory(
	&Clan{},
).Attr("GameID", func(args factory.Args) (interface{}, error) {
	return uuid.NewV4().String(), nil
}).Attr("PublicID", func(args factory.Args) (interface{}, error) {
	return uuid.NewV4().String(), nil
}).Attr("Name", func(args factory.Args) (interface{}, error) {
	return randomdata.FullName(randomdata.RandomGender), nil
}).Attr("Metadata", func(args factory.Args) (interface{}, error) {
	return "{}", nil
})

//MembershipFactory is responsible for constructing test membership instances
var MembershipFactory = factory.NewFactory(
	&Membership{},
).SeqInt("GameID", func(n int) (interface{}, error) {
	return fmt.Sprintf("game-%d", n), nil
})

//GetClanWithMemberships returns a clan filled with the number of memberships specified
func GetClanWithMemberships(
	db DB, numberOfMemberships int, gameID string, clanPublicID string,
) (*Clan, *Player, []*Player, []*Membership, error) {
	owner := PlayerFactory.MustCreateWithOption(map[string]interface{}{
		"GameID": gameID,
	}).(*Player)
	err := db.Insert(owner)
	if err != nil {
		return nil, nil, nil, nil, err
	}

	var players []*Player

	for i := 0; i < numberOfMemberships; i++ {
		player := PlayerFactory.MustCreateWithOption(map[string]interface{}{
			"GameID": owner.GameID,
		}).(*Player)
		err = db.Insert(player)
		if err != nil {
			return nil, nil, nil, nil, err
		}
		players = append(players, player)
	}

	clan := ClanFactory.MustCreateWithOption(map[string]interface{}{
		"GameID":   owner.GameID,
		"PublicID": clanPublicID,
		"OwnerID":  owner.ID,
		"Metadata": "{\"x\": 1}",
	}).(*Clan)
	err = db.Insert(clan)
	if err != nil {
		return nil, nil, nil, nil, err
	}

	var memberships []*Membership

	for i := 0; i < numberOfMemberships; i++ {
		membership := MembershipFactory.MustCreateWithOption(map[string]interface{}{
			"GameID":      owner.GameID,
			"PlayerID":    players[i].ID,
			"ClanID":      clan.ID,
			"RequestorID": owner.ID,
			"Metadata":    "{\"x\": 1}",
		}).(*Membership)

		err = db.Insert(membership)
		if err != nil {
			return nil, nil, nil, nil, err
		}

		memberships = append(memberships, membership)
	}

	return clan, owner, players, memberships, nil
}
