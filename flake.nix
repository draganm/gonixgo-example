{
  description = "gonixgo-example";
  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-26.05";

    systems.url = "github:nix-systems/default";

    gonixgo = {
      url = "github:draganm/gonixgo/v0.2.0";
      inputs.nixpkgs.follows = "nixpkgs";
      inputs.systems.follows = "systems";
    };

  };

  outputs = { self, nixpkgs, systems, gonixgo, ... }@inputs:
    let
      eachSystem = f:
        nixpkgs.lib.genAttrs (import systems)
        (system: f system nixpkgs.legacyPackages.${system});
    in {

      # Evaluating a gonixgo package runs `gonixgo resolve` through
      # builtins.exec, so every command that touches `packages` needs
      #   --option allow-unsafe-native-code-during-evaluation true
      packages = eachSystem (system: pkgs:
        let
          # This pkgs builds gonixgo's own tool and the Go program.
          goEnv = gonixgo.lib.mkGoEnv { inherit pkgs; };
        in {
          default = goEnv.buildGoApplication {
            pname = "greet";
            version = "0.1.0";
            src = ./.;
            subPackages = [ "cmd/greet" ];
            ldflags = [ "-X main.version=0.1.0" ];
            # internal/zstd is a cgo package: its `#cgo pkg-config: libzstd`
            # line needs pkg-config and the library. Only that package's
            # compile and the link of greet get them.
            packageOverrides."github.com/draganm/gonixgo-example/internal/zstd" = {
              buildInputs = [ pkgs.zstd ];
              nativeBuildInputs = [ pkgs.pkg-config ];
            };
            meta.license = pkgs.lib.licenses.mit;
          };
        });

      # The shell carries the built program, so entering it (or direnv
      # reloading it after a source change, see .envrc) rebuilds `greet`.
      # It therefore needs the same evaluation option as `packages`.
      devShells = eachSystem (system: pkgs: {
        default = pkgs.mkShell {
          shellHook = ''
            # Set here the env vars you want to be available in the shell
          '';
          hardeningDisable = [ "all" ];

          packages = [ pkgs.go pkgs.pkg-config self.packages.${system}.default ];
          # For `go build` and `go test` of the cgo package in the shell.
          buildInputs = [ pkgs.zstd ];
        };
      });
    };
}
