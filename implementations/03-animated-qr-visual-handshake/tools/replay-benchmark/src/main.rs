use fountain_core::decoder::decode;
use fountain_core::encoder::encode;

fn should_drop(seq: u32, percent: u8) -> bool {
    if percent == 0 {
        return false;
    }
    (seq.wrapping_mul(37).wrapping_add(17) % 100) < percent as u32
}

fn parse_drop_list(raw: Option<String>) -> Vec<u8> {
    let Some(raw) = raw else {
        return vec![0, 10, 20, 30, 40, 50];
    };

    let mut out: Vec<u8> = raw
        .split(',')
        .filter_map(|v| v.trim().parse::<u8>().ok())
        .map(|v| v.min(99))
        .collect();

    if out.is_empty() {
        out = vec![0, 10, 20, 30, 40, 50];
    }
    out.sort_unstable();
    out.dedup();
    out
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
    let drop_percents = parse_drop_list(args.next());
    let runs = args
        .next()
        .and_then(|v| v.parse::<u32>().ok())
        .unwrap_or(20)
        .max(1);

    let encoded = encode(payload.as_bytes(), frame_size);

    println!("payload_len={}", payload.len());
    println!("frame_size={}", frame_size);
    println!("encoded_symbols={}", encoded.len());
    println!("runs_per_drop={}", runs);
    println!("drop_percents={:?}", drop_percents);
    println!("drop_percent,recovery_rate,decoded_match_rate,avg_received_symbols");

    for drop_percent in drop_percents {
        let mut success_count = 0u32;
        let mut match_count = 0u32;
        let mut received_total = 0u32;

        for run in 0..runs {
            let sent: Vec<_> = encoded
                .iter()
                .filter(|s| !should_drop(s.sequence.wrapping_add(run), drop_percent))
                .cloned()
                .collect();

            received_total += sent.len() as u32;
            match decode(&sent) {
                Ok(decoded) => {
                    success_count += 1;
                    if decoded == payload.as_bytes() {
                        match_count += 1;
                    }
                }
                Err(_) => {}
            }
        }

        let recovery_rate = success_count as f64 / runs as f64;
        let match_rate = match_count as f64 / runs as f64;
        let avg_received = received_total as f64 / runs as f64;

        println!(
            "{},{:.3},{:.3},{:.2}",
            drop_percent, recovery_rate, match_rate, avg_received
        );
    }
}
