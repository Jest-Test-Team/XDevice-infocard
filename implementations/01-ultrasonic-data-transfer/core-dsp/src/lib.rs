pub mod demodulator;
pub mod framing;
pub mod modulator;
pub mod protocol;

pub use demodulator::{demodulate_samples, DemodulatorConfig};
pub use modulator::{modulate_bytes, ModulatorConfig};
pub use protocol::{
    decode_payload, decode_payload_bytes, encode_payload, encode_payload_bytes, ProtocolError,
};

#[cfg(test)]
mod tests {
    use super::{
        decode_payload, demodulate_samples, encode_payload, modulate_bytes, DemodulatorConfig,
        ModulatorConfig,
    };

    #[test]
    fn round_trip_payload() {
        let input = "hello-sonic";
        let frame = encode_payload(input).expect("encode should succeed");
        let output = decode_payload(&frame).expect("decode should succeed");
        assert_eq!(input, output);
    }

    #[test]
    fn round_trip_modulate_demodulate_bytes() {
        let input: Vec<u8> = (0..=255).collect();
        let samples = modulate_bytes(&input, &ModulatorConfig::default());
        let output = demodulate_samples(&samples, &DemodulatorConfig::default());
        assert_eq!(input, output);
    }

    #[test]
    fn round_trip_binary_payload() {
        let input: Vec<u8> = (0..=255).rev().collect();
        let frame = encode_payload_bytes(&input).expect("encode should succeed");
        let output = decode_payload_bytes(&frame).expect("decode should succeed");
        assert_eq!(input, output);
    }
}
