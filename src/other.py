import db
import sqlite3 as sl
import datetime

con = sl.connect('bhimps.db', detect_types=sl.PARSE_DECLTYPES)
db.set_timer(con, "pavlov", datetime.datetime.now())
