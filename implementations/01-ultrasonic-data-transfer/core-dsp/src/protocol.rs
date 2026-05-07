use crate::framing::{deframe_payload, frame_payload, FramingError};

#[derive(Debug, Clone, PartialEq, Eq)]
pub enum ProtocolError {
    InvalidFrame,
    InvalidUtf8,
    PayloadTooLarge,
}

impl From<FramingError> for ProtocolError {
    fn from(_: FramingError) -> Self {
        ProtocolError::InvalidFrame
    }
}

const MAX_PAYLOAD_LEN: usize = u16::MAX as usize;

pub fn encode_payload(payload: &str) -> Result<Vec<u8>, ProtocolError> {
    let raw = payload.as_bytes();
    if raw.len() > MAX_PAYLOAD_LEN {
        return Err(ProtocolError::PayloadTooLarge);
    }
    Ok(frame_payload(raw))
}

pub fn decode_payload(frame: &[u8]) -> Result<String, ProtocolError> {
    let raw = deframe_payload(frame)?;
    String::from_utf8(raw).map_err(|_| ProtocolError::InvalidUtf8)
}
