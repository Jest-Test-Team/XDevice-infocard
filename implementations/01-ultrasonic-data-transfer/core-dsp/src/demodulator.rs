use std::f32::consts::TAU;

const BIT0_FREQ_HZ: f32 = 18_000.0;
const BIT1_FREQ_HZ: f32 = 19_500.0;

#[derive(Debug, Clone, Copy, PartialEq, Eq)]
pub enum DemodMethod {
    Correlator,
    Goertzel,
}

#[derive(Debug, Clone)]
pub struct DemodulatorConfig {
    pub sample_rate_hz: u32,
    pub symbol_rate_bps: u32,
    pub demod_method: DemodMethod,
    pub adaptive_threshold_enabled: bool,
    pub noise_floor_alpha: f32,
    pub noise_floor_scale: f32,
    pub min_energy_gap: f32,
}

impl Default for DemodulatorConfig {
    fn default() -> Self {
        Self {
            sample_rate_hz: 48_000,
            symbol_rate_bps: 1_200,
            demod_method: DemodMethod::Correlator,
            adaptive_threshold_enabled: true,
            noise_floor_alpha: 0.1,
            noise_floor_scale: 2.0,
            min_energy_gap: 1.0,
        }
    }
}

pub fn demodulate_samples(samples: &[f32], config: &DemodulatorConfig) -> Vec<u8> {
    let samples_per_bit = (config.sample_rate_hz / config.symbol_rate_bps.max(1)).max(1) as usize;
    if samples_per_bit == 0 {
        return Vec::new();
    }

    let sample_rate = config.sample_rate_hz.max(1) as f32;
    let bit_count = samples.len() / samples_per_bit;
    let mut bits = Vec::with_capacity(bit_count);
    let mut noise_floor = 0.0_f32;
    let mut last_bit = 0_u8;
    for i in 0..bit_count {
        let start = i * samples_per_bit;
        let end = start + samples_per_bit;
        let symbol = &samples[start..end];
        let energy_0 = tone_energy(symbol, sample_rate, BIT0_FREQ_HZ, config.demod_method);
        let energy_1 = tone_energy(symbol, sample_rate, BIT1_FREQ_HZ, config.demod_method);
        let bit = classify_symbol(energy_0, energy_1, config, &mut noise_floor, last_bit);
        bits.push(bit);
        last_bit = bit;
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

fn classify_symbol(
    energy_0: f32,
    energy_1: f32,
    config: &DemodulatorConfig,
    noise_floor: &mut f32,
    last_bit: u8,
) -> u8 {
    let (strong_energy, weak_energy, strong_bit) = if energy_1 >= energy_0 {
        (energy_1, energy_0, 1_u8)
    } else {
        (energy_0, energy_1, 0_u8)
    };

    if !config.adaptive_threshold_enabled {
        return strong_bit;
    }

    let alpha = config.noise_floor_alpha.clamp(0.0, 1.0);
    *noise_floor = (1.0 - alpha) * *noise_floor + alpha * weak_energy;
    let required_gap = config.min_energy_gap.max(0.0) + config.noise_floor_scale.max(0.0) * *noise_floor;
    let gap = strong_energy - weak_energy;

    if gap < required_gap {
        last_bit
    } else {
        strong_bit
    }
}

fn tone_energy(
    symbol_samples: &[f32],
    sample_rate_hz: f32,
    target_freq_hz: f32,
    method: DemodMethod,
) -> f32 {
    match method {
        DemodMethod::Correlator => tone_energy_correlator(symbol_samples, sample_rate_hz, target_freq_hz),
        DemodMethod::Goertzel => tone_energy_goertzel(symbol_samples, sample_rate_hz, target_freq_hz),
    }
}

fn tone_energy_correlator(symbol_samples: &[f32], sample_rate_hz: f32, target_freq_hz: f32) -> f32 {
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

fn tone_energy_goertzel(symbol_samples: &[f32], sample_rate_hz: f32, target_freq_hz: f32) -> f32 {
    let omega = TAU * target_freq_hz / sample_rate_hz;
    let coeff = 2.0 * omega.cos();
    let mut s_prev = 0.0_f32;
    let mut s_prev2 = 0.0_f32;

    for &x in symbol_samples {
        let s = x + coeff * s_prev - s_prev2;
        s_prev2 = s_prev;
        s_prev = s;
    }

    s_prev2 * s_prev2 + s_prev * s_prev - coeff * s_prev * s_prev2
}

#[cfg(test)]
mod tests {
    use super::{classify_symbol, DemodMethod, DemodulatorConfig};

    #[test]
    fn adaptive_threshold_holds_last_bit_when_symbol_is_ambiguous() {
        let cfg = DemodulatorConfig::default();
        let mut noise_floor = 10.0_f32;
        let bit = classify_symbol(100.0, 104.0, &cfg, &mut noise_floor, 0);
        assert_eq!(0, bit);
    }

    #[test]
    fn non_adaptive_threshold_keeps_stronger_bit() {
        let cfg = DemodulatorConfig {
            adaptive_threshold_enabled: false,
            ..DemodulatorConfig::default()
        };
        let mut noise_floor = 10.0_f32;
        let bit = classify_symbol(100.0, 104.0, &cfg, &mut noise_floor, 0);
        assert_eq!(1, bit);
    }

    #[test]
    fn goertzel_and_correlator_agree_on_clear_tone() {
        let cfg = DemodulatorConfig::default();
        let symbol_len = (cfg.sample_rate_hz / cfg.symbol_rate_bps.max(1)).max(1) as usize;
        let sample_rate = cfg.sample_rate_hz as f32;
        let tone = 19_500.0_f32;
        let mut symbol = Vec::with_capacity(symbol_len);
        for n in 0..symbol_len {
            let phase = std::f32::consts::TAU * tone * (n as f32 / sample_rate);
            symbol.push(0.8 * phase.sin());
        }

        let e_corr = super::tone_energy(&symbol, sample_rate, tone, DemodMethod::Correlator);
        let e_goertzel = super::tone_energy(&symbol, sample_rate, tone, DemodMethod::Goertzel);
        assert!(e_corr > 0.0);
        let rel_err = (e_corr - e_goertzel).abs() / e_corr;
        assert!(rel_err < 0.05, "expected methods to be close, rel_err={rel_err}");
    }
}
