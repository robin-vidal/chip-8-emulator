{ pkgs ? import <nixpkgs> {} }:

pkgs.mkShell {
  buildInputs = with pkgs; [
    go
    libx11
    libxcursor
    libxrandr
    libxinerama
    libxi
    libxext
    xorg.libXxf86vm
    mesa
    libGL
    pkg-config
  ];

  shellHook = ''
    export LD_LIBRARY_PATH=${pkgs.mesa}/lib:${pkgs.libGL}/lib:$LD_LIBRARY_PATH
  '';
}