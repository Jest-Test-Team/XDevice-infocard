pub mod decoder;
pub mod encoder;
pub mod symbol;

#[cfg(test)]
mod tests {
    use crate::decoder::decode;
    use crate::encoder::encode;
    use crate::symbol::Symbol;

    #[test]
    fn roundtrip_encode_decode() {
        let payload = b"visual-handshake";
        let symbols = encode(payload, 6);
        let decoded = decode(&symbols).expect("decode should succeed");
        assert_eq!(decoded, payload);
    }

    #[test]
    fn decode_recovers_multiple_missing_symbols_with_mixed_repair() {
        let payload = b"visual-handshake-runtime-integration";
        let mut symbols = encode(payload, 4);
        symbols.retain(|s| s.sequence != 1 && s.sequence != 3);
        let decoded = decode(&symbols).expect("decode should recover");
        assert_eq!(decoded, payload);
    }

    #[test]
    fn decode_rejects_incomplete_symbol_set() {
        let payload = b"visual-handshake";
        let mut symbols = encode(payload, 4);
        symbols.retain(|s| s.sequence >= 4);
        let err = decode(&symbols).expect_err("decode should fail");
        assert_eq!(err, "incomplete symbol set");
    }

    #[test]
    fn decode_rejects_mixed_transfer_ids() {
        let symbols = vec![
            Symbol::new(1, 10, 0, 3, vec![0], b"hello".to_vec()),
            Symbol::new(2, 10, 1, 3, vec![1], b"world".to_vec()),
        ];
        let err = decode(&symbols).expect_err("decode should fail");
        assert_eq!(err, "mixed transfer ids");
    }

    #[test]
    fn decode_accepts_out_of_order_symbols() {
        let symbols = vec![
            Symbol::new(7, 3, 2, 3, vec![2], b"c".to_vec()),
            Symbol::new(7, 3, 0, 3, vec![0], b"a".to_vec()),
            Symbol::new(7, 3, 1, 3, vec![1], b"b".to_vec()),
        ];
        let decoded = decode(&symbols).expect("decode should succeed");
        assert_eq!(decoded, b"abc");
    }

    #[test]
    fn decode_accepts_identical_duplicate_symbols() {
        let symbols = vec![
            Symbol::new(9, 6, 0, 3, vec![0], b"foo".to_vec()),
            Symbol::new(9, 6, 0, 3, vec![0], b"foo".to_vec()),
            Symbol::new(9, 6, 1, 3, vec![1], b"bar".to_vec()),
        ];
        let decoded = decode(&symbols).expect("decode should succeed");
        assert_eq!(decoded, b"foobar");
    }

    #[test]
    fn decode_rejects_conflicting_duplicate_symbols() {
        let symbols = vec![
            Symbol::new(11, 6, 0, 3, vec![0], b"foo".to_vec()),
            Symbol::new(11, 6, 0, 3, vec![0], b"zzz".to_vec()),
            Symbol::new(11, 6, 1, 3, vec![1], b"bar".to_vec()),
        ];
        let err = decode(&symbols).expect_err("decode should fail");
        assert_eq!(err, "conflicting duplicate symbol");
    }

    #[test]
    fn decode_rejects_mixed_payload_lengths() {
        let symbols = vec![
            Symbol::new(21, 6, 0, 3, vec![0], b"foo".to_vec()),
            Symbol::new(21, 7, 1, 3, vec![1], b"bar".to_vec()),
        ];
        let err = decode(&symbols).expect_err("decode should fail");
        assert_eq!(err, "mixed payload lengths");
    }

    #[test]
    fn decode_rejects_invalid_total_count_zero() {
        let symbols = vec![Symbol::new(31, 3, 0, 0, vec![0], b"abc".to_vec())];
        let err = decode(&symbols).expect_err("decode should fail");
        assert_eq!(err, "invalid total count");
    }

    #[test]
    fn decode_preserves_trailing_zero_bytes_with_complete_set() {
        let payload = vec![0x51, 0x52, 0x53, 0x00, 0x00];
        let symbols = encode(&payload, 4);
        let decoded = decode(&symbols).expect("decode should succeed");
        assert_eq!(decoded, payload);
    }

    #[test]
    fn decode_preserves_trailing_zero_bytes_when_recovering_missing_chunks() {
        let payload = vec![0x41, 0x42, 0x43, 0x44, 0x00, 0x00, 0x00, 0x00];
        let mut symbols = encode(&payload, 4);
        symbols.retain(|s| s.sequence != 1);
        let decoded = decode(&symbols).expect("decode should recover");
        assert_eq!(decoded, payload);
    }
}
