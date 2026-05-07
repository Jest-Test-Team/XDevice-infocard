use crate::symbol::Symbol;
use std::collections::BTreeMap;

pub fn decode(symbols: &[Symbol]) -> Result<Vec<u8>, &'static str> {
    if symbols.is_empty() {
        return Ok(Vec::new());
    }

    let first = &symbols[0];
    let expected_total = first.total;
    let transfer_id = first.transfer_id;

    let mut ordered: BTreeMap<u32, &Symbol> = BTreeMap::new();
    for symbol in symbols {
        if symbol.total != expected_total {
            return Err("inconsistent total count");
        }
        if symbol.transfer_id != transfer_id {
            return Err("mixed transfer ids");
        }
        if let Some(existing) = ordered.get(&symbol.sequence) {
            if existing.data != symbol.data {
                return Err("conflicting duplicate symbol");
            }
            continue;
        }
        ordered.insert(symbol.sequence, symbol);
    }

    if ordered.len() as u32 != expected_total {
        return Err("incomplete symbol set");
    }

    for expected in 0..expected_total {
        if !ordered.contains_key(&expected) {
            return Err("missing sequence");
        }
    }

    let mut output = Vec::new();
    for symbol in ordered {
        output.extend_from_slice(&symbol.1.data);
    }

    Ok(output)
}
