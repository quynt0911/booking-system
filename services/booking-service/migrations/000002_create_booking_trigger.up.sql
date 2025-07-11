-- Create notification function
CREATE OR REPLACE FUNCTION notify_booking_changes()
RETURNS trigger AS $$
BEGIN
    IF TG_OP = 'DELETE' THEN
        -- Format: operation:booking_id:expert_id:user_id
        PERFORM pg_notify('booking_changes', 
            TG_OP || ':' || OLD.id || ':' || OLD.expert_id || ':' || OLD.user_id
        );
    END IF;
    RETURN NULL;
END;
$$ LANGUAGE plpgsql;

-- Create trigger
DROP TRIGGER IF EXISTS booking_changes_trigger ON bookings;
CREATE TRIGGER booking_changes_trigger
    AFTER DELETE ON bookings
    FOR EACH ROW
    EXECUTE FUNCTION notify_booking_changes(); 