#[derive(Debug, Clone)]
pub struct DemodulatorConfig {
    pub sample_rate_hz: u32,
    pub symbol_rate_bps: u32,
}

impl Default for DemodulatorConfig {
    fn default() -> Self {
        Self {
            sample_rate_hz: 48_000,
            symbol_rate_bps: 1_200,
        }
    }
}

pub fn demodulate_samples(samples: &[f32], _config: &DemodulatorConfig) -> Vec<u8> {
    // Starter placeholder: reverse the simple byte<->sample mapping.
    samples
        .iter()
        .map(|s| (((s + 1.0) / 2.0) * 255.0).clamp(0.0, 255.0) as u8)
        .collect()
}
