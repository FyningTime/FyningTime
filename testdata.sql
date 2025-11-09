insert into workday (date) values ('2025-08-08');
insert into workday (date) values ('2025-08-09');

insert into worktime (type, time, workday) values ('Begin', '2025-08-08 07:45:33', (select id from workday where date = '2025-08-08'));
insert into worktime (type, time, workday) values ('End', '2025-08-08 16:30:12', (select id from workday where date = '2025-08-08'));

insert into worktime (type, time, workday) values ('Begin', '2025-08-09 07:15:47', (select id from workday where date = '2025-08-09'));
insert into worktime (type, time, workday) values ('End', '2025-08-09 12:10:13', (select id from workday where date = '2025-08-09'));
insert into worktime (type, time, workday) values ('Begin', '2025-08-09 12:45:33', (select id from workday where date = '2025-08-09'));
insert into worktime (type, time, workday) values ('End', '2025-08-09 16:00:12', (select id from workday where date = '2025-08-09'));
