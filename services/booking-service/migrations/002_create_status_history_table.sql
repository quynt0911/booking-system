-- Create status_history table
CREATE TABLE IF NOT EXISTS booking_status_history (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    booking_id UUID NOT NULL,
    old_status booking_status,
    new_status booking_status NOT NULL,
    changed_by UUID NOT NULL,
    reason TEXT,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    
    -- Add foreign key constraint
    CONSTRAINT fk_booking
        FOREIGN KEY (booking_id)
        REFERENCES bookings(id)
        ON DELETE CASCADE,
    CONSTRAINT fk_changed_by
        FOREIGN KEY (changed_by)
        REFERENCES users(id)
        ON DELETE SET NULL
);

-- Create indexes
CREATE INDEX idx_status_history_booking_id ON booking_status_history(booking_id);
CREATE INDEX idx_status_history_created_at ON booking_status_history(created_at);
CREATE INDEX idx_status_history_changed_by ON booking_status_history(changed_by);

-- Add trigger to automatically create status history entry
CREATE OR REPLACE FUNCTION create_status_history()
RETURNS TRIGGER AS $$
BEGIN
    IF OLD.status IS DISTINCT FROM NEW.status THEN
        INSERT INTO booking_status_history (
            booking_id,
            old_status,
            new_status,
            changed_by,
            reason
        ) VALUES (
            NEW.id,
            OLD.status,
            NEW.status,
            COALESCE(NEW.updated_by, 0), -- Default to 0 if not set
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