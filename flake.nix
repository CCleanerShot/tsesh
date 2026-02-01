{
  description = "A very basic flake";

  inputs = {
    nixpkgs.url = "github:nixos/nixpkgs?ref=nixos-unstable";
  };

  outputs = { self, nixpkgs }: 
  let 
    systems = [
      "x86_64-linux"
    ];
  in 
  {
    devShells = nixpkgs.lib.genAttrs systems (system: 
      let
        pkgs = nixpkgs.legacyPackages.${system};
      in {
        default = pkgs.mkShell {
          packages = with pkgs; [ go tmux opencode ];
        };
      }
    );
  };
}
