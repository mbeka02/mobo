-- name: GetDashboardStats :one
SELECT
  COALESCE(SUM(p.amount), 0)::text as total_revenue,
  COALESCE(SUM(r.number_of_seats), 0)::bigint as tickets_sold,
  COUNT(DISTINCT s.venue_id)::int as active_venues,
  (SELECT COUNT(*) FROM movies WHERE deleted_at IS NULL)::int as active_movies
FROM payments p
JOIN reservations r ON r.id = p.reservation_id
JOIN showtimes s ON s.id = r.showtime_id
WHERE p.payment_status = 'completed'
  AND r.deleted_at IS NULL;

-- name: GetMonthlyRevenue :many
SELECT
  DATE_TRUNC('month', p.paid_at)::date as month,
  COALESCE(SUM(p.amount), 0)::text as revenue
FROM payments p
WHERE p.payment_status = 'completed'
  AND p.paid_at >= DATE_TRUNC('month', NOW()) - INTERVAL '5 months'
GROUP BY DATE_TRUNC('month', p.paid_at)
ORDER BY month ASC;
