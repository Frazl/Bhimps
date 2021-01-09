def get_user_score(con, uid):
  c = con.cursor()
  data = c.execute("SELECT * FROM USERSCORES WHERE id = '%s'" % uid)
  user_score = data.fetchone()[1]
  return user_score

def get_scores(con, n=10, high=True):
  c = con.cursor()
  order = 'DESC'
  if not high: order='ASC'
  data = c.execute("SELECT * FROM USERSCORES WHERE score != 0 ORDER BY score %s LIMIT %s" % (order, n))
  res = []
  for row in data:
    res.append(row)
  return res

def update_user_score(con, uid, score):
  c = con.cursor()
  c.execute("UPDATE USERSCORES SET score= %s WHERE id = '%s'" % (score, uid))
  con.commit()

def get_timer(con, timer_id):
  c = con.cursor()
  data = c.execute("SELECT * FROM TIMERS WHERE id = '%s'" % timer_id)
  timestamp = data.fetchone()[1]
  return timestamp

def set_timer(con, timer_id, timestamp):
  c = con.cursor()
  c.execute("INSERT or REPLACE INTO TIMERS (id, time) VALUES(?, ?)", (timer_id, timestamp))
  con.commit()
