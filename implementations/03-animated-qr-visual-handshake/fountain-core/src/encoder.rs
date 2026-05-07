use crate::symbol::Symbol;

pub fn encode(payload: &[u8], frame_size: usize) -> Vec<Symbol> {
    if frame_size == 0 {
        return Vec::new();
    }

    let data_total = payload.len().div_ceil(frame_size) as u32;
    let total = data_total + 1;
    let transfer_id = transfer_id_for(payload);
    let payload_len = payload.len();
    let mut symbols: Vec<Symbol> = payload
        .chunks(frame_size)
        .enumerate()
        .map(|(index, chunk)| {
            Symbol::new(transfer_id, payload_len, index as u32, total, chunk.to_vec())
        })
        .collect();

    // Add one parity symbol so decoder can recover one missing data chunk.
    let parity = build_parity_symbol(transfer_id, payload, frame_size, total, data_total);
    symbols.push(parity);
    symbols
}

fn transfer_id_for(payload: &[u8]) -> u64 {
    let mut hash: u64 = 1469598103934665603;
    for &b in payload {
        hash ^= b as u64;
        hash = hash.wrapping_mul(1099511628211);
    }
    hash
}

fn build_parity_symbol(
    transfer_id: u64,
    payload: &[u8],
    frame_size: usize,
    total: u32,
    data_total: u32,
) -> Symbol {
    let mut parity = vec![0u8; frame_size];
    for chunk in payload.chunks(frame_size) {
        for (i, b) in chunk.iter().enumerate() {
            parity[i] ^= *b;
        }
    }
    Symbol::new(transfer_id, payload.len(), data_total, total, parity)
}
