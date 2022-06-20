-- +goose Up
-- +goose StatementBegin
-- helper function for generating random datetime
Create or replace function random_datetime() returns text as
$$
declare
begin
    return timestamp '2022-06-20' +
           random() * (timestamp '2023-01-01' -
                       timestamp '2022-06-20');
end;
$$ language plpgsql;

-- insert N users
insert into users (first_name, second_name, email, user_zone)
select left(md5(random()::text), 5)                                                   as first_name,
       left(md5(random()::text), 7)                                                   as second_name,
       concat(left(md5(random()::text), 4), '@', left(md5(random()::text), 5), '.ru') as email,
       'Europe/Moscow'                                                                as user_zone
from generate_series(1, 200) s(i)
on conflict do nothing
returning id;

-- insert N meetings (ONLY)
insert into meetings (meet_name, description, start_date, start_time)
select concat(left(md5(random()::text), 4), ' ', left(md5(random()::text), 5)) as meet_name,
       md5(random()::text)                                                     as description,
       cast(random_datetime() as date)                                         as start_date,
       cast(random_datetime() as time)                                         as start_time
from generate_series(1, 1000) s(i)
on conflict do nothing
returning id;

-- fill missing end_date and end_time
UPDATE meetings
SET end_date = start_date,
    end_time = start_time + '1 hours'
WHERE end_date IS NULL;

-- select random user id
Create or replace function random_valid_user_id() returns int as
$$
declare
begin
    return (select id from users order by random() limit 1);
end;
$$ language plpgsql;

--

-- select random meeting id
Create or replace function random_valid_meeting_id() returns int as
$$
declare
begin
    return (select id from meetings order by random() limit 1);
end;
$$ language plpgsql;

--

-- random status
Create or replace function random_status() returns text as
$$
declare
begin
    return (array [status_enum 'requested', status_enum 'approved', status_enum 'declined'])[floor(random() * 3 + 1)];
end;
$$ language plpgsql;

-- insert user_meetings links
insert into user_meetings (user_id, meeting_id, status)
select random_valid_user_id(), random_valid_meeting_id(), cast(random_status() as status_enum)
from generate_series(1, 2000) s(i)
on conflict do nothing
returning (user_id, meeting_id, status);

-- show average numbers of users for every meeting
select avg(t.count)
from (select count(um.user_id)
      from user_meetings um
               join meetings m on um.meeting_id = m.id
      group by um.meeting_id) t;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
truncate user_meetings;
truncate meetings cascade;
truncate users cascade;
-- +goose StatementEnd
