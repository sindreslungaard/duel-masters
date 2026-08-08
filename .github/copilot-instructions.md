# GitHub Copilot repository instructions

Read and follow [`../AGENTS.md`](../AGENTS.md) in full before proposing or making repository changes. It is the canonical instruction set for this project.

For every new or changed card, Copilot must:

- use the top-level `cards[].id` in `sim/DuelMastersCards.json` for the canonical UUID and use the same file for the exact name, base power, mana cost, civilization, races/type, printings, and complete effect text; do not confuse the UUID with a printing's collector-number `id`;
- inspect cards with similar timing and effects and search all of `sim/game/fx` before writing custom logic;
- reason explicitly about handler order, zones, pre/post events, cancellation, nested events, `ScheduleAfter`, persistent-effect cleanup, multiple effect sources, evolution piles, targeting protection, shields, and destruction replacements;
- add scenario-based automated tests, including applicable interaction and empty/cancel/cleanup cases; extend `sim/tests/scenario` with reusable bounded helpers when the current harness is insufficient;
- treat `Player.Action` as an unbuffered, single-consumer protocol: prefer existing `fx` selection helpers, never add ad hoc goroutines or channel closes, balance every prompt/wait path, and test channel/goroutine changes with timeouts and the race detector;
- run focused tests and `go test ./...`; run `go test -race ./...` for action, scheduling, persistent-effect, disposal, socket, or goroutine changes.

Do not generate a card implementation from printed text alone. Existing engine behavior and similar implementations must be inspected first.
