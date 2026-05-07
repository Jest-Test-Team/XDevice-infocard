use std::f32::consts::TAU;

#[derive(Debug, Clone)]
pub struct ModulatorConfig {
    pub sample_rate_hz: u32,
    pub symbol_rate_bps: u32,
}

impl Default for ModulatorConfig {
    fn default() -> Self {
        Self {
            sample_rate_hz: 48_000,
            symbol_rate_bps: 1_200,
        }
    }
}

pub fn modulate_bytes(payload: &[u8], config: &ModulatorConfig) -> Vec<f32> {
    modulate_bytes_with_config(payload, config)
}

const BIT0_FREQ_HZ: f32 = 18_000.0;
const BIT1_FREQ_HZ: f32 = 19_500.0;
const TONE_AMPLITUDE: f32 = 0.8;

fn modulate_bytes_with_config(payload: &[u8], config: &ModulatorConfig) -> Vec<f32> {
    let samples_per_bit = (config.sample_rate_hz / config.symbol_rate_bps.max(1)).max(1) as usize;
    let mut samples = Vec::with_capacity(payload.len() * 8 * samples_per_bit);
    let sample_rate = config.sample_rate_hz.max(1) as f32;
    let mut phase = 0.0_f32;

    for byte in payload {
        for bit_idx in (0..8).rev() {
            let bit = (byte >> bit_idx) & 1;
            let freq_hz = if bit == 1 { BIT1_FREQ_HZ } else { BIT0_FREQ_HZ };
            let phase_step = TAU * (freq_hz / sample_rate);
            for _ in 0..samples_per_bit {
                samples.push(TONE_AMPLITUDE * phase.sin());
                phase += phase_step;
                if phase >= TAU {
                    phase -= TAU;
                }
            }
        }
    }
    samples
}
