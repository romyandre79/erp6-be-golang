-- Record Lock Tracking Table
-- Stores active locks on records to prevent concurrent editing (SAP-style)

CREATE TABLE IF NOT EXISTS recordlock (
    recordlockid INT AUTO_INCREMENT PRIMARY KEY,
    tablename VARCHAR(100) NOT NULL COMMENT 'Name of the table being locked (e.g., invoice, purchase_order)',
    recordid INT NOT NULL COMMENT 'ID of the record being locked',
    lockedby INT NOT NULL COMMENT 'User ID who locked the record',
    lockedat TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT 'When the lock was acquired',
    locktype VARCHAR(50) DEFAULT 'edit' COMMENT 'Type of lock: edit, delete, etc.',
    sessionid VARCHAR(255) COMMENT 'Session ID for automatic cleanup',
    
    INDEX idx_table_record (tablename, recordid),
    INDEX idx_locked_by (lockedby),
    INDEX idx_locked_at (lockedat),
    
    UNIQUE KEY unique_lock (tablename, recordid) COMMENT 'Only one lock per record'
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='Tracks active record locks for concurrent editing prevention';
