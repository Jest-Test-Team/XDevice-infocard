use crate::symbol::Symbol;
use std::collections::{BTreeMap, BTreeSet};

#[derive(Clone)]
struct Equation {
    vars: BTreeSet<u32>,
    mixed: Vec<u8>,
}

pub fn decode(symbols: &[Symbol]) -> Result<Vec<u8>, &'static str> {
    if symbols.is_empty() {
        return Ok(Vec::new());
    }

    let first = &symbols[0];
    let expected_total = first.total;
    if expected_total == 0 {
        return Err("invalid total count");
    }

    let transfer_id = first.transfer_id;
    let payload_len = first.payload_len;

    let mut unique: BTreeMap<u32, &Symbol> = BTreeMap::new();
    for symbol in symbols {
        validate_symbol_shape(symbol, expected_total, transfer_id, payload_len)?;

        if let Some(existing) = unique.get(&symbol.sequence) {
            if existing.data != symbol.data || existing.neighbors != symbol.neighbors {
                return Err("conflicting duplicate symbol");
            }
            continue;
        }
        unique.insert(symbol.sequence, symbol);
    }

    let data_total = infer_data_total(expected_total, unique.values().copied())?;
    if data_total == 0 {
        return Ok(Vec::new());
    }

    let frame_size = unique
        .values()
        .next()
        .map(|s| s.data.len())
        .ok_or("incomplete symbol set")?;

    let mut resolved: Vec<Option<Vec<u8>>> = vec![None; data_total as usize];
    let mut equations: Vec<Equation> = Vec::new();

    for symbol in unique.values().copied() {
        if symbol.neighbors.is_empty() {
            return Err("invalid neighbors");
        }
        if symbol.data.len() != frame_size {
            return Err("inconsistent symbol size");
        }
        for &n in &symbol.neighbors {
            if n >= data_total {
                return Err("invalid neighbors");
            }
        }

        if symbol.neighbors.len() == 1 {
            let idx = symbol.neighbors[0] as usize;
            insert_or_check(&mut resolved[idx], symbol.data.clone())?;
            continue;
        }

        equations.push(Equation {
            vars: symbol.neighbors.iter().copied().collect(),
            mixed: symbol.data.clone(),
        });
    }

    peel_decode(&mut resolved, &mut equations, frame_size)?;

    if resolved.iter().any(Option::is_none) {
        return Err("incomplete symbol set");
    }

    let mut output = Vec::with_capacity((data_total as usize) * frame_size);
    for chunk in resolved {
        output.extend_from_slice(chunk.as_ref().ok_or("incomplete symbol set")?);
    }

    if output.len() < payload_len {
        return Err("payload length mismatch");
    }
    output.truncate(payload_len);
    Ok(output)
}

fn validate_symbol_shape(
    symbol: &Symbol,
    expected_total: u32,
    transfer_id: u64,
    payload_len: usize,
) -> Result<(), &'static str> {
    if symbol.total != expected_total {
        return Err("inconsistent total count");
    }
    if symbol.transfer_id != transfer_id {
        return Err("mixed transfer ids");
    }
    if symbol.payload_len != payload_len {
        return Err("mixed payload lengths");
    }
    if symbol.sequence >= expected_total {
        return Err("invalid sequence");
    }
    Ok(())
}

fn infer_data_total<'a>(expected_total: u32, symbols: impl Iterator<Item = &'a Symbol>) -> Result<u32, &'static str> {
    let mut max_neighbor: Option<u32> = None;
    for symbol in symbols {
        for &n in &symbol.neighbors {
            max_neighbor = Some(max_neighbor.map_or(n, |m| m.max(n)));
        }
    }

    let max_data_seq = max_neighbor.ok_or("incomplete symbol set")?;
    let data_total = max_data_seq + 1;
    if data_total > expected_total {
        return Err("invalid total count");
    }
    Ok(data_total)
}

fn insert_or_check(slot: &mut Option<Vec<u8>>, chunk: Vec<u8>) -> Result<(), &'static str> {
    if let Some(existing) = slot {
        if existing != &chunk {
            return Err("conflicting duplicate symbol");
        }
        return Ok(());
    }
    *slot = Some(chunk);
    Ok(())
}

fn peel_decode(
    resolved: &mut [Option<Vec<u8>>],
    equations: &mut [Equation],
    frame_size: usize,
) -> Result<(), &'static str> {
    let mut changed = true;
    while changed {
        changed = false;

        for eq in equations.iter_mut() {
            let known: Vec<u32> = eq
                .vars
                .iter()
                .copied()
                .filter(|idx| resolved[*idx as usize].is_some())
                .collect();

            for idx in known {
                if let Some(chunk) = &resolved[idx as usize] {
                    xor_into(&mut eq.mixed, chunk);
                    eq.vars.remove(&idx);
                }
            }

            if eq.vars.is_empty() {
                if eq.mixed.iter().any(|b| *b != 0) {
                    return Err("conflicting duplicate symbol");
                }
                continue;
            }

            if eq.vars.len() == 1 {
                let idx = *eq.vars.iter().next().ok_or("incomplete symbol set")?;
                let slot = &mut resolved[idx as usize];
                insert_or_check(slot, eq.mixed.clone())?;
                eq.vars.clear();
                eq.mixed = vec![0u8; frame_size];
                changed = true;
            }
        }
    }

    Ok(())
}

fn xor_into(target: &mut [u8], src: &[u8]) {
    for (dst, b) in target.iter_mut().zip(src.iter()) {
        *dst ^= *b;
    }
}
