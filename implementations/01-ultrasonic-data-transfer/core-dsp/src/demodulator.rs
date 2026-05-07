use std::f32::consts::TAU;

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

    const BIT0_FREQ_HZ: f32 = 18_000.0;
    const BIT1_FREQ_HZ: f32 = 19_500.0;

    let sample_rate = config.sample_rate_hz.max(1) as f32;
    let bit_count = samples.len() / samples_per_bit;
    let mut bits = Vec::with_capacity(bit_count);
    for i in 0..bit_count {
        let start = i * samples_per_bit;
        let end = start + samples_per_bit;
        let energy_0 = tone_energy(&samples[start..end], sample_rate, BIT0_FREQ_HZ);
        let energy_1 = tone_energy(&samples[start..end], sample_rate, BIT1_FREQ_HZ);
        bits.push(if energy_1 >= energy_0 { 1_u8 } else { 0_u8 });
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

fn tone_energy(symbol_samples: &[f32], sample_rate_hz: f32, target_freq_hz: f32) -> f32 {
    let phase_step = TAU * target_freq_hz / sample_rate_hz;
    let mut i_acc = 0.0_f32;
    let mut q_acc = 0.0_f32;

    for (n, &s) in symbol_samples.iter().enumerate() {
        let phase = phase_step * n as f32;
        i_acc += s * phase.cos();
        q_acc -= s * phase.sin();
    }
    i_acc * i_acc + q_acc * q_acc
}
