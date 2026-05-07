use crate::symbol::Symbol;
use std::collections::BTreeMap;

pub fn decode(symbols: &[Symbol]) -> Result<Vec<u8>, &'static str> {
    if symbols.is_empty() {
        return Ok(Vec::new());
    }

    let first = &symbols[0];
    let expected_total = first.total;
    if expected_total == 0 {
        return Err("invalid total count");
    }
    let data_total = expected_total - 1;
    let transfer_id = first.transfer_id;

    let mut ordered: BTreeMap<u32, &Symbol> = BTreeMap::new();
    let mut parity: Option<&Symbol> = None;
    for symbol in symbols {
        if symbol.total != expected_total {
            return Err("inconsistent total count");
        }
        if symbol.transfer_id != transfer_id {
            return Err("mixed transfer ids");
        }
        if symbol.sequence >= expected_total {
            return Err("invalid sequence");
        }
        if symbol.sequence == data_total {
            if let Some(existing) = parity {
                if existing.data != symbol.data {
                    return Err("conflicting duplicate symbol");
                }
            } else {
                parity = Some(symbol);
            }
            continue;
        }
        if let Some(existing) = ordered.get(&symbol.sequence) {
            if existing.data != symbol.data {
                return Err("conflicting duplicate symbol");
            }
            continue;
        }
        ordered.insert(symbol.sequence, symbol);
    }

    if ordered.len() as u32 == data_total {
        return assemble_output(&ordered, data_total);
    }

    if ordered.len() as u32 + 1 == data_total {
        let parity_symbol = parity.ok_or("incomplete symbol set")?;
        let missing_sequence = (0..data_total)
            .find(|seq| !ordered.contains_key(seq))
            .ok_or("missing sequence")?;
        let recovered = recover_missing_chunk(&ordered, parity_symbol, data_total);
        let mut owned: BTreeMap<u32, Vec<u8>> = BTreeMap::new();
        for (k, v) in ordered {
            owned.insert(k, v.data.clone());
        }
        owned.insert(missing_sequence, recovered);
        return assemble_output_owned(&owned, data_total);
    }

    Err("incomplete symbol set")
}

fn recover_missing_chunk(
    ordered: &BTreeMap<u32, &Symbol>,
    parity_symbol: &Symbol,
    data_total: u32,
) -> Vec<u8> {
    let frame_size = parity_symbol.data.len();
    let mut recovered = parity_symbol.data.clone();
    for seq in 0..data_total {
        if let Some(chunk) = ordered.get(&seq) {
            for i in 0..frame_size.min(chunk.data.len()) {
                recovered[i] ^= chunk.data[i];
            }
        }
    }
    recovered
}

fn assemble_output(ordered: &BTreeMap<u32, &Symbol>, data_total: u32) -> Result<Vec<u8>, &'static str> {
    let mut output = Vec::new();
    for expected in 0..data_total {
        let symbol = ordered.get(&expected).ok_or("missing sequence")?;
        output.extend_from_slice(&symbol.data);
    }
    Ok(trim_trailing_padding(output))
}

fn assemble_output_owned(
    ordered: &BTreeMap<u32, Vec<u8>>,
    data_total: u32,
) -> Result<Vec<u8>, &'static str> {
    let mut output = Vec::new();
    for expected in 0..data_total {
        let symbol = ordered.get(&expected).ok_or("missing sequence")?;
        output.extend_from_slice(symbol);
    }
    Ok(trim_trailing_padding(output))
}

fn trim_trailing_padding(mut output: Vec<u8>) -> Vec<u8> {
    while output.last().copied() == Some(0) {
        output.pop();
    }
    output
}
