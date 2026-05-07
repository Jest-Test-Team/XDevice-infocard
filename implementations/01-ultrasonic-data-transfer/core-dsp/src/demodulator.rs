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

pub fn demodulate_samples(samples: &[f32], config: &DemodulatorConfig) -> Vec<u8> {
    let samples_per_bit = (config.sample_rate_hz / config.symbol_rate_bps.max(1)).max(1) as usize;
    if samples_per_bit == 0 {
        return Vec::new();
    }

    let bit_count = samples.len() / samples_per_bit;
    let mut bits = Vec::with_capacity(bit_count);
    for i in 0..bit_count {
        let start = i * samples_per_bit;
        let end = start + samples_per_bit;
        let avg = samples[start..end].iter().copied().sum::<f32>() / samples_per_bit as f32;
        bits.push(if avg >= 0.0 { 1_u8 } else { 0_u8 });
    }

    let byte_count = bits.len() / 8;
    let mut out = Vec::with_capacity(byte_count);
    for i in 0..byte_count {
        let mut byte = 0_u8;
        for bit in 0..8 {
            byte = (byte << 1) | bits[i * 8 + bit];
        }
        out.push(byte);
    }
    out
}
