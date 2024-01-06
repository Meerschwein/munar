let
  pkgs = import <nixpkgs> {system = "x86_64-linux";};
in
  pkgs.mkShell {
    packages = with pkgs; [
      delve
      go
      go-tools
      gofumpt
      gopls
      govulncheck
    ];
  }
