import sqlite3

conn = sqlite3.connect("database.db")
print("Connected to database succesfully")
c = conn.cursor()

c.execute("drop table sessions")
c.execute("drop table users")

conn.commit()
conn.close()
print("Migrated down")
