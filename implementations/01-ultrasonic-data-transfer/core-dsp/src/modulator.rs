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

fn modulate_bytes_with_config(payload: &[u8], config: &ModulatorConfig) -> Vec<f32> {
    let samples_per_bit = (config.sample_rate_hz / config.symbol_rate_bps.max(1)).max(1) as usize;
    let mut samples = Vec::with_capacity(payload.len() * 8 * samples_per_bit);
    for byte in payload {
        for bit_idx in (0..8).rev() {
            let bit = (byte >> bit_idx) & 1;
            let level = if bit == 1 { 0.8_f32 } else { -0.8_f32 };
            for _ in 0..samples_per_bit {
                samples.push(level);
            }
        }
    }
    samples
}
