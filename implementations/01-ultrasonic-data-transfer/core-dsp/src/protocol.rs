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
    encode_payload_bytes(raw)
}

pub fn decode_payload(frame: &[u8]) -> Result<String, ProtocolError> {
    let raw = decode_payload_bytes(frame)?;
    String::from_utf8(raw).map_err(|_| ProtocolError::InvalidUtf8)
}

pub fn encode_payload_bytes(payload: &[u8]) -> Result<Vec<u8>, ProtocolError> {
    if payload.len() > MAX_PAYLOAD_LEN {
        return Err(ProtocolError::PayloadTooLarge);
    }
    Ok(frame_payload(payload))
}

pub fn decode_payload_bytes(frame: &[u8]) -> Result<Vec<u8>, ProtocolError> {
    deframe_payload(frame).map_err(|_| ProtocolError::InvalidFrame)
}

#[cfg(test)]
mod tests {
    use super::{decode_payload, decode_payload_bytes, encode_payload_bytes, ProtocolError};

    #[test]
    fn decode_payload_bytes_rejects_invalid_frame() {
        let err = decode_payload_bytes(&[0x00, 0x01, 0x02]).expect_err("should fail");
        assert_eq!(err, ProtocolError::InvalidFrame);
    }

    #[test]
    fn decode_payload_rejects_non_utf8_payload() {
        let frame = encode_payload_bytes(&[0xff, 0xfe, 0xfd]).expect("encode should work");
        let err = decode_payload(&frame).expect_err("should fail utf8 decode");
        assert_eq!(err, ProtocolError::InvalidUtf8);
    }
}
