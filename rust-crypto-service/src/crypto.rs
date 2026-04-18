use aes_gcm::{
    aead::{Aead, KeyInit},
    Aes256Gcm, Nonce,
};
use argon2::{Algorithm, Argon2, Params, Version};
use rand::{rngs::OsRng, RngCore};
use thiserror::Error;

#[allow(dead_code)]
pub const SALT_LEN: usize = 16;
pub const NONCE_LEN: usize = 12;
pub const KEY_LEN: usize = 32;
const ARGON2_MEMORY_KIB: u32 = 64 * 1024;
const ARGON2_ITERATIONS: u32 = 3;
const ARGON2_LANES: u32 = 2;

#[derive(Debug, Error)]
pub enum CryptoError {
    #[error("key derivation failed")]
    KeyDerivation,
    #[error("encryption failed")]
    EncryptionFailed,
    #[error("decryption failed")]
    DecryptionFailed,
    #[error("invalid input")]
    InvalidInput,
}

pub fn generate_random_bytes(n: usize) -> Vec<u8> {
    let mut data = vec![0u8; n];
    OsRng.fill_bytes(&mut data);
    data
}

pub fn derive_key(password: &str, salt: &[u8]) -> Result<[u8; KEY_LEN], CryptoError> {
    let params = Params::new(
        ARGON2_MEMORY_KIB,
        ARGON2_ITERATIONS,
        ARGON2_LANES,
        Some(KEY_LEN),
    )
    .map_err(|_| CryptoError::KeyDerivation)?;
    let argon2 = Argon2::new(Algorithm::Argon2i, Version::V0x13, params);
    let mut key = [0u8; KEY_LEN];
    argon2
        .hash_password_into(password.as_bytes(), salt, &mut key)
        .map_err(|_| CryptoError::KeyDerivation)?;
    Ok(key)
}

pub fn encrypt_with_key(key: &[u8], plaintext: &[u8]) -> Result<(Vec<u8>, Vec<u8>), CryptoError> {
    if key.len() != KEY_LEN {
        return Err(CryptoError::InvalidInput);
    }

    let nonce = generate_random_bytes(NONCE_LEN);
    let cipher = Aes256Gcm::new_from_slice(key).map_err(|_| CryptoError::EncryptionFailed)?;
    let ciphertext = cipher
        .encrypt(Nonce::from_slice(&nonce), plaintext)
        .map_err(|_| CryptoError::EncryptionFailed)?;
    Ok((nonce, ciphertext))
}

pub fn decrypt_with_key(
    key: &[u8],
    nonce: &[u8],
    ciphertext: &[u8],
) -> Result<Vec<u8>, CryptoError> {
    if key.len() != KEY_LEN || nonce.len() != NONCE_LEN {
        return Err(CryptoError::InvalidInput);
    }

    let cipher = Aes256Gcm::new_from_slice(key).map_err(|_| CryptoError::DecryptionFailed)?;
    cipher
        .decrypt(Nonce::from_slice(nonce), ciphertext)
        .map_err(|_| CryptoError::DecryptionFailed)
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn roundtrip_works() {
        let salt = [1u8; SALT_LEN];
        let key = derive_key("password", &salt).expect("key");
        let (nonce, ciphertext) = encrypt_with_key(&key, b"hello").expect("encrypt");
        let plaintext = decrypt_with_key(&key, &nonce, &ciphertext).expect("decrypt");
        assert_eq!(plaintext, b"hello");
    }
}
