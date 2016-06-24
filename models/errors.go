// khan
// https://github.com/topfreegames/khan
//
// Licensed under the MIT license:
// http://www.opensource.org/licenses/mit-license
// Copyright © 2016 Top Free Games <backend@tfgco.com>

package models

import "fmt"

// ModelNotFoundError identifies that a given model was not found in the Database with the given ID
type ModelNotFoundError struct {
	Type string
	ID   interface{}
}

func (e *ModelNotFoundError) Error() string {
	return fmt.Sprintf("%s was not found with id: %v", e.Type, e.ID)
}

// EmptyGameIDError identifies that a request was made for a model without the proper game id
type EmptyGameIDError struct {
	Type string
}

func (e *EmptyGameIDError) Error() string {
	return fmt.Sprintf("Game ID is required to retrieve %s!", e.Type)
}

// ClanReachedMaxMembersError identifies that a given clan already reached the max number of members allowed
type ClanReachedMaxMembersError struct {
	ID interface{}
}

func (e *ClanReachedMaxMembersError) Error() string {
	return fmt.Sprintf("Clan %s reached max members", e.ID)
}

// PlayerReachedMaxClansError identifies that a given player already reached the max number of clans allowed
type PlayerReachedMaxClansError struct {
	ID interface{}
}

func (e *PlayerReachedMaxClansError) Error() string {
	return fmt.Sprintf("Player %s reached max clans", e.ID)
}

// PlayerCannotCreateMembershipError identifies that a given player is not allowed to create a membership
type PlayerCannotCreateMembershipError struct {
	PlayerID interface{}
	ClanID   interface{}
}

func (e *PlayerCannotCreateMembershipError) Error() string {
	return fmt.Sprintf("Player %v cannot create membership for clan %v", e.PlayerID, e.ClanID)
}

// PlayerCannotPerformMembershipActionError identifies that a given player is not allowed to promote/demote another member
type PlayerCannotPerformMembershipActionError struct {
	Action      string
	PlayerID    interface{}
	ClanID      interface{}
	RequestorID interface{}
}

func (e *PlayerCannotPerformMembershipActionError) Error() string {
	return fmt.Sprintf("Player %v cannot %s membership for player %s and clan %v", e.RequestorID, e.Action, e.PlayerID, e.ClanID)
}

// CannotApproveOrDenyMembershipAlreadyProcessedError identifies that a membership that is already processed cannot be approved or denied
type CannotApproveOrDenyMembershipAlreadyProcessedError struct {
	Action string
}

func (e *CannotApproveOrDenyMembershipAlreadyProcessedError) Error() string {
	return fmt.Sprintf("Cannot %s membership that was already approved or denied", e.Action)
}

// CannotPromoteOrDemoteInvalidMemberError identifies that a given player is not allowed to promote/demote a member
type CannotPromoteOrDemoteInvalidMemberError struct {
	Action string
}

func (e *CannotPromoteOrDemoteInvalidMemberError) Error() string {
	return fmt.Sprintf("Cannot %s membership that is denied or not yet approved", e.Action)
}

// CannotPromoteOrDemoteMemberLevelError identifies that a given member is already max level and cannot be promoted
type CannotPromoteOrDemoteMemberLevelError struct {
	Action string
	Level  int
}

func (e *CannotPromoteOrDemoteMemberLevelError) Error() string {
	return fmt.Sprintf("Cannot %s member that is already level %d", e.Action, e.Level)
}

// InvalidMembershipActionError identifies that a given action is not valid
type InvalidMembershipActionError struct {
	Action string
}

func (e *InvalidMembershipActionError) Error() string {
	return fmt.Sprintf("%s a membership is not a valid action.", e.Action)
}

// InvalidLevelForGameError identifies that a given level is not valid for the given game
type InvalidLevelForGameError struct {
	GameID string
	Level  interface{}
}

func (e *InvalidLevelForGameError) Error() string {
	return fmt.Sprintf("Level %v is not valid for game %s.", e.Level, e.GameID)
}

// ClanHasNoMembersError identifies that a clan has no members
type ClanHasNoMembersError struct {
	ClanID interface{}
}

func (e *ClanHasNoMembersError) Error() string {
	return fmt.Sprintf("Clan %v has no members", e.ClanID)
}

// EmptySearchTermError identifies that a search term was not provided
type EmptySearchTermError struct{}

func (e *EmptySearchTermError) Error() string {
	return "A search term was not provided to find a clan."
}
