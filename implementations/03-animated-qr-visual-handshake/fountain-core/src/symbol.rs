#[derive(Debug, Clone, PartialEq, Eq)]
pub struct Symbol {
    pub transfer_id: u64,
    pub payload_len: usize,
    pub sequence: u32,
    pub total: u32,
    pub neighbors: Vec<u32>,
    pub data: Vec<u8>,
}

impl Symbol {
    pub fn new(
        transfer_id: u64,
        payload_len: usize,
        sequence: u32,
        total: u32,
        neighbors: Vec<u32>,
        data: Vec<u8>,
    ) -> Self {
        Self {
            transfer_id,
            payload_len,
            sequence,
            total,
            neighbors,
            data,
        }
    }
}
