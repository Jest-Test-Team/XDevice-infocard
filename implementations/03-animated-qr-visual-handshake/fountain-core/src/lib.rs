pub mod decoder;
pub mod encoder;
pub mod symbol;

#[cfg(test)]
mod tests {
    use crate::decoder::decode;
    use crate::encoder::encode;
    use crate::symbol::Symbol;

    #[test]
    fn roundtrip_placeholder_encode_decode() {
        let payload = b"visual-handshake";
        let symbols = encode(payload, 6);
        let decoded = decode(&symbols).expect("decode should succeed");
        assert_eq!(decoded, payload);
    }

    #[test]
    fn decode_rejects_incomplete_symbol_set() {
        let payload = b"visual-handshake";
        let mut symbols = encode(payload, 4);
        symbols.pop();
        symbols.pop();
        let err = decode(&symbols).expect_err("decode should fail");
        assert_eq!(err, "incomplete symbol set");
    }

    #[test]
    fn decode_rejects_mixed_transfer_ids() {
        let symbols = vec![
            Symbol::new(1, 10, 0, 2, b"hello".to_vec()),
            Symbol::new(2, 10, 1, 2, b"world".to_vec()),
        ];
        let err = decode(&symbols).expect_err("decode should fail");
        assert_eq!(err, "mixed transfer ids");
    }

    #[test]
    fn decode_accepts_out_of_order_symbols() {
        let symbols = vec![
            Symbol::new(7, 3, 2, 4, b"c".to_vec()),
            Symbol::new(7, 3, 0, 4, b"a".to_vec()),
            Symbol::new(7, 3, 3, 4, b"`".to_vec()),
            Symbol::new(7, 3, 1, 4, b"b".to_vec()),
        ];
        let decoded = decode(&symbols).expect("decode should succeed");
        assert_eq!(decoded, b"abc");
    }

    #[test]
    fn decode_accepts_identical_duplicate_symbols() {
        let symbols = vec![
            Symbol::new(9, 6, 0, 3, b"foo".to_vec()),
            Symbol::new(9, 6, 0, 3, b"foo".to_vec()),
            Symbol::new(9, 6, 1, 3, b"bar".to_vec()),
            Symbol::new(9, 6, 2, 3, b"\x04\x0e\x1d".to_vec()),
        ];
        let decoded = decode(&symbols).expect("decode should succeed");
        assert_eq!(decoded, b"foobar");
    }

    #[test]
    fn decode_rejects_conflicting_duplicate_symbols() {
        let symbols = vec![
            Symbol::new(11, 6, 0, 3, b"foo".to_vec()),
            Symbol::new(11, 6, 0, 3, b"zzz".to_vec()),
            Symbol::new(11, 6, 1, 3, b"bar".to_vec()),
            Symbol::new(11, 6, 2, 3, b"\x04\x0e\x1d".to_vec()),
        ];
        let err = decode(&symbols).expect_err("decode should fail");
        assert_eq!(err, "conflicting duplicate symbol");
    }

    #[test]
    fn decode_recovers_one_missing_symbol_with_parity() {
        let payload = b"visual-handshake";
        let mut symbols = encode(payload, 4);
        symbols.remove(1);
        let decoded = decode(&symbols).expect("decode should recover");
        assert_eq!(decoded, payload);
    }

    #[test]
    fn decode_rejects_mixed_payload_lengths() {
        let symbols = vec![
            Symbol::new(21, 6, 0, 3, b"foo".to_vec()),
            Symbol::new(21, 7, 1, 3, b"bar".to_vec()),
            Symbol::new(21, 6, 2, 3, b"\x04\x0e\x1d".to_vec()),
        ];
        let err = decode(&symbols).expect_err("decode should fail");
        assert_eq!(err, "mixed payload lengths");
    }

    #[test]
    fn decode_rejects_invalid_total_count_zero() {
        let symbols = vec![Symbol::new(31, 3, 0, 0, b"abc".to_vec())];
        let err = decode(&symbols).expect_err("decode should fail");
        assert_eq!(err, "invalid total count");
    }
}
