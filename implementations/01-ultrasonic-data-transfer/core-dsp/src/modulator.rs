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

pub fn modulate_bytes(payload: &[u8], _config: &ModulatorConfig) -> Vec<f32> {
    // MVP transport mapping: one byte to one normalized PCM sample.
    payload
        .iter()
        .map(|b| (*b as f32 / 255.0) * 2.0 - 1.0)
        .collect()
}
