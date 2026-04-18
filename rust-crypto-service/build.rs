fn main() {
    tonic_build::configure()
        .build_server(true)
        .build_client(false)
        .compile_protos(&["proto/crypto.proto"], &["proto"])
        .expect("failed to compile protos");
}
