# gonixgo-example

A small Go program built with [gonixgo](https://github.com/draganm/gonixgo):
one Nix derivation per Go package, and no Nix file to update when `go.mod`
changes.

The program, `greet`, has a main package, two local packages (one embeds a
file) and a third-party dependency, `github.com/fatih/color`.

## Build and run

gonixgo runs a program during Nix evaluation through `builtins.exec`, which
Nix only provides when asked:

```bash
nix build --option allow-unsafe-native-code-during-evaluation true
./result/bin/greet -name you

nix run --option allow-unsafe-native-code-during-evaluation true . -- -name you
```

With that option on, any Nix expression you evaluate can run programs as
you. Pass it per command, for flakes you trust, rather than enabling it in
`nix.conf`.

The first build compiles gonixgo's tool and the Go standard library; later
builds reuse them.

## What the flake does

```nix
goEnv = gonixgo.lib.mkGoEnv { inherit pkgs; };

packages.default = goEnv.buildGoApplication {
  pname = "greet";
  version = "0.1.0";
  src = ./.;
  subPackages = [ "cmd/greet" ];
  ldflags = [ "-X main.version=0.1.0" ];
};
```

There is no `vendorHash` and no lockfile of module hashes. While Nix
evaluates, gonixgo runs `go list` on this source, hashes the modules it finds
in your Go module cache, and adds them to the Nix store. The `pkgs` you pass
builds gonixgo's own tool and performs the Go build.

## One derivation per package

The result exposes every package, module and binary:

```bash
# The packages in the build.
nix eval --option allow-unsafe-native-code-during-evaluation true \
  --json .#default.packages --apply builtins.attrNames

# Build one package alone.
nix build --option allow-unsafe-native-code-during-evaluation true \
  --no-link --print-out-paths \
  '.#default.packages."github.com/draganm/gonixgo-example/internal/greeting"'
```

Edit `internal/greeting/greeting.go` and build again: Nix recompiles that
package and `cmd/greet`, which imports it, and links. The other local
package and the four third-party packages are not rebuilt.

## Changing dependencies

```bash
go get github.com/fatih/color@latest
go mod tidy
nix build --option allow-unsafe-native-code-during-evaluation true
```

Nothing else needs regenerating.

## Development shell

`nix develop --option allow-unsafe-native-code-during-evaluation true` (or
`direnv allow`) gives a shell with the same Go the build uses and with
`greet` itself, built by gonixgo, on `PATH`. `go build`, `go test` and
`go run ./cmd/greet` work as usual. gonixgo does not run tests yet.

With direnv, `.envrc` watches `go.mod`, `go.sum`, `cmd/` and `internal/`:
after you change a source file, the next prompt in the directory reloads the
shell, and the `greet` on `PATH` is rebuilt from the new source, again only
the packages that changed. Nix sees only the files git tracks, so `git add`
a new file before it shows up in the build.

## Notes

- `nix develop`, `nix flake check` and `nix flake show` evaluate `packages`,
  so they need the option too.
- gonixgo 0.1.0 builds pure-Go programs. Packages that use cgo and modules
  with `replace` directives are rejected with an explanation; see gonixgo's
  README for the current limits.
