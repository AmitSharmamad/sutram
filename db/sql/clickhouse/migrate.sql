USE sutram;

INSERT INTO instrument
SELECT *
FROM file('instrument.csv', 'CSVWithNames');


INSERT INTO instrument_master
SELECT *
FROM file('instrument_master.csv', 'CSVWithNames');