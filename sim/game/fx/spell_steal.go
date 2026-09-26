package fx

import (
	"duel-masters/game/match"
)

// SpellStealCast casts spell as though caster controlled it: "you"/"your
// opponent" in the spell's own text mean caster and caster's opponent for the
// whole resolution, even though spell's real owner never changes. Once the
// cast attempt is over, the card is unconditionally handed to its real
// owner's graveyard.
//
// This is the "Spell Steal" mechanic Bluum Erkis, Flare Guardian introduced.
// It always casts spell — never offers it as an optional cast — and none of
// the reactions to "a player used a shield trigger" see it, because CastSpell
// here is not marked FromShield.
//
// The card always ends up in owner's graveyard no matter what happened during
// the cast attempt: whether the spell resolved normally (which already moved
// it to caster's graveyard), whether a static "can't cast" restriction (such
// as Alcadeias, Lord of Spirits) prevented the cast outright (leaving it
// sitting in caster's hand), or whether the spell's own effect tried to
// redirect its own destination (such as Boomerang Comet's "put it into your
// mana zone instead of your graveyard") — that redirect never applies here,
// because caster was never going to let it reach a graveyard on its own.
//
// spell must currently belong to caster's opponent and not be in a zone
// caster already has a card of the same identity in (it is briefly relocated
// to caster's hand for the duration of the cast).
func SpellStealCast(spell *match.Card, caster *match.Player, ctx *match.Context) {

	owner := spell.Player

	if _, err := owner.GiveCardTo(spell.ID, spell.Zone, caster, match.HAND); err != nil {
		return
	}

	ctx.Match.CastSpell(spell, false)

	// CastSpell resolves synchronously, prompts included, so by the time it
	// returns the cast attempt is fully over. Force the card into its real
	// owner's graveyard regardless of where it currently sits.
	caster.GiveCardTo(spell.ID, spell.Zone, owner, match.GRAVEYARD)
}
