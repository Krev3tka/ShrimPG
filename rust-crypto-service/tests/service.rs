use shrimpg_crypto_service::crypto::{
    KEY_LEN, NONCE_LEN, SALT_LEN, decrypt_with_key, derive_key, encrypt_with_key,
};

#[test]
fn crypto_api_roundtrip() {
    let salt = [7u8; SALT_LEN];
    let key = derive_key("password", &salt).expect("key");
    assert_eq!(key.len(), KEY_LEN);

    let (nonce, ciphertext) = encrypt_with_key(&key, b"secret").expect("encrypt");
    assert_eq!(nonce.len(), NONCE_LEN);

    let plaintext = decrypt_with_key(&key, &nonce, &ciphertext).expect("decrypt");
    assert_eq!(plaintext, b"secret");
}
