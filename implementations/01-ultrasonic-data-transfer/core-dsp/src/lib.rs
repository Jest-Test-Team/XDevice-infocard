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
    use std::f32::consts::TAU;

    use super::{
        decode_payload, decode_payload_bytes, demodulate_samples, encode_payload,
        encode_payload_bytes, modulate_bytes, DemodulatorConfig, ModulatorConfig,
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

    #[test]
    fn demodulate_tolerates_moderate_deterministic_noise() {
        let input = b"fsk-noise-check";
        let mut samples = modulate_bytes(input, &ModulatorConfig::default());
        for (i, s) in samples.iter_mut().enumerate() {
            let noise = ((i as u32).wrapping_mul(1_103_515_245).wrapping_add(12_345) % 1000) as f32
                / 1000.0;
            *s += (noise - 0.5) * 0.35;
        }
        let output = demodulate_samples(&samples, &DemodulatorConfig::default());
        assert_eq!(input.to_vec(), output);
    }

    #[test]
    fn demodulate_exposes_bit_errors_under_heavy_corruption() {
        let input: Vec<u8> = vec![0b1000_0000, 0b0101_0101, 0b1100_0011];
        let mut samples = modulate_bytes(&input, &ModulatorConfig::default());
        let cfg = ModulatorConfig::default();
        let samples_per_bit = (cfg.sample_rate_hz / cfg.symbol_rate_bps.max(1)).max(1) as usize;
        let sample_rate = cfg.sample_rate_hz as f32;
        let wrong_freq_hz = 18_000.0_f32;

        for n in 0..samples_per_bit {
            let phase = TAU * wrong_freq_hz * (n as f32 / sample_rate);
            samples[n] = 0.8 * phase.sin();
        }
        let output = demodulate_samples(&samples, &DemodulatorConfig::default());
        assert_eq!(input.len(), output.len());
        assert!(
            output != input,
            "expected at least one decoded bit to differ after heavy corruption"
        );
    }
}
