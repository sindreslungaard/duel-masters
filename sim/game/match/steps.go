package match

// BeginTurnStep ...
// Resolve any summoning sickness from creatures in the battle zone.
type BeginTurnStep struct{}

// UntapStep ...
// Your creatures in the battle zone and cards in your mana zone are untapped. This is forced.
type UntapStep struct{}

// StartOfTurnStep ...
// Any abilities that trigger at "the start of your turn" are resolved now.
type StartOfTurnStep struct{}

// DrawStep ...
// You draw a card. This is forced. If you have no deck, you lose when you draw the last card.
type DrawStep struct{}

// ChargeStep ...
// You may put a card into your mana zone. This is optional.
type ChargeStep struct{}

// MainStep ...
// You may use cards, such as summoning creatures, casting spells, generating and crossing cross gear or fortifying castles.
// You can do these actions as many times as you want, in any order, as long as you can pay their costs.
type MainStep struct{}

// AttackStep ...
// You can attack with creatures or use Tap Abilities. You may do this in any order as long as you have eligible creatures to attack.
type AttackStep struct{}

// EndStep ...
// The turn finishes after you have no more creatures to attack with.
type EndStep struct{}

// EndOfTurnStep ...
// Any abilities that trigger at "the end of your turn" are resolved now.
type EndOfTurnStep struct{}
