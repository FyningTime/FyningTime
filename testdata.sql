insert into workday (date) values ('2024-08-08');
insert into workday (date) values ('2024-08-09');

insert into worktime (type, time, workday) values ('Begin', '2024-08-08 07:56:47', (select id from workday where date = '2024-08-08'));
insert into worktime (type, time, workday) values ('End', '2024-08-08 12:10:13', (select id from workday where date = '2024-08-08'));
insert into worktime (type, time, workday) values ('Begin', '2024-08-08 12:45:33', (select id from workday where date = '2024-08-08'));
insert into worktime (type, time, workday) values ('End', '2024-08-08 16:30:12', (select id from workday where date = '2024-08-08'));

insert into worktime (type, time, workday) values ('Begin', '2024-08-09 07:56:47', (select id from workday where date = '2024-08-09'));
insert into worktime (type, time, workday) values ('End', '2024-08-09 12:10:13', (select id from workday where date = '2024-08-09'));
insert into worktime (type, time, workday) values ('Begin', '2024-08-09 12:45:33', (select id from workday where date = '2024-08-09'));
insert into worktime (type, time, workday) values ('End', '2024-08-09 16:30:12', (select id from workday where date = '2024-08-09'));