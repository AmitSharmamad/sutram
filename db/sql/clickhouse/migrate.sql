USE sutram;

INSERT INTO instrument
SELECT *
FROM file('instrument_data_*.csv', 'CSVWithNames');


INSERT INTO instrument_master
SELECT *
FROM file('instrument_master.csv', 'CSVWithNames');