{
  description = "A very basic flake";
  inputs = {
    utils.url = "github:numtide/flake-utils";
  };
  outputs = { self, nixpkgs, utils }: utils.lib.eachSystem utils.lib.allSystems (system: let
    pkgs = nixpkgs.legacyPackages.${system};
  in {
    devShells.default = pkgs.mkShell {
      packages = with pkgs; [
        go_1_18
        go-tools
        gopls
      ];
    };


  });
}
