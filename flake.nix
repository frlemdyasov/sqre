{
  description = "development environmnet for sqre";

  inputs = {
    nixpkgs.url = "github:nixos/nixpkgs?ref=nixos-unstable";
  };

  outputs =
    { self, nixpkgs }:
    {
      devShells.x86_64-linux.default =
        let
          pkgs = nixpkgs.legacyPackages.x86_64-linux;
        in
          pkgs.mkShell {
            packages = with pkgs; [
              cairo
              gdk-pixbuf
              glib
              go
              gobject-introspection
              graphene
              gtk4
              pango
              pkg-config
            ];
          };
    };
}
