#[derive(Debug, Clone, PartialEq, Eq)]
pub struct Symbol {
    pub sequence: u32,
    pub total: u32,
    pub data: Vec<u8>,
}

impl Symbol {
    pub fn new(sequence: u32, total: u32, data: Vec<u8>) -> Self {
        Self { sequence, total, data }
    }
}
