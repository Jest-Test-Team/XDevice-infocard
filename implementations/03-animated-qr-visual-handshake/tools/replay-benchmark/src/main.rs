use fountain_core::decoder::decode;
use fountain_core::encoder::encode;

fn should_drop(seq: u32, percent: u8) -> bool {
    if percent == 0 {
        return false;
    }
    (seq.wrapping_mul(37).wrapping_add(17) % 100) < percent as u32
}

fn main() {
    let mut args = std::env::args().skip(1);
    let payload = args
        .next()
        .unwrap_or_else(|| "animated-qr-visual-handshake".to_string());
    let frame_size = args
        .next()
        .and_then(|v| v.parse::<usize>().ok())
        .unwrap_or(12);
    let drop_percent = args
        .next()
        .and_then(|v| v.parse::<u8>().ok())
        .unwrap_or(20)
        .min(99);

    let encoded = encode(payload.as_bytes(), frame_size);
    let sent: Vec<_> = encoded
        .iter()
        .filter(|s| !should_drop(s.sequence, drop_percent))
        .cloned()
        .collect();

    let result = decode(&sent);
    println!("payload_len={}", payload.len());
    println!("frame_size={}", frame_size);
    println!("drop_percent={}", drop_percent);
    println!("encoded_symbols={}", encoded.len());
    println!("received_symbols={}", sent.len());

    match result {
        Ok(decoded) => {
            println!("decode_success=true");
            println!("decoded_matches={}", decoded == payload.as_bytes());
        }
        Err(err) => {
            println!("decode_success=false");
            println!("error={}", err);
        }
    }
}
