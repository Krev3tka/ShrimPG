FROM rust:1.94 AS builder
WORKDIR /app
COPY Cargo.toml Cargo.lock* ./
COPY proto ./proto
COPY src ./src
COPY build.rs ./build.rs
RUN cargo build --release

FROM debian:bookworm-slim
WORKDIR /app
COPY --from=builder /app/target/release/shrimpg-crypto-service /usr/local/bin/shrimpg-crypto-service
EXPOSE 50051
CMD ["shrimpg-crypto-service"]
