// Code generated by SQLBoiler 4.16.2 (https://github.com/volatiletech/sqlboiler). DO NOT EDIT.
// This file is meant to be re-generated in place and/or deleted at any time.

package models

import "testing"

// TestToOne tests cannot be run in parallel
// or deadlocks can occur.
func TestToOne(t *testing.T) {
	t.Run("AuthToUserUsingImpersonateAs", testAuthToOneUserUsingImpersonateAs)
	t.Run("AuthToUserUsingUser", testAuthToOneUserUsingUser)
	t.Run("ConfiantToPublisherUsingPublisher", testConfiantToOnePublisherUsingPublisher)
	t.Run("DemandParnterPlacementToDemandPartnerUsingDemandPartner", testDemandParnterPlacementToOneDemandPartnerUsingDemandPartner)
	t.Run("DpoRuleToDpoUsingDemandPartner", testDpoRuleToOneDpoUsingDemandPartner)
	t.Run("DpoRuleToPublisherUsingDpoRulePublisher", testDpoRuleToOnePublisherUsingDpoRulePublisher)
	t.Run("FloorToPublisherUsingFloorPublisher", testFloorToOnePublisherUsingFloorPublisher)
	t.Run("PixalateToPublisherUsingPublisher", testPixalateToOnePublisherUsingPublisher)
	t.Run("PublisherDomainToPublisherUsingPublisher", testPublisherDomainToOnePublisherUsingPublisher)
	t.Run("TargetingToPublisherUsingTargetingPublisher", testTargetingToOnePublisherUsingTargetingPublisher)
	t.Run("UserPlatformRoleToUserUsingUser", testUserPlatformRoleToOneUserUsingUser)
}

// TestOneToOne tests cannot be run in parallel
// or deadlocks can occur.
func TestOneToOne(t *testing.T) {}

// TestToMany tests cannot be run in parallel
// or deadlocks can occur.
func TestToMany(t *testing.T) {
	t.Run("DemandPartnerToDemandParnterPlacements", testDemandPartnerToManyDemandParnterPlacements)
	t.Run("DpoToDemandPartnerDpoRules", testDpoToManyDemandPartnerDpoRules)
	t.Run("PublisherToConfiants", testPublisherToManyConfiants)
	t.Run("PublisherToDpoRules", testPublisherToManyDpoRules)
	t.Run("PublisherToFloors", testPublisherToManyFloors)
	t.Run("PublisherToPixalates", testPublisherToManyPixalates)
	t.Run("PublisherToPublisherDomains", testPublisherToManyPublisherDomains)
	t.Run("PublisherToTargetings", testPublisherToManyTargetings)
	t.Run("UserToImpersonateAsAuths", testUserToManyImpersonateAsAuths)
	t.Run("UserToAuths", testUserToManyAuths)
	t.Run("UserToUserPlatformRoles", testUserToManyUserPlatformRoles)
}

// TestToOneSet tests cannot be run in parallel
// or deadlocks can occur.
func TestToOneSet(t *testing.T) {
	t.Run("AuthToUserUsingImpersonateAsAuths", testAuthToOneSetOpUserUsingImpersonateAs)
	t.Run("AuthToUserUsingAuths", testAuthToOneSetOpUserUsingUser)
	t.Run("ConfiantToPublisherUsingConfiants", testConfiantToOneSetOpPublisherUsingPublisher)
	t.Run("DemandParnterPlacementToDemandPartnerUsingDemandParnterPlacements", testDemandParnterPlacementToOneSetOpDemandPartnerUsingDemandPartner)
	t.Run("DpoRuleToDpoUsingDemandPartnerDpoRules", testDpoRuleToOneSetOpDpoUsingDemandPartner)
	t.Run("DpoRuleToPublisherUsingDpoRules", testDpoRuleToOneSetOpPublisherUsingDpoRulePublisher)
	t.Run("FloorToPublisherUsingFloors", testFloorToOneSetOpPublisherUsingFloorPublisher)
	t.Run("PixalateToPublisherUsingPixalates", testPixalateToOneSetOpPublisherUsingPublisher)
	t.Run("PublisherDomainToPublisherUsingPublisherDomains", testPublisherDomainToOneSetOpPublisherUsingPublisher)
	t.Run("TargetingToPublisherUsingTargetings", testTargetingToOneSetOpPublisherUsingTargetingPublisher)
	t.Run("UserPlatformRoleToUserUsingUserPlatformRoles", testUserPlatformRoleToOneSetOpUserUsingUser)
}

// TestToOneRemove tests cannot be run in parallel
// or deadlocks can occur.
func TestToOneRemove(t *testing.T) {
	t.Run("AuthToUserUsingImpersonateAsAuths", testAuthToOneRemoveOpUserUsingImpersonateAs)
	t.Run("DpoRuleToPublisherUsingDpoRules", testDpoRuleToOneRemoveOpPublisherUsingDpoRulePublisher)
	t.Run("TargetingToPublisherUsingTargetings", testTargetingToOneRemoveOpPublisherUsingTargetingPublisher)
}

// TestOneToOneSet tests cannot be run in parallel
// or deadlocks can occur.
func TestOneToOneSet(t *testing.T) {}

// TestOneToOneRemove tests cannot be run in parallel
// or deadlocks can occur.
func TestOneToOneRemove(t *testing.T) {}

// TestToManyAdd tests cannot be run in parallel
// or deadlocks can occur.
func TestToManyAdd(t *testing.T) {
	t.Run("DemandPartnerToDemandParnterPlacements", testDemandPartnerToManyAddOpDemandParnterPlacements)
	t.Run("DpoToDemandPartnerDpoRules", testDpoToManyAddOpDemandPartnerDpoRules)
	t.Run("PublisherToConfiants", testPublisherToManyAddOpConfiants)
	t.Run("PublisherToDpoRules", testPublisherToManyAddOpDpoRules)
	t.Run("PublisherToFloors", testPublisherToManyAddOpFloors)
	t.Run("PublisherToPixalates", testPublisherToManyAddOpPixalates)
	t.Run("PublisherToPublisherDomains", testPublisherToManyAddOpPublisherDomains)
	t.Run("PublisherToTargetings", testPublisherToManyAddOpTargetings)
	t.Run("UserToImpersonateAsAuths", testUserToManyAddOpImpersonateAsAuths)
	t.Run("UserToAuths", testUserToManyAddOpAuths)
	t.Run("UserToUserPlatformRoles", testUserToManyAddOpUserPlatformRoles)
}

// TestToManySet tests cannot be run in parallel
// or deadlocks can occur.
func TestToManySet(t *testing.T) {
	t.Run("PublisherToDpoRules", testPublisherToManySetOpDpoRules)
	t.Run("PublisherToTargetings", testPublisherToManySetOpTargetings)
	t.Run("UserToImpersonateAsAuths", testUserToManySetOpImpersonateAsAuths)
}

// TestToManyRemove tests cannot be run in parallel
// or deadlocks can occur.
func TestToManyRemove(t *testing.T) {
	t.Run("PublisherToDpoRules", testPublisherToManyRemoveOpDpoRules)
	t.Run("PublisherToTargetings", testPublisherToManyRemoveOpTargetings)
	t.Run("UserToImpersonateAsAuths", testUserToManyRemoveOpImpersonateAsAuths)
}