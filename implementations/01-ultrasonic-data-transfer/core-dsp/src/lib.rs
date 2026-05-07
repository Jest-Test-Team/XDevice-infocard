pub mod demodulator;
pub mod framing;
pub mod modulator;
pub mod protocol;

pub use demodulator::{demodulate_samples, DemodulatorConfig};
pub use modulator::{modulate_bytes, ModulatorConfig};
pub use protocol::{decode_payload, encode_payload, ProtocolError};

#[cfg(test)]
mod tests {
    use super::{decode_payload, encode_payload};

    #[test]
    fn round_trip_payload() {
        let input = "hello-sonic";
        let frame = encode_payload(input).expect("encode should succeed");
        let output = decode_payload(&frame).expect("decode should succeed");
        assert_eq!(input, output);
    }
}
