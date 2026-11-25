# Boomerang

A re-engineered, terminal-based version of the *Boomerang* draft-and-write card game, written in Go as part of the Software Engineering home exam at LTU. The original single-file Java implementation (`BoomerangAustralia.java`) has been split into separate packages for game rules, game flow, players, networking and user interface. The goal of the split is to make the code easier to modify, extend and test.

The game supports the **Australia** edition required by the exam (rules 1–12), and the optional future modifications: **alternating pass direction** (requirement 13) and the **Europe** and **USA** editions (requirement 14).

## Requirements

The only requirement is the **Go toolchain, version 1.27.3 or newer**. No other software is needed. The single third-party library ([argparse](https://github.com/akamensky/argparse)) is downloaded automatically by Go the first time the code is built.

Go is not shipped with Windows or macOS, so install it first:

| OS      | How to install Go |
|---------|-------------------|
| Windows | Download and run the `.msi` installer from <https://go.dev/dl/>. |
| macOS   | Download and run the `.pkg` installer from <https://go.dev/dl/>, or run `brew install go` if you use Homebrew. |
| Linux   | Use your package manager, or the archive from <https://go.dev/dl/>. |

Open a **new** terminal after installing (Command Prompt or PowerShell on Windows, Terminal on macOS), then check the installation:

```sh
go version
```

## Building

Run these commands from the project root (the folder that contains `go.mod`).

**macOS / Linux**

```sh
go build -o bin/boomerang ./cmd/boomerang
```

**Windows**

```powershell
go build -o bin\boomerang.exe .\cmd\boomerang
```

You can also skip the build step and run the game straight from source with `go run ./cmd/boomerang [options]`.

## Running the game

The program has two modes:

- **Host:** starts a game. The person running the host always plays as player 1 in the same terminal. The host waits until every remote player has connected, then fills any remaining seats with bots and starts the game.
- **Client:** joins a game on a host over TCP.

The host must be started before any client connects.

### Options

| Option | Default | Description |
|--------|---------|-------------|
| `-p`, `--players` | `2` | Number of human players, including the host (1–4). Any humans other than the host connect as remote clients. |
| `-b`, `--bots` | `0` | Number of bots that fill the remaining seats. Humans and bots together must be 2–4. |
| `-e`, `--edition` | `AU` | Edition to play: `AU` (Australia), `EU` (Europe) or `US` (USA). |
| `-d`, `--direction` | `left` | Pass direction: `left` (next player, the standard rule), `right` (previous player) or `alternate` (switch every round, requirement 13). |
| `-r`, `--rounds` | `4` | Number of rounds (4–10). The standard rules use 4. |
| `-l`, `--listen` | `7777` | TCP port the host listens on for remote players. |
| `-c`, `--connect` | | Join a host at `address` or `address:port`. If no port is given, 7777 is used. When this is set, all other options are ignored. |
| `-h`, `--help` | | Show the help text. |

### Examples

On Windows, replace `./bin/boomerang` with `.\bin\boomerang.exe`.

Play alone against one bot (comparable to `java BoomerangAustralia 1 1` in the original):

```sh
./bin/boomerang -p 1 -b 1
```

Host a game with yourself, one remote player and one bot. The host waits for the remote player on port 7777:

```sh
./bin/boomerang -p 2 -b 1
```

Join that game from another terminal or computer (comparable to `java BoomerangAustralia 127.0.0.1`):

```sh
./bin/boomerang -c 127.0.0.1
```

If the host uses a different port with `-l`, add it to the address, for example `-c 127.0.0.1:8000`.

Play the Europe edition with four players, alternating the pass direction each round:

```sh
./bin/boomerang -p 1 -b 3 -e EU -d alternate
```

If a remote player loses the connection during the game, a bot takes over that seat so the game can finish.

## Running the unit tests

The tests use Go's built-in `testing` package, so no extra test framework is needed. Run them from the project root (the commands are the same on Windows, macOS and Linux):

```sh
go test ./...
```

Each package is reported as `ok` (all tests passed) or `FAIL`, together with the test that failed. Useful variants:

```sh
go test -v ./...                         # list every test with PASS/FAIL
go test -v ./internal/engine             # run the tests in one package
go test -v -run TestCatchCard ./...      # run tests whose names match a pattern
go test -count=1 ./...                   # ignore cached results and run everything again
go test -cover ./...                     # show test coverage per package
```

### Requirement coverage

Every test is commented with the requirement number it verifies. The table shows where to find the tests for each rule.

| Req. | What is verified | Tests |
|------|------------------|-------|
| 1 | 2–4 players in total, and games with 2, 3 and 4 players run to completion | `cmd/boomerang/config_test.go` `TestPlayerCountLimits`, `internal/engine/game_test.go` `TestPlayWithTwoToFourPlayers` |
| 2 | The deck has 28 cards that match the exam's card table, each with two symbols | `internal/edition/australia/cards_test.go` |
| 3 | The deck is shuffled (a permutation that depends on the RNG) | `internal/core/deck_test.go`, `TestDealSevenUniqueCardsEachRound` |
| 4 | Each player is dealt 7 unique cards | `internal/core/deck_test.go` `TestDeckDealGivesDisjointHands`, `TestDealSevenUniqueCardsEachRound` |
| 5 | The throw card is chosen first and stays hidden until the round ends | `TestThrowCardIsHiddenUntilRoundEnd`, `TestMoveCardToDraft` |
| 6 | Hands are passed to the next player, and the last player passes to the first | `TestHandsPassToNextPlayer`, `TestPassLeftTargetsNextPlayer` |
| 7–8 | One card is picked and shown per pass until one card remains, and picks are revealed only after every player has chosen | `TestDraftedCardsAreShown`, `TestPicksRevealedTogether` |
| 9 | The catch card goes to the previous player | `TestCatchCardGoesToPreviousPlayer`, `TestCatchTargetIsOppositeOfPassing` |
| 10a–e | Throw & catch, sites, region bonus (first to complete and simultaneous), collections, animals and activities (table, one per round, once per game) | `internal/edition/australia/scoringRules_test.go`, `TestUpdateEndOfRound`, `TestScoredActivitiesCarryOver`, `internal/agent/botController_test.go` |
| 11 | Every round is scored before the next round starts | `TestFourRoundsScored`, `TestScoreRoundAppliesAllRules` |
| 12 | Four rounds by default, highest total wins, ties go to the most throw & catch points | `TestDefaultsToFourRounds`, `TestPlayReportsWinner`, `australia` `TestGetWinner`, `TestScoreFinalUsesGetWinner` |
| 13 | Alternating pass direction | `TestHandsPassAlternating`, `TestPassAlternateSwitchesEachRound` |
| 14 | Europe and USA editions: decks and scoring | `internal/edition/europe/europe_test.go`, `internal/edition/usa/usa_test.go`, `internal/edition/init_test.go` |

The networking layer has its own tests (`internal/net`, `internal/agent/remoteController_test.go`), which cover message encoding, event ordering and the bot taking over when a connection is lost.

## Project structure

```
cmd/boomerang/          Entry point: parses options and connects all the parts (composition root)
internal/core/          Domain types: Card, Symbol, Deck, PlayerState, ScoreSheet, PassingStrategy, events
internal/edition/       Edition interface and lookup (GetEdition)
  australia/            Australia deck, regions, symbols and scoring rules
  europe/               Europe edition
  usa/                  USA edition
  shared/               Scoring helpers shared by the editions
internal/scoring/       Applies an edition's rules to produce round and final results
internal/engine/        Game loop: deal, throw, draft, pass, catch, score
internal/agent/         Player controllers: human, bot and remote (network) player
internal/net/           TCP transport, listener and message protocol
internal/ui/            Text-based user interface
docs/diagrams/          Design diagrams (PlantUML sources and rendered PNGs)
```

## Design overview

The design documentation is in `docs/diagrams/`:

- `packages.png`: package dependencies. The dependencies point inwards towards `core`, and only `cmd/boomerang` knows about every concrete package.
- `classes.png`: the main types and interfaces and how they relate.
- `round.png`: sequence diagram showing the communication flow during one round.

These are the extension points that support the quality attributes:

- **New edition (extensibility, req. 14):** create a package under `internal/edition/` that implements the `Edition` interface (`Cards()`, `Rules()`, `GetWinner()`), then register it in `edition.GetEdition` and in the `--edition` option. The engine and scoring packages do not need to change.
- **New scoring rule:** implement `core.ScoringRule` (`Score(player, allStates) int`) and add it to the edition's `Rules()` map under a `Category`.
- **New pass direction (modifiability, req. 13):** add a case to `core.NewPassingStrategy`, which implements the `PassingStrategy` interface (`Target`, `CatchTarget`).
- **New kind of player or another network protocol (modifiability):** implement `engine.PlayerController`. The engine does not know whether a player is human, a bot or remote. All network code is behind the `net.Transport` interface.
- **Testability:** the engine receives its edition, its controllers, a seeded `*rand.Rand` and the event server through its constructor (dependency injection). This lets the tests run whole games with scripted players and predictable shuffles, without a terminal or a network connection.

The original Java code, the rulebooks and the exam text are kept locally in `existing/` and `exam_docs/`. Both folders are excluded from version control.

## License

See [LICENSE](LICENSE) (Creative Commons Attribution 4.0 International).
