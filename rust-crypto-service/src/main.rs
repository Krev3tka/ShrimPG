mod crypto;
mod proto;

use std::{env, net::SocketAddr};

use proto::crypto_service_server::{CryptoService, CryptoServiceServer};
use proto::{
    DecryptWithKeyRequest, DecryptWithKeyResponse, DeriveKeyRequest, DeriveKeyResponse,
    EncryptWithKeyRequest, EncryptWithKeyResponse, GenerateRandomBytesRequest,
    GenerateRandomBytesResponse,
};
use tonic::{Request, Response, Status, transport::Server};

#[derive(Default)]
struct CryptoSvc;

#[tonic::async_trait]
impl CryptoService for CryptoSvc {
    async fn generate_random_bytes(
        &self,
        request: Request<GenerateRandomBytesRequest>,
    ) -> Result<Response<GenerateRandomBytesResponse>, Status> {
        let n = request.into_inner().n as usize;
        Ok(Response::new(GenerateRandomBytesResponse {
            data: crypto::generate_random_bytes(n),
        }))
    }

    async fn derive_key(
        &self,
        request: Request<DeriveKeyRequest>,
    ) -> Result<Response<DeriveKeyResponse>, Status> {
        let req = request.into_inner();
        let key = crypto::derive_key(&req.password, &req.salt)
            .map_err(|e| Status::internal(e.to_string()))?;
        Ok(Response::new(DeriveKeyResponse { key: key.to_vec() }))
    }

    async fn encrypt_with_key(
        &self,
        request: Request<EncryptWithKeyRequest>,
    ) -> Result<Response<EncryptWithKeyResponse>, Status> {
        let req = request.into_inner();
        let (nonce, ciphertext) = crypto::encrypt_with_key(&req.key, &req.plaintext)
            .map_err(|e| Status::internal(e.to_string()))?;
        Ok(Response::new(EncryptWithKeyResponse {
            salt: Vec::new(),
            nonce,
            ciphertext,
        }))
    }

    async fn decrypt_with_key(
        &self,
        request: Request<DecryptWithKeyRequest>,
    ) -> Result<Response<DecryptWithKeyResponse>, Status> {
        let req = request.into_inner();
        let plaintext = crypto::decrypt_with_key(&req.key, &req.nonce, &req.ciphertext)
            .map_err(|e| Status::internal(e.to_string()))?;
        Ok(Response::new(DecryptWithKeyResponse { plaintext }))
    }
}

#[tokio::main]
async fn main() -> Result<(), Box<dyn std::error::Error>> {
    let addr: SocketAddr = env::var("CRYPTO_SERVICE_ADDRESS")
        .unwrap_or_else(|_| "0.0.0.0:50051".to_string())
        .parse()?;

    Server::builder()
        .add_service(CryptoServiceServer::new(CryptoSvc))
        .serve(addr)
        .await?;

    Ok(())
}
