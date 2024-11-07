{
  inputs.nixpkgs.url = "github:nixos/nixpkgs/nixos-unstable";

  outputs = { self, nixpkgs }: 
    let
      forAllSystems = function:
      nixpkgs.lib.genAttrs [
      "x86_64-linux"
      "x86_64-darwin"
      "aarch64-linux"
      ] (system:
        function (import nixpkgs {
          inherit system;
        }));
    in {
      packages = forAllSystems (pkgs: {
        default = pkgs.callPackage ./default.nix { } ;
      });
  };
}
