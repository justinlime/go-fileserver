{ buildGoModule, lib}:
  buildGoModule {
    pname = "go-fileserver";
    version = "0.1.5";
    src = ./.;
    vendorHash = "sha256-escCquo9BzDRj+dH9N7+8v/Qo5+v8LqVGHrJLuZsqog=";
    buildPhase = ''
      go build -o go-fileserver 
      mkdir -p $out/bin
      install -m755 go-fileserver $out/bin/go-filserver
    '';
}
