-- Create status_history table
CREATE TABLE IF NOT EXISTS status_history (
    id SERIAL PRIMARY KEY,
    booking_id INTEGER NOT NULL,
    old_status booking_status,
    new_status booking_status NOT NULL,
    changed_by INTEGER NOT NULL, -- ID of user/expert who changed the status
    changed_by_type VARCHAR(10) NOT NULL, -- 'user' or 'expert'
    reason TEXT,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    
    -- Add foreign key constraint
    CONSTRAINT fk_booking
        FOREIGN KEY (booking_id)
        REFERENCES bookings(id)
        ON DELETE CASCADE
);

-- Create indexes
CREATE INDEX idx_status_history_booking_id ON status_history(booking_id);
CREATE INDEX idx_status_history_created_at ON status_history(created_at);

-- Add trigger to automatically create status history entry
CREATE OR REPLACE FUNCTION create_status_history()
RETURNS TRIGGER AS $$
BEGIN
    IF OLD.status IS DISTINCT FROM NEW.status THEN
        INSERT INTO status_history (
            booking_id,
            old_status,
            new_status,
            changed_by,
            changed_by_type,
            reason
        ) VALUES (
            NEW.id,
            OLD.status,
            NEW.status,
            COALESCE(NEW.updated_by, 0), -- Default to 0 if not set
            COALESCE(NEW.updated_by_type, 'system'),
            NEW.status_change_reason
        );
    END IF;
    RETURN NEW;
END;
$$ language 'plpgsql';

-- Add columns for tracking who changed the status
ALTER TABLE bookings
ADD COLUMN updated_by INTEGER,
ADD COLUMN updated_by_type VARCHAR(10),
ADD COLUMN status_change_reason TEXT;

-- Create trigger for status history
CREATE TRIGGER booking_status_change
    AFTER UPDATE ON bookings
    FOR EACH ROW
    EXECUTE FUNCTION create_status_history(); 