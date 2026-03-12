fn main() {
    let protoc = protoc_bin_vendored::protoc_bin_path().expect("protoc");
    unsafe {
        std::env::set_var("PROTOC", protoc);
    }

    tonic_build::configure()
        .build_server(true)
        .build_client(true)
        .compile_protos(
            &["../../packages/shared-proto/customs/platform/v1/platform.proto"],
            &["../../packages/shared-proto"],
        )
        .expect("compile protos");
}
