package domain

import (
	"testing"

	"github.com/coinflipgamesllc/api.playtest-coop.com/domain/game"
)

func TestNewGame(t *testing.T) {
	g := NewGame("First Game", User{ID: 123})
	if g.Title != "First Game" {
		t.Error("Title incorrect on new game")
	}

	if g.Status != game.Prototype {
		t.Error("Default status is not Prototype on new game")
	}

	if len(g.Designers) != 1 && g.Designers[0].ID != 123 {
		t.Error("Primary designer not set on new game")
	}

	if g.Stats.MinPlayers != 1 || g.Stats.MaxPlayers != 5 || g.Stats.MinAge != 8 || g.Stats.EstimatedPlaytime != 30 {
		t.Error("Game stats are not set to expected defaults on new game")
	}
}

func TestMayBeUpdatedBy(t *testing.T) {
	var tests = []struct {
		game          *Game
		user          *User
		expectAllowed bool
	}{
		{&Game{Designers: []User{User{ID: 2}, User{ID: 3}}}, &User{ID: 2}, true},
		{&Game{Designers: []User{User{ID: 2}}}, &User{ID: 5}, false},
		{&Game{Designers: []User{User{ID: 2}}}, nil, false},
	}

	for _, tt := range tests {
		actual := tt.game.MayBeUpdatedBy(tt.user)
		if tt.expectAllowed != actual {
			t.Errorf("Editing permission incorrect")
		}
	}
}

func TestRename(t *testing.T) {
	var tests = []struct {
		game         *Game
		newName      string
		expectedName string
	}{
		{&Game{Title: "Original Name"}, "New Name", "New Name"},
		{&Game{Title: "Original Name"}, "", "Original Name"},
		{&Game{Title: "Original Name"}, "Original Name", "Original Name"},
	}

	for _, tt := range tests {
		tt.game.Rename(tt.newName)
		actual := tt.game.Title
		if tt.expectedName != actual {
			t.Errorf("Rename incorrect")
		}
	}
}

func TestUpdateOverview(t *testing.T) {
	var tests = []struct {
		game             *Game
		newOverview      string
		expectedOverview string
	}{
		{&Game{Overview: "Original Overview"}, "New Overview", "New Overview"},
		{&Game{Overview: "Original Overview"}, "", "Original Overview"},
		{&Game{Overview: "Original Overview"}, "Original Overview", "Original Overview"},
	}

	for _, tt := range tests {
		tt.game.UpdateOverview(tt.newOverview)
		actual := tt.game.Overview
		if tt.expectedOverview != actual {
			t.Errorf("UpdateOverview incorrect")
		}
	}
}

func TestUpdateStatus(t *testing.T) {
	var tests = []struct {
		game           *Game
		newStatus      string
		expectedStatus game.Status
		expectedError  error
	}{
		{&Game{Status: game.Prototype}, "Published", game.Published, nil},
		{&Game{Status: game.Prototype}, "Not a status", game.Prototype, game.InvalidStatus{}},
	}

	for _, tt := range tests {
		err := tt.game.UpdateStatus(tt.newStatus)
		if tt.expectedError != nil {
			if _, ok := err.(game.InvalidStatus); !ok {
				t.Errorf("Expected error on invalid status, got none")
			}
		}

		actual := tt.game.Status
		if actual != tt.expectedStatus {
			t.Errorf("String '%s' did not produce expected status. Got '%s'", tt.newStatus, actual)
		}
	}
}

func TestAddDesigner(t *testing.T) {
	var tests = []struct {
		game              *Game
		user              *User
		expectedDesigners []uint
	}{
		{&Game{Designers: nil}, &User{ID: 1}, []uint{1}},
		{&Game{Designers: []User{}}, &User{ID: 1}, []uint{1}},
		{&Game{Designers: []User{User{ID: 1}}}, &User{ID: 2}, []uint{1, 2}},
		{&Game{Designers: []User{User{ID: 1}}}, &User{ID: 1}, []uint{1}},
	}

	for _, tt := range tests {
		tt.game.AddDesigner(tt.user)

		actualIDs := []uint{}
		for _, u := range tt.game.Designers {
			actualIDs = append(actualIDs, u.ID)
		}

		if !EqualUintSlice(actualIDs, tt.expectedDesigners) {
			t.Errorf("Mismatched designers after add")
		}
	}
}

func TestReplaceDesigners(t *testing.T) {
	var tests = []struct {
		game              *Game
		designers         []User
		expectedDesigners []uint
	}{
		{&Game{Designers: nil}, []User{User{ID: 1}}, []uint{1}},
		{&Game{Designers: []User{}}, []User{User{ID: 1}}, []uint{1}},
		{&Game{Designers: []User{User{ID: 1}}}, []User{User{ID: 2}}, []uint{2}},
		{&Game{Designers: []User{User{ID: 1}}}, []User{User{ID: 1}}, []uint{1}},
	}

	for _, tt := range tests {
		tt.game.ReplaceDesigners(tt.designers)

		actualIDs := []uint{}
		for _, u := range tt.game.Designers {
			actualIDs = append(actualIDs, u.ID)
		}

		if !EqualUintSlice(actualIDs, tt.expectedDesigners) {
			t.Errorf("Mismatched designers after replace")
		}
	}
}

func EqualUintSlice(a, b []uint) bool {
	if len(a) != len(b) {
		return false
	}
	for i, v := range a {
		if v != b[i] {
			return false
		}
	}
	return true
}
