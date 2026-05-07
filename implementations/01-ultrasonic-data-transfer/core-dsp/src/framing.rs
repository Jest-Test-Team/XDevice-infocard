pub const PREAMBLE: [u8; 2] = [0xAA, 0x55];
pub const PROTOCOL_VERSION: u8 = 1;

#[derive(Debug, Clone, PartialEq, Eq)]
pub enum FramingError {
    FrameTooShort,
    InvalidPreamble,
    InvalidVersion,
    LengthMismatch,
    ChecksumMismatch,
}

pub fn frame_payload(payload: &[u8]) -> Vec<u8> {
    let checksum = checksum32(payload);
    let mut frame = Vec::with_capacity(2 + 1 + 2 + payload.len() + 4);
    frame.extend_from_slice(&PREAMBLE);
    frame.push(PROTOCOL_VERSION);
    frame.extend_from_slice(&(payload.len() as u16).to_le_bytes());
    frame.extend_from_slice(payload);
    frame.extend_from_slice(&checksum.to_le_bytes());
    frame
}

pub fn deframe_payload(frame: &[u8]) -> Result<Vec<u8>, FramingError> {
    if frame.len() < 9 {
        return Err(FramingError::FrameTooShort);
    }

    if frame[0..2] != PREAMBLE {
        return Err(FramingError::InvalidPreamble);
    }

    if frame[2] != PROTOCOL_VERSION {
        return Err(FramingError::InvalidVersion);
    }

    let declared_len = u16::from_le_bytes([frame[3], frame[4]]) as usize;
    let actual_len = frame.len() - 9;
    if declared_len != actual_len {
        return Err(FramingError::LengthMismatch);
    }

    let payload = &frame[5..(5 + declared_len)];
    let declared_checksum = u32::from_le_bytes([
        frame[5 + declared_len],
        frame[6 + declared_len],
        frame[7 + declared_len],
        frame[8 + declared_len],
    ]);
    let actual_checksum = checksum32(payload);
    if declared_checksum != actual_checksum {
        return Err(FramingError::ChecksumMismatch);
    }

    Ok(payload.to_vec())
}

fn checksum32(data: &[u8]) -> u32 {
    let mut sum: u32 = 0;
    for &b in data {
        sum = sum.wrapping_add(b as u32);
    }
    sum
}
