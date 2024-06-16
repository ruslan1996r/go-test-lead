SELECT
    c.id,
    c.name,
    c.start_date as start_date,
    c.end_date as end_date,
    c.priority,
    c.lead_capacity,
    l.lead_id,
    l.start_date as lead_start,
    l.end_date as lead_end
FROM clients AS c
LEFT JOIN leads as l on c.id = l.client_id