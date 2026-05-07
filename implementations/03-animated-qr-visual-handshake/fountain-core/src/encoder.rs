use crate::symbol::Symbol;

pub fn encode(payload: &[u8], frame_size: usize) -> Vec<Symbol> {
    if frame_size == 0 {
        return Vec::new();
    }

    let data_total = payload.len().div_ceil(frame_size) as u32;
    if data_total == 0 {
        return Vec::new();
    }

    let repair_total = repair_symbol_count(data_total);
    let total = data_total + repair_total;
    let transfer_id = transfer_id_for(payload);
    let payload_len = payload.len();

    let data_chunks = chunk_with_padding(payload, frame_size);

    let mut symbols: Vec<Symbol> = data_chunks
        .iter()
        .enumerate()
        .map(|(index, chunk)| {
            Symbol::new(
                transfer_id,
                payload_len,
                index as u32,
                total,
                vec![index as u32],
                chunk.clone(),
            )
        })
        .collect();

    for repair_seq in 0..repair_total {
        let sequence = data_total + repair_seq;
        let neighbors = choose_neighbors(transfer_id, data_total, repair_seq);
        let mut mixed = vec![0u8; frame_size];
        for &idx in &neighbors {
            xor_into(&mut mixed, &data_chunks[idx as usize]);
        }
        symbols.push(Symbol::new(
            transfer_id,
            payload_len,
            sequence,
            total,
            neighbors,
            mixed,
        ));
    }

    symbols
}

fn repair_symbol_count(data_total: u32) -> u32 {
    // Light overhead with better-than-single-parity recovery odds.
    data_total.clamp(2, 8)
}

fn transfer_id_for(payload: &[u8]) -> u64 {
    let mut hash: u64 = 1469598103934665603;
    for &b in payload {
        hash ^= b as u64;
        hash = hash.wrapping_mul(1099511628211);
    }
    hash
}

fn chunk_with_padding(payload: &[u8], frame_size: usize) -> Vec<Vec<u8>> {
    payload
        .chunks(frame_size)
        .map(|chunk| {
            let mut out = vec![0u8; frame_size];
            out[..chunk.len()].copy_from_slice(chunk);
            out
        })
        .collect()
}

fn choose_neighbors(transfer_id: u64, data_total: u32, repair_seq: u32) -> Vec<u32> {
    if data_total == 1 {
        return vec![0];
    }

    let mut state = transfer_id ^ ((repair_seq as u64).wrapping_mul(0x9E37_79B9_7F4A_7C15));
    let max_degree = data_total.min(4);
    let degree = 2 + (next_u32(&mut state) % (max_degree - 1));

    let mut chosen = Vec::with_capacity(degree as usize);
    while chosen.len() < degree as usize {
        let candidate = next_u32(&mut state) % data_total;
        if !chosen.contains(&candidate) {
            chosen.push(candidate);
        }
    }
    chosen.sort_unstable();
    chosen
}

fn next_u32(state: &mut u64) -> u32 {
    *state ^= *state >> 12;
    *state ^= *state << 25;
    *state ^= *state >> 27;
    ((*state).wrapping_mul(2685821657736338717) >> 32) as u32
}

fn xor_into(target: &mut [u8], src: &[u8]) {
    for (dst, b) in target.iter_mut().zip(src.iter()) {
        *dst ^= *b;
    }
}
