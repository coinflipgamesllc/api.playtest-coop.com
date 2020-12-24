package game

// Mechanic is the name of a construct used in a game system
type Mechanic string

// AvailableMechanics lists all the mechanic options available
func AvailableMechanics() []string {
	return []string{
		"Action Blocking",
		"Action Points",
		"Action/Role Selection",
		"Area Control",
		"Area Enclosure",
		"Area Majority",
		"Asymmetric",
		"Auction/Bidding",
		"Betting",
		"Campaign",
		"Card Driven",
		"Cube Pusher",
		"Deck Construction",
		"Deck Improvement",
		"Deduction",
		"Dexterity",
		"Dice Rolling",
		"Drafting",
		"Drawing",
		"Engine Building",
		"Folk on a Map",
		"Grid/Area Movement",
		"Hand Management",
		"Hex and Counter",
		"Hidden Deployment",
		"Hidden Movement",
		"Maintenance Cost",
		"Memory",
		"Multi-Use Cards",
		"Negotiation",
		"One Vs. Many",
		"Pattern Building",
		"Player Elimination",
		"Point To Point Movement",
		"Press Your Luck",
		"Real Time",
		"Resource Management",
		"Role Playing",
		"Roll to Move",
		"Rondel",
		"Route Building",
		"Scenario-Driven",
		"Set Collection",
		"Social Deduction",
		"Stocks",
		"Take That",
		"Tile Placement",
		"Time Track",
		"Trading",
		"Trick-Taking",
		"Variable Player Order",
		"Voting",
		"Worker Placement",
		"X And Write",
	}
}
