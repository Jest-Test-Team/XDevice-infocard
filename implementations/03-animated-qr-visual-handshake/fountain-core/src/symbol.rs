#[derive(Debug, Clone, PartialEq, Eq)]
pub struct Symbol {
    pub transfer_id: u64,
    pub sequence: u32,
    pub total: u32,
    pub data: Vec<u8>,
}

impl Symbol {
    pub fn new(transfer_id: u64, sequence: u32, total: u32, data: Vec<u8>) -> Self {
        Self {
            transfer_id,
            sequence,
            total,
            data,
        }
    }
}
