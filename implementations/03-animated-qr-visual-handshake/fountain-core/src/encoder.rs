use crate::symbol::Symbol;

pub fn encode(payload: &[u8], frame_size: usize) -> Vec<Symbol> {
    if frame_size == 0 {
        return Vec::new();
    }

    let total = payload.len().div_ceil(frame_size) as u32;
    payload
        .chunks(frame_size)
        .enumerate()
        .map(|(index, chunk)| Symbol::new(index as u32, total, chunk.to_vec()))
        .collect()
}
