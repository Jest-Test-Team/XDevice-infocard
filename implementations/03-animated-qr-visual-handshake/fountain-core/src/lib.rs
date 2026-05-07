pub mod decoder;
pub mod encoder;
pub mod symbol;

#[cfg(test)]
mod tests {
    use crate::decoder::decode;
    use crate::encoder::encode;

    #[test]
    fn roundtrip_placeholder_encode_decode() {
        let payload = b"visual-handshake";
        let symbols = encode(payload, 6);
        let decoded = decode(&symbols).expect("decode should succeed");
        assert_eq!(decoded, payload);
    }
}
