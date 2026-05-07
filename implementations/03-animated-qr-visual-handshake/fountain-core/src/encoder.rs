use crate::symbol::Symbol;

pub fn encode(payload: &[u8], frame_size: usize) -> Vec<Symbol> {
    if frame_size == 0 {
        return Vec::new();
    }

    let total = payload.len().div_ceil(frame_size) as u32;
    let transfer_id = transfer_id_for(payload);
    payload
        .chunks(frame_size)
        .enumerate()
        .map(|(index, chunk)| Symbol::new(transfer_id, index as u32, total, chunk.to_vec()))
        .collect()
}

fn transfer_id_for(payload: &[u8]) -> u64 {
    let mut hash: u64 = 1469598103934665603;
    for &b in payload {
        hash ^= b as u64;
        hash = hash.wrapping_mul(1099511628211);
    }
    hash
}
