-- Drop trigger
DROP TRIGGER IF EXISTS booking_changes_trigger ON bookings;

-- Drop function
DROP FUNCTION IF EXISTS notify_booking_changes(); 