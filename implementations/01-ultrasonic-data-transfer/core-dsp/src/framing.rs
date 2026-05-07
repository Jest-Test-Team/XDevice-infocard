pub const PREAMBLE: [u8; 2] = [0xAA, 0x55];

#[derive(Debug, Clone, PartialEq, Eq)]
pub enum FramingError {
    FrameTooShort,
    InvalidPreamble,
    LengthMismatch,
}

pub fn frame_payload(payload: &[u8]) -> Vec<u8> {
    let mut frame = Vec::with_capacity(2 + 2 + payload.len());
    frame.extend_from_slice(&PREAMBLE);
    frame.extend_from_slice(&(payload.len() as u16).to_le_bytes());
    frame.extend_from_slice(payload);
    frame
}

pub fn deframe_payload(frame: &[u8]) -> Result<Vec<u8>, FramingError> {
    if frame.len() < 4 {
        return Err(FramingError::FrameTooShort);
    }

    if frame[0..2] != PREAMBLE {
        return Err(FramingError::InvalidPreamble);
    }

    let declared_len = u16::from_le_bytes([frame[2], frame[3]]) as usize;
    let actual_len = frame.len() - 4;
    if declared_len != actual_len {
        return Err(FramingError::LengthMismatch);
    }

    Ok(frame[4..].to_vec())
}
