use crate::symbol::Symbol;

pub fn decode(symbols: &[Symbol]) -> Result<Vec<u8>, &'static str> {
    if symbols.is_empty() {
        return Ok(Vec::new());
    }

    let expected_total = symbols[0].total;
    if symbols.len() as u32 != expected_total {
        return Err("incomplete symbol set");
    }

    let mut ordered = symbols.to_vec();
    ordered.sort_by_key(|s| s.sequence);

    let mut output = Vec::new();
    for symbol in ordered {
        output.extend_from_slice(&symbol.data);
    }

    Ok(output)
}
