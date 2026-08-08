# Repository instructions

These instructions apply to every automated coding agent working in this repository. `AGENTS.md` is the canonical copy; `CLAUDE.md` imports it and the Copilot instructions point to it.

## Before changing code

- Read the relevant implementation, its tests, and the engine code that invokes it. Do not infer engine behavior from a card's wording alone.
- Preserve unrelated working-tree changes. Keep card changes narrowly scoped and follow `CONTRIBUTING.md` (normally one card per change for new contributors).
- Prefer existing abstractions and repository conventions. Search with `rg` before adding a helper, condition, event, selection loop, or goroutine.
- For any behavior involving `Player.Action`, goroutines, event scheduling, match disposal, or sockets, treat deadlocks and goroutine leaks as correctness failures. Follow the channel rules below for all development, not only card work.

## Mandatory workflow for implementing or changing a card

### 1. Establish the exact card contract

Use `sim/DuelMastersCards.json` as the local source of truth for the card's:

- exact name;
- printed/base power (the numeric part of values such as `3000+`);
- mana cost;
- civilization;
- type and race(s);
- complete rules text, including optional words such as "may," ownership, target restrictions, duration, replacement wording, and ordering.

The JSON does not supply the simulator UUID. Find that in the appropriate set map in `sim/game/cards/repository.go`. Confirm the printing/set before editing. Do not silently "correct" the local JSON from memory or an external card database.

Translate every independent sentence or clause into an explicit implementation obligation before coding. Record at least: trigger/event, controller making each choice, eligible cards and zones, minimum/maximum selections, whether cancellation is legal, movement/destruction source, duration, cleanup condition, and interaction with replacement or prevention effects.

### 2. Find and compare similar cards first

This is mandatory even when the effect looks simple.

- Search `sim/DuelMastersCards.json` for distinctive rules phrases to find cards with the same printed mechanic.
- Search `sim/game/cards` for the relevant event, condition, helper, and wording. Compare at least one existing implementation with the same timing and one with the same operation when available.
- Search all of `sim/game/fx` before writing custom logic. Common behavior already exists for creatures, spells, evolution, blockers, shield triggers, breakers, power attacker, slayer, targeting, draw/search, movement, destruction, tap/untap, attack restrictions, persistent effects, and player choices.
- Inspect the helper implementation, not only its name. Confirm its timing, optionality, target owner, selection behavior when there are too few/zero targets, destruction context, and cleanup semantics match the card exactly.
- Prefer a newer, tested implementation when examples disagree, but do not blindly copy it. Check the engine path and printed text yourself. If examples expose inconsistent behavior, fix or document the shared issue rather than introducing another variant accidentally.

Useful searches include:

```text
rg -n 'distinctive card text' sim/DuelMastersCards.json
rg -n 'EventName|ConditionName|HelperName' sim/game/cards sim/game/fx
rg -n 'CardName|constructorName' sim/game/cards/repository.go sim/game/cards sim/tests
rg -n '^func ' sim/game/fx
```

### 3. Put the code in the established locations

- Card constructors live under `sim/game/cards/<set>/`, grouped by race; spells live in that set's `spells.go`.
- Set all printed fields explicitly using `civ` and `family` constants: `Name`, `Power` for creatures, `Civ`, `Family`, `ManaCost`, and `ManaRequirement`.
- A normal creature starts with `c.Use(fx.Creature, ...)`; a normal spell starts with `c.Use(fx.Spell, ...)`. Add `fx.Evolution`, `fx.ShieldTrigger`, `fx.Charger`, breaker helpers, and other standard mechanics as required by the exact text.
- `c.Use` order is semantic. Default handlers commonly validate, interrupt, mutate events, or schedule default behavior. Keep the established order from genuinely similar cards and reason about cancellation and scheduled callbacks before changing it.
- Register the constructor under the exact existing UUID in the correct set map in `sim/game/cards/repository.go`. The executable and scenario harness register those maps; do not add ad hoc `match.AddCard` calls in card files.
- If several cards need the same operation, add or improve a focused helper in `sim/game/fx` and test the helper's boundaries. Do not create a card-specific helper merely to hide complex logic.

### 4. Model the engine's event semantics correctly

`Match.HandleFx` is synchronous and order-sensitive:

1. it snapshots cards from both players' zones, with the active player first;
2. it runs match persistent effects;
3. it runs each card's `c.Use` handlers in registration order;
4. it runs `Context.ScheduleAfter` callbacks;
5. cancellation or match ending stops the remaining flow.

Important consequences:

- A card's handlers are considered while it is in the battle zone, spell zone, hand, shields, hidden zone, mana, graveyard, and deck. Every non-zone-agnostic effect must guard the required zone, normally `card.Zone == match.BATTLEZONE`, or use a vetted hook that does so.
- Effects from multiple cards are processed active-player first and in container/card order, while persistent effects are stored in a map and must not rely on map iteration order. Do not implement simultaneous interactions whose correctness depends on an incidental handler order.
- `ScheduleAfter` means "after the normal handler traversal for this context," not "later in another goroutine." Use it where the effect must see all event modifiers or where default behavior should occur only after replacement/prevention handlers. A cancelled context does not run remaining scheduled callbacks.
- `InterruptFlow` cancels the entire current context. Use it only for an actual prevention/replacement/cancel outcome. Mutating an event slice (for example, removing one shield from a multi-break) is different from cancelling the whole event.
- Re-entrant calls such as `MoveCard`, `Destroy`, `Battle`, `BreakShields`, and `GetPower` dispatch more contexts synchronously. Re-check zones after effects that can move participants; never assume a selected attacker, blocker, source, or target is still in play.

Use the correct timing event:

- `MoveCard` is the pre-move event and can be interrupted; `CardMoved` is post-move.
- `CardPlayedEvent` occurs after mana selection but before the card enters its destination.
- `SpellCast` represents both normal casts and shield-trigger casts and includes `FromShield`/player identity.
- `CreatureDestroyed` is the pre-destruction/replacement point; the later battle-zone-to-graveyard `CardMoved` is the reliable post-destruction observation.
- `AttackPlayer`/`AttackCreature` are attempts. `AttackConfirmed` is after validation and tapping. Player attacks then pass through shield selection, blocker selection, `Block`, possibly `Battle`, and finally shield breaking. Choose the hook matching the printed "attacks," "isn't blocked," "becomes blocked," or "battles" wording.
- `BreakShieldEvent` is pre-break and mutable/preventable; `BrokenShieldEvent` is post-break. Shield-trigger processing is a separate sequence and may move or replay cards.
- `EndStep` and `EndOfTurnStep` are distinct. Existing temporary effects usually expire on `EndOfTurnStep`; inspect similar effects when end-of-turn triggers can still destroy or move cards.
- `Summoned` excludes an evolution base reappearing from `HIDDENZONE`; `InTheBattlezone` includes that transition and is appropriate for continuous "while in the battle zone" setup.

### 5. Preserve state, stacking, and cleanup invariants

- Conditions are a multi-source collection, not a boolean set. Use `AddUniqueSourceCondition` for persistent/granted effects and a stable source, normally `card.ID`.
- Remove only the contribution you own with `RemoveSpecificConditionBySource` or `RemoveConditionBySource`. Do not use `RemoveCondition` for a granted effect when another card may grant the same condition.
- `fx.Creature` and `fx.Spell` clear conditions at `EndOfTurnStep`; many intrinsic conditions are rebuilt on `UntapStep`. Test effects across a turn boundary rather than assuming a constructor-time condition persists.
- For a continuous effect, install it once on `fx.InTheBattlezone`, reapply source-tagged state as needed, clean it from every affected card when the source leaves the battle zone, and call the persistent effect's `exit`. Also handle newly arriving eligible creatures and multiple copies of the source.
- For "this turn" persistent effects, exit at the correct turn event and remove any conditions they added. If end-of-turn effects must still observe the modifier, use the established `ScheduleAfter(exit)` pattern.
- Never add the same persistent effect on every unrelated event. Gate installation to the source entering play, a spell resolving, a tap ability resolving, or the exact printed trigger.
- Use `ctx.Match.GetPower(card, attacking)` for dynamic power comparisons. `Card.Power` is only printed base power. During `Battle`, power has already been captured in `AttackerPower` and `DefenderPower`; battle-specific modifiers must update those event fields, while longer-lived bonuses need source-tagged conditions or a persistent `GetPowerEvent` modifier.
- Do not cache a derived list across events unless the card text freezes that set. Re-evaluate continuous effects as the battle zone changes.
- Use card identity (`ID` or pointer) and owner, not name, for sources and targets. Duplicate copies are normal.

### 6. Respect zones, movement, evolution, and targeting

- Move cards through `Player.MoveCard`, `Match.MoveCard`, `Destroy`, `BreakShields`, or an existing `fx` helper as appropriate so pre/post events and replacement effects run. Directly editing zone slices bypasses game rules.
- `Player.MoveCard` returns `(card, nil)` when a pre-move event cancelled the move. Verify the resulting `card.Zone` when subsequent logic depends on success; `err == nil` alone is not proof that movement occurred.
- A successful move resets `Tapped` to false. Only override tap state when the rules explicitly require it and after confirming the move happened.
- Pass the real effect source ID/card and the correct `CreatureDestroyedContext`. Other cards distinguish battle, spell, slayer, and miscellaneous destruction.
- Evolution bases move to `HIDDENZONE` and are attached to the top card. When the evolution leaves the battle zone, its pile follows. Test effects that move, destroy, bounce, or replace destruction of evolution creatures; do not treat hidden cards as ordinary active creatures.
- Use `fx.SelectFilter`, `SelectMultipart`, or their backside variants for targeting. These helpers enforce `CantBeSelectedByOpp`; hand-written targeting can accidentally bypass protection.
- Keep "choose" separate from non-targeting/global effects. Do not run selectability filtering for effects that do not choose, and do not bypass it for effects that do.
- Distinguish the player making a choice from the owner of the container. Many discard/sacrifice effects are chosen by the opponent, not by the effect's controller.
- Preserve hidden information: use backside selection for shields or other face-down cards, and reveal/show only what the printed effect permits.
- Snapshot a collection before moving several cards out of it. Avoid ranging a live zone slice while mutating that same zone.

### 7. Make optionality and empty-target behavior exact

- "May"/"up to" normally requires a cancellable action. Mandatory text must not be cancellable merely for convenience.
- Selection bounds must adapt safely when fewer legal cards exist. Inspect the selection helper's behavior and test zero targets, fewer than the requested maximum, and exactly one forced target.
- Do not open an action prompt when no legal choice exists. Existing `fx` selection helpers return an empty collection in that case.
- An optional cost and its payoff are one transaction semantically: do not apply the payoff if the player cancels, cannot fully pay, or a selected cost card moved before payment.
- Validate selections against the exact offered snapshot and re-fetch from the expected zone before applying them. Never trust client `Count`, card IDs, or cancellation flags without validation.

## Player.Action channels and goroutines

`Player.Action` is an unbuffered channel. The parser delivers action messages from a parallel event, while card resolution normally blocks the sequential event loop waiting for exactly one response. A small mistake can leave the match event loop, a sender, or a test goroutine blocked forever.

Rules for all code:

- Prefer `fx.Select`, `SelectFilter`, `SelectBackside`, `SelectMultipart`, `BinaryQuestion`, `SelectCount`, `OrderCards`, and multiple-choice helpers. They centralize validation and balance `NewAction`/`CloseAction` plus opponent `Wait`/`EndWait` messages.
- Never create a goroutine merely to wait on or send to `Player.Action`. Never close the channel except through the existing player/match disposal owner. Never replace it with a buffered channel to mask a protocol bug.
- Allow only one outstanding prompt/receiver per player. Do not install persistent handlers that can prompt concurrently, prompt from power calculation/filter predicates, or prompt from a `ParallelEvent` path.
- Before waiting, prove there is a legal response path. Handle empty lists synchronously. Every cancellation, invalid selection, moved target, early return, match end, and disconnect path must leave action and wait UI state balanced.
- Use `defer CloseAction`/`defer EndWait` where the existing helper pattern supports it. In manual loops, validate card count, membership, current zone, and optional cancellation before mutating state; warnings must continue the same loop without spawning another waiter.
- Do not hold a mutex while sending, receiving, prompting, invoking `HandleFx`, or doing network I/O. These calls can block or re-enter engine code.
- Remember that disposing a player closes `Action`; receives from a closed channel immediately return zero values. A loop that treats that value as merely invalid can spin forever. New low-level prompt code must handle channel closure explicitly or, preferably, improve and reuse a shared helper.
- The event loop already schedules work. Avoid nested goroutines in card/effect code. If general development truly needs background work, give it explicit ownership, cancellation, bounded lifetime, and a test proving disposal terminates it.

## Required testing

Every implemented behavior should have an automated test when technically possible. "The scenario helper does not support it" is a reason to extend `sim/tests/scenario`, not a reason to leave a card untested.

### Test at the right level

- Put card integration tests in `sim/tests/cards/<snake_case_card_name>_test.go` and use the UUID from `repository.go` as a named constant.
- Prefer `scenario.New`, real actions, and public match behavior over calling a handler directly. This covers handler order, mana payment, zones, event scheduling, prompts, and state broadcasts.
- Seed deterministic state with `SpawnCard`; for deck-order effects, destroy/rebuild the deck so the exact top cards are known. Do not assert random outcomes.
- Capture `scenario.MessageCount` before the operation that opens a prompt, then wait/submit/cancel through scenario helpers. Do not use arbitrary sleeps. Every test-triggered prompt must receive a response so the sequential event loop is not leaked.
- If a useful action is missing, add a generic scenario method (for example, waiting for and inspecting an action, attacking, blocking, using a tap ability, or selecting a target) with a bounded timeout and membership validation. Do not add a one-card-only testing shortcut.

### Minimum behavior matrix

At minimum, cover the card's printed stats/standard mechanics and its happy path. Add all applicable interaction cases:

- no legal targets, fewer than the maximum, optional cancellation, and mandatory choice;
- source leaves play before a delayed/continuous effect resolves;
- target leaves or changes zone during nested effects;
- two copies of the source, stacking bonuses, and removing only one source;
- turn transition and exact expiration time;
- controller versus opponent choices and both players' turns;
- blockers, cannot-be-selected/attacked/blocked effects, dynamic power modifiers, and summoning sickness;
- destruction replacement, slayer, battle destruction, and the correct destruction context;
- shield breaking prevention/modification, multiple breakers, shield triggers, and cards that return themselves;
- evolution piles and `HIDDENZONE` transitions;
- invalid or stale client selections for any new low-level action code;
- match disposal/cancellation for any channel or goroutine change.

Pairwise tests with a representative existing special-effect card are more valuable than isolated condition assertions. Assert observable outcomes (zones, tap state, effective power, choices, cleanup, and continued turn progress), not only that a condition was added.

### Validation commands

Run from `sim`:

```text
go test ./tests/cards -run 'TestCardName' -count=1 -timeout 30s
go test ./... -count=1 -timeout 120s
```

For changes to actions, event scheduling, persistent effects, disposal, sockets, or goroutines, also run:

```text
go test -race ./... -count=1 -timeout 180s
```

Run `gofmt` on changed Go files. A timeout, leaked waiter, data race, or flaky test is a failed implementation, even if the card's isolated happy path works.

## Final review checklist for a card

Before declaring the work complete, verify all of the following:

- Constructor metadata matches the exact `DuelMastersCards.json` entry and the UUID/set registration is correct.
- Each rules-text clause maps to the correct event and zone, including "may," "up to," "other," owner/controller, duration, and replacement timing.
- Similar cards and every relevant `fx` helper were inspected; duplicated logic is justified.
- `c.Use` ordering, cancellation, nested events, and `ScheduleAfter` behavior were reasoned through.
- Continuous/temporary conditions use unique source IDs and clean up without removing another card's contribution.
- Movement, destruction, power, targeting protection, hidden information, shields, and evolution interactions use engine APIs.
- Empty, insufficient, stale, and cancelled selections cannot block the event loop.
- No new unowned goroutine or unmatched `Player.Action` send/receive/close exists.
- Focused, interaction, full-suite, and (when relevant) race tests pass.
