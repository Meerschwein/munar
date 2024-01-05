let
  pkgs = import <nixpkgs> {system = "x86_64-linux";};
in
  pkgs.mkShell {
    packages = with pkgs; [
      go

      delve
      go-tools
      gofumpt
      gopls

      just
    ];
  }
